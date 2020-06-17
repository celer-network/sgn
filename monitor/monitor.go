package monitor

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
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
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/rpc/client"
	tTypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

const (
	prefixEthMonitor = "em"
)

var (
	initiateWithdrawRewardEvent = fmt.Sprintf("%s.%s='%s'", validator.ModuleName, sdk.AttributeKeyAction, validator.ActionInitiateWithdraw)
	slashEvent                  = fmt.Sprintf("%s.%s='%s'", slash.EventTypeSlash, sdk.AttributeKeyAction, slash.ActionPenalty)
)

type EthMonitor struct {
	ethClient       *mainchain.EthClient
	operator        *transactor.Transactor
	db              dbm.DB
	monitorService  *monitor.Service
	dposContract    monitor.Contract
	sgnContract     monitor.Contract
	ledgerContract  monitor.Contract
	verifiedChanges *bigcache.BigCache
	blkNum          *big.Int
	secureBlkNum    uint64
	isValidator     bool
}

func NewEthMonitor(ethClient *mainchain.EthClient, operator *transactor.Transactor, db dbm.DB) {
	monitorDb := dbm.NewPrefixDB(db, []byte(prefixEthMonitor))
	dal := newWatcherDAL(monitorDb)
	watchService := watcher.NewWatchService(ethClient.Client, dal, viper.GetUint64(common.FlagEthPollInterval))
	if watchService == nil {
		log.Fatalln("Cannot create watch service")
	}

	monitorService := monitor.NewService(watchService, 0 /* blockDelay */, true /* enabled */)
	monitorService.Init()

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

	m := EthMonitor{
		ethClient:       ethClient,
		operator:        operator,
		db:              db,
		monitorService:  monitorService,
		blkNum:          monitorService.GetCurrentBlockNumber(),
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
	go m.monitorWithdrawReward()
	go m.monitorSlash()
}

func (m *EthMonitor) monitorBlockHead() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		<-ticker.C
		blkNum := m.monitorService.GetCurrentBlockNumber()
		if blkNum.Cmp(m.blkNum) == 0 {
			continue
		}

		m.blkNum = blkNum
		m.secureBlkNum = blkNum.Uint64() - viper.GetUint64(common.FlagEthConfirmCount)
		m.processQueue()
		m.verifyActiveChanges()
	}
}

func (m *EthMonitor) monitorUpdateSidechainAddr() {
	_, err := m.monitorService.Monitor(
		&monitor.Config{
			EventName:  string(UpdateSidechainAddr),
			Contract:   m.sgnContract,
			StartBlock: m.blkNum,
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

func (m *EthMonitor) monitorConfirmParamProposal() {
	_, err := m.monitorService.Monitor(
		&monitor.Config{
			EventName:  string(ConfirmParamProposal),
			Contract:   m.dposContract,
			StartBlock: m.blkNum,
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

func (m *EthMonitor) monitorDelegate() {
	_, err := m.monitorService.Monitor(
		&monitor.Config{
			EventName:  string(Delegate),
			Contract:   m.dposContract,
			StartBlock: m.blkNum,
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

func (m *EthMonitor) monitorCandidateUnbonded() {
	_, err := m.monitorService.Monitor(
		&monitor.Config{
			EventName:  string(CandidateUnbonded),
			Contract:   m.dposContract,
			StartBlock: m.blkNum,
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

func (m *EthMonitor) monitorValidatorChange() {
	_, err := m.monitorService.Monitor(
		&monitor.Config{
			EventName:  string(ValidatorChange),
			Contract:   m.dposContract,
			StartBlock: m.blkNum,
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

func (m *EthMonitor) monitorIntendWithdrawSgn() {
	_, err := m.monitorService.Monitor(
		&monitor.Config{
			EventName:  string(IntendWithdraw),
			Contract:   m.dposContract,
			StartBlock: m.blkNum,
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

func (m *EthMonitor) monitorIntendWithdrawChannel() {
	_, err := m.monitorService.Monitor(
		&monitor.Config{
			EventName:  string(IntendWithdraw),
			Contract:   m.ledgerContract,
			StartBlock: m.blkNum,
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

func (m *EthMonitor) monitorIntendSettle() {
	_, err := m.monitorService.Monitor(
		&monitor.Config{
			EventName:  string(IntendSettle),
			Contract:   m.ledgerContract,
			StartBlock: m.blkNum,
		},
		func(cb monitor.CallbackID, eLog ethtypes.Log) {
			log.Infof("Catch event IntendSettle, tx hash: %x", eLog.TxHash)
			event := NewEvent(IntendSettle, eLog)
			err2 := m.db.Set(GetPullerKey(eLog), event.MustMarshal())
			if err2 != nil {
				log.Errorln("db Set err", err2)
			}
			err2 = m.db.Set(GetPusherKey(eLog), event.MustMarshal())
			if err2 != nil {
				log.Errorln("db Set err", err2)
			}
		})
	if err != nil {
		log.Fatal(err)
	}
}

func (m *EthMonitor) monitorWithdrawReward() {
	m.monitorTendermintEvent(initiateWithdrawRewardEvent, func(e abci.Event) {
		if !m.isValidator {
			return
		}

		event := sdk.StringifyEvent(e)
		if event.Attributes[0].Value == validator.ActionInitiateWithdraw {
			m.handleInitiateWithdrawReward(event.Attributes[1].Value)
		}
	})
}

func (m *EthMonitor) monitorSlash() {
	m.monitorTendermintEvent(slashEvent, func(e abci.Event) {
		if !m.isValidator {
			return
		}

		event := sdk.StringifyEvent(e)

		if event.Attributes[0].Value == slash.ActionPenalty {
			nonce, err := strconv.ParseUint(event.Attributes[1].Value, 10, 64)
			if err != nil {
				log.Errorln("Parse penalty nonce error", err)
				return
			}

			penaltyEvent := NewPenaltyEvent(nonce)
			m.handlePenalty(penaltyEvent)
			err = m.db.Set(GetPenaltyKey(penaltyEvent.Nonce), penaltyEvent.MustMarshal())
			if err != nil {
				log.Errorln("db Set err", err)
			}
		}
	})
}

func (m *EthMonitor) monitorTendermintEvent(eventTag string, handleEvent func(event abci.Event)) {
	client, err := client.NewHTTP(m.operator.CliCtx.NodeURI, "/websocket")
	if err != nil {
		log.Errorln("Fail to start create http client", err)
		return
	}

	err = client.Start()
	if err != nil {
		log.Errorln("Fail to start ws client", err)
		return
	}
	defer client.Stop()

	txs, err := client.Subscribe(context.Background(), "monitor", eventTag)
	if err != nil {
		log.Errorln("ws client subscribe error", err)
		return
	}

	for e := range txs {
		switch data := e.Data.(type) {
		case tTypes.EventDataNewBlock:
			for _, event := range data.ResultBeginBlock.Events {
				handleEvent(event)
			}
			for _, event := range data.ResultEndBlock.Events {
				handleEvent(event)
			}
		case tTypes.EventDataTx:
			for _, event := range data.TxResult.Result.Events {
				handleEvent(event)
			}
		}
	}
}
