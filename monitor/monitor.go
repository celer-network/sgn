package monitor

import (
	"fmt"
	"time"

	"github.com/allegro/bigcache"
	"github.com/celer-network/goutils/eth/monitor"
	"github.com/celer-network/goutils/eth/watcher"
	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/transactor"
	"github.com/celer-network/sgn/x/slash"
	"github.com/celer-network/sgn/x/validator"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/viper"
	dbm "github.com/tendermint/tm-db"
)

const (
	prefixMonitor = "mon"
)

var (
	initiateWithdrawRewardEvent = fmt.Sprintf("%s.%s='%s'", validator.ModuleName, sdk.AttributeKeyAction, validator.ActionInitiateWithdraw)
	slashEvent                  = fmt.Sprintf("%s.%s='%s'", slash.EventTypeSlash, sdk.AttributeKeyAction, slash.ActionPenalty)
)

type Monitor struct {
	ethClient       *mainchain.EthClient
	operator        *transactor.Transactor
	db              dbm.DB
	ethMonitor      *monitor.Service
	dposContract    monitor.Contract
	sgnContract     monitor.Contract
	ledgerContract  monitor.Contract
	verifiedChanges *bigcache.BigCache
	secureBlkNum    uint64
	isValidator     bool
}

func NewMonitor(ethClient *mainchain.EthClient, operator *transactor.Transactor, db dbm.DB) {
	monitorDb := dbm.NewPrefixDB(db, []byte(prefixMonitor))
	dal := newWatcherDAL(monitorDb)
	watchService := watcher.NewWatchService(ethClient.Client, dal, viper.GetUint64(common.FlagEthPollInterval))
	if watchService == nil {
		log.Fatalln("Cannot create watch service")
	}

	ethMonitor := monitor.NewService(watchService, 0 /* blockDelay */, true /* enabled */)
	ethMonitor.Init()

	dposCandidateInfo, err := ethClient.DPoS.GetCandidateInfo(&bind.CallOpts{}, ethClient.Address)
	if err != nil {
		log.Fatalln("GetCandidateInfo err", err)
	}

	dposContract := NewMonitorContractInfo(ethClient.DPoSAddress, mainchain.DPoSABI)
	sgnContract := NewMonitorContractInfo(ethClient.SGNAddress, mainchain.SGNABI)
	ledgerContract := NewMonitorContractInfo(ethClient.LedgerAddress, mainchain.CelerLedgerABI)

	verifiedChanges, err := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
	if err != nil {
		log.Fatalln("NewBigCache err", err)
	}

	m := Monitor{
		ethClient:       ethClient,
		operator:        operator,
		db:              db,
		ethMonitor:      ethMonitor,
		dposContract:    dposContract,
		sgnContract:     sgnContract,
		ledgerContract:  ledgerContract,
		verifiedChanges: verifiedChanges,
		isValidator:     mainchain.IsBonded(dposCandidateInfo),
	}

	go m.monitorBlockHead()
	go m.monitorUpdateSidechainAddr()
	go m.monitorDelegate()
	go m.monitorCandidateUnbonded()
	go m.monitorValidatorChange()
	go m.monitorIntendWithdrawSgn()
	go m.monitorIntendSettle()
	go m.monitorIntendWithdrawChannel()
	go m.monitorWithdrawReward()
	go m.monitorSlash()
}

func (m *Monitor) monitorBlockHead() {
	// TODO: configure check interval
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	blkNum := m.ethMonitor.GetCurrentBlockNumber().Uint64()
	for {
		<-ticker.C
		newblk := m.ethMonitor.GetCurrentBlockNumber().Uint64()
		if blkNum == newblk {
			continue
		}

		blkNum = newblk
		m.secureBlkNum = blkNum - viper.GetUint64(common.FlagEthConfirmCount)
		m.processQueue()
		m.verifyActiveChanges()
	}
}

func (m *Monitor) processQueue() {
	m.processEventQueue()
	m.processPullerQueue()
	m.processGuardQueue()
	m.processPenaltyQueue()
}

func (m *Monitor) monitorUpdateSidechainAddr() {
	_, err := m.ethMonitor.Monitor(
		&monitor.Config{
			EventName:  string(UpdateSidechainAddr),
			Contract:   m.sgnContract,
			StartBlock: m.ethMonitor.GetCurrentBlockNumber(),
		},
		func(cb monitor.CallbackID, eLog ethtypes.Log) {
			log.Infof("Catch event UpdateSidechainAddr, tx hash: %x", eLog.TxHash)
			event := NewEvent(UpdateSidechainAddr, eLog)
			err2 := m.db.Set(GetPullerKey(eLog), event.MustMarshal())
			if err2 != nil {
				log.Errorln("db Set err", err2)
			}
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}

func (m *Monitor) monitorConfirmParamProposal() {
	_, err := m.ethMonitor.Monitor(
		&monitor.Config{
			EventName:  string(ConfirmParamProposal),
			Contract:   m.dposContract,
			StartBlock: m.ethMonitor.GetCurrentBlockNumber(),
		},
		func(cb monitor.CallbackID, eLog ethtypes.Log) {
			log.Infof("Catch event ConfirmParamProposal, tx hash: %x", eLog.TxHash)
			event := NewEvent(ConfirmParamProposal, eLog)
			err2 := m.db.Set(GetPullerKey(eLog), event.MustMarshal())
			if err2 != nil {
				log.Errorln("db Set err", err2)
			}
		})
	if err != nil {
		log.Fatal(err)
	}
}

func (m *Monitor) monitorDelegate() {
	_, err := m.ethMonitor.Monitor(
		&monitor.Config{
			EventName:  string(Delegate),
			Contract:   m.dposContract,
			StartBlock: m.ethMonitor.GetCurrentBlockNumber(),
		},
		func(cb monitor.CallbackID, eLog ethtypes.Log) {
			log.Infof("Catch event Delegate, tx hash: %x", eLog.TxHash)
			event := NewEvent(Delegate, eLog)
			err2 := m.db.Set(GetEventKey(eLog), event.MustMarshal())
			if err2 != nil {
				log.Errorln("db Set err", err2)
			}
		})
	if err != nil {
		log.Fatal(err)
	}
}

func (m *Monitor) monitorCandidateUnbonded() {
	_, err := m.ethMonitor.Monitor(
		&monitor.Config{
			EventName:  string(CandidateUnbonded),
			Contract:   m.dposContract,
			StartBlock: m.ethMonitor.GetCurrentBlockNumber(),
		},
		func(cb monitor.CallbackID, eLog ethtypes.Log) {
			log.Infof("Catch event CandidateUnbonded, tx hash: %x", eLog.TxHash)
			event := NewEvent(CandidateUnbonded, eLog)
			err2 := m.db.Set(GetEventKey(eLog), event.MustMarshal())
			if err2 != nil {
				log.Errorln("db Set err", err2)
			}
		})
	if err != nil {
		log.Fatal(err)
	}
}

func (m *Monitor) monitorValidatorChange() {
	_, err := m.ethMonitor.Monitor(
		&monitor.Config{
			EventName:  string(ValidatorChange),
			Contract:   m.dposContract,
			StartBlock: m.ethMonitor.GetCurrentBlockNumber(),
		},
		func(cb monitor.CallbackID, eLog ethtypes.Log) {
			log.Infof("Catch event ValidatorChange, tx hash: %x", eLog.TxHash)
			event := NewEvent(ValidatorChange, eLog)
			err2 := m.db.Set(GetEventKey(eLog), event.MustMarshal())
			if err2 != nil {
				log.Errorln("db Set err", err2)
			}
		})
	if err != nil {
		log.Fatal(err)
	}
}

func (m *Monitor) monitorIntendWithdrawSgn() {
	_, err := m.ethMonitor.Monitor(
		&monitor.Config{
			EventName:  string(IntendWithdraw),
			Contract:   m.dposContract,
			StartBlock: m.ethMonitor.GetCurrentBlockNumber(),
		},
		func(cb monitor.CallbackID, eLog ethtypes.Log) {
			log.Infof("Catch event IntendWithdrawSgn, tx hash: %x", eLog.TxHash)
			event := NewEvent(IntendWithdrawSgn, eLog)
			err2 := m.db.Set(GetEventKey(eLog), event.MustMarshal())
			if err2 != nil {
				log.Errorln("db Set err", err2)
			}
		})
	if err != nil {
		log.Fatal(err)
	}
}

func (m *Monitor) monitorIntendWithdrawChannel() {
	_, err := m.ethMonitor.Monitor(
		&monitor.Config{
			EventName:  string(IntendWithdraw),
			Contract:   m.ledgerContract,
			StartBlock: m.ethMonitor.GetCurrentBlockNumber(),
		},
		func(cb monitor.CallbackID, eLog ethtypes.Log) {
			log.Infof("Catch event IntendWithdrawChannel, tx hash: %x", eLog.TxHash)
			event := NewEvent(IntendWithdrawChannel, eLog)
			err2 := m.db.Set(GetEventKey(eLog), event.MustMarshal())
			if err2 != nil {
				log.Errorln("db Set err", err2)
			}
		})
	if err != nil {
		log.Fatal(err)
	}
}

func (m *Monitor) monitorIntendSettle() {
	_, err := m.ethMonitor.Monitor(
		&monitor.Config{
			EventName:  string(IntendSettle),
			Contract:   m.ledgerContract,
			StartBlock: m.ethMonitor.GetCurrentBlockNumber(),
		},
		func(cb monitor.CallbackID, eLog ethtypes.Log) {
			log.Infof("Catch event IntendSettle, tx hash: %x", eLog.TxHash)
			event := NewEvent(IntendSettle, eLog)
			err2 := m.db.Set(GetPullerKey(eLog), event.MustMarshal())
			if err2 != nil {
				log.Errorln("db Set err", err2)
			}
			err2 = m.db.Set(GetGuardKey(eLog), event.MustMarshal())
			if err2 != nil {
				log.Errorln("db Set err", err2)
			}
		})
	if err != nil {
		log.Fatal(err)
	}
}
