package monitor

import (
	"fmt"
	"sync"
	"time"

	"github.com/allegro/bigcache"
	"github.com/celer-network/goutils/eth/monitor"
	"github.com/celer-network/goutils/eth/watcher"
	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/transactor"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/viper"
	dbm "github.com/tendermint/tm-db"
)

const (
	prefixMonitor = "mon"
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
	isValidator     bool
	dbLock          sync.RWMutex
}

func NewMonitor(ethClient *mainchain.EthClient, operator *transactor.Transactor, db dbm.DB) {
	monitorDb := dbm.NewPrefixDB(db, []byte(prefixMonitor))
	dal := newWatcherDAL(monitorDb)
	watchService := watcher.NewWatchService(ethClient.Client, dal, viper.GetUint64(common.FlagEthPollInterval))
	if watchService == nil {
		log.Fatalln("Cannot create watch service")
	}

	blkDelay := viper.GetUint64(common.FlagEthBlockDelay)
	ethMonitor := monitor.NewService(watchService, blkDelay, true /* enabled */)
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

	go m.checkBlockHead()

	go m.monitorDPoSDelegate()
	go m.monitorDPoSValidatorChange()
	go m.monitorDPoSIntendWithdraw()
	go m.monitorDPoSCandidateUnbonded()
	go m.monitorDPoSConfirmParamProposal()
	go m.monitorSGNUpdateSidechainAddr()
	go m.monitorCelerLedgerIntendSettle()
	go m.monitorCelerLedgerIntendWithdraw()

	go m.monitorSidechainWithdrawReward()
	go m.monitorSidechainSlash()
}

func (m *Monitor) checkBlockHead() {
	// TODO: configure check interval,
	// different queues could be checked at different times
	// e.g., guard queue does not need to be checked so frequently
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
		m.processQueue()
		m.verifyActiveChanges()
	}
}

func (m *Monitor) processQueue() {
	m.processPullerQueue()
	m.processGuardQueue()
	m.processPenaltyQueue()
}

func (m *Monitor) monitorSGNUpdateSidechainAddr() {
	_, err := m.ethMonitor.Monitor(
		&monitor.Config{
			EventName:  string(UpdateSidechainAddr),
			Contract:   m.sgnContract,
			StartBlock: m.ethMonitor.GetCurrentBlockNumber(),
		},
		func(cb monitor.CallbackID, eLog ethtypes.Log) {
			log.Infof("Catch event UpdateSidechainAddr, tx hash: %x", eLog.TxHash)
			event := NewEvent(UpdateSidechainAddr, eLog)
			dberr := m.dbSet(GetPullerKey(eLog), event.MustMarshal())
			if dberr != nil {
				log.Errorln("db Set err", dberr)
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
			dberr := m.dbSet(GetPullerKey(eLog), event.MustMarshal())
			if dberr != nil {
				log.Errorln("db Set err", dberr)
			}
		})
	if err != nil {
		log.Fatal(err)
	}
}

func (m *Monitor) monitorDPoSCandidateUnbonded() {
	_, err := m.ethMonitor.Monitor(
		&monitor.Config{
			EventName:  string(CandidateUnbonded),
			Contract:   m.dposContract,
			StartBlock: m.ethMonitor.GetCurrentBlockNumber(),
		},
		func(cb monitor.CallbackID, eLog ethtypes.Log) {
			log.Infof("Catch event CandidateUnbonded, tx hash: %x", eLog.TxHash)
			event := NewEvent(CandidateUnbonded, eLog)
			dberr := m.dbSet(GetPullerKey(eLog), event.MustMarshal())
			if dberr != nil {
				log.Errorln("db Set err", dberr)
			}
		})
	if err != nil {
		log.Fatal(err)
	}
}

func (m *Monitor) monitorDPoSConfirmParamProposal() {
	_, err := m.ethMonitor.Monitor(
		&monitor.Config{
			EventName:  string(ConfirmParamProposal),
			Contract:   m.dposContract,
			StartBlock: m.ethMonitor.GetCurrentBlockNumber(),
		},
		func(cb monitor.CallbackID, eLog ethtypes.Log) {
			log.Infof("Catch event ConfirmParamProposal, tx hash: %x", eLog.TxHash)
			event := NewEvent(ConfirmParamProposal, eLog)
			dberr := m.dbSet(GetPullerKey(eLog), event.MustMarshal())
			if dberr != nil {
				log.Errorln("db Set err", dberr)
			}
		})
	if err != nil {
		log.Fatal(err)
	}
}

func (m *Monitor) monitorDPoSValidatorChange() {
	_, err := m.ethMonitor.Monitor(
		&monitor.Config{
			EventName:  string(ValidatorChange),
			Contract:   m.dposContract,
			StartBlock: m.ethMonitor.GetCurrentBlockNumber(),
		},
		func(cb monitor.CallbackID, eLog ethtypes.Log) {
			logmsg := fmt.Sprintf("Catch event ValidatorChange, tx hash: %x", eLog.TxHash)
			validatorChange, perr := m.ethClient.DPoS.ParseValidatorChange(eLog)
			if perr != nil {
				log.Errorf("%s. parse event err: %s", logmsg, perr)
				return
			}
			if validatorChange.ChangeType == mainchain.AddValidator {
				// self init sync if add validator
				if validatorChange.EthAddr == m.ethClient.Address {
					log.Infof("%s. Init my own validator.", logmsg)
					m.isValidator = true
					m.syncValidator(validatorChange.EthAddr)
					m.setTransactors()
				}
			} else {
				// self only put removal event to puller queue
				log.Infof("%s, eth addr: %x", logmsg, validatorChange.EthAddr)
				if validatorChange.EthAddr == m.ethClient.Address {
					m.isValidator = false
				}
				event := NewEvent(ValidatorChange, eLog)
				dberr := m.dbSet(GetPullerKey(eLog), event.MustMarshal())
				if dberr != nil {
					log.Errorln("db Set err", dberr)
				}
			}
		})
	if err != nil {
		log.Fatal(err)
	}
}

func (m *Monitor) monitorDPoSIntendWithdraw() {
	_, err := m.ethMonitor.Monitor(
		&monitor.Config{
			EventName:  string(IntendWithdraw),
			Contract:   m.dposContract,
			StartBlock: m.ethMonitor.GetCurrentBlockNumber(),
		},
		func(cb monitor.CallbackID, eLog ethtypes.Log) {
			log.Infof("Catch event IntendWithdrawDpos, tx hash: %x", eLog.TxHash)
			event := NewEvent(IntendWithdrawDpos, eLog)
			dberr := m.dbSet(GetPullerKey(eLog), event.MustMarshal())
			if dberr != nil {
				log.Errorln("db Set err", dberr)
			}
		})
	if err != nil {
		log.Fatal(err)
	}
}

func (m *Monitor) monitorCelerLedgerIntendWithdraw() {
	_, err := m.ethMonitor.Monitor(
		&monitor.Config{
			EventName:  string(IntendWithdraw),
			Contract:   m.ledgerContract,
			StartBlock: m.ethMonitor.GetCurrentBlockNumber(),
		},
		func(cb monitor.CallbackID, eLog ethtypes.Log) {
			log.Infof("Catch event IntendWithdrawChannel, tx hash: %x", eLog.TxHash)
			event := NewEvent(IntendWithdrawChannel, eLog)
			dberr := m.dbSet(GetPullerKey(eLog), event.MustMarshal())
			if dberr != nil {
				log.Errorln("db Set err", dberr)
			}
			dberr = m.dbSet(GetGuardKey(eLog), event.MustMarshal())
			if dberr != nil {
				log.Errorln("db Set err", dberr)
			}
		})
	if err != nil {
		log.Fatal(err)
	}
}

func (m *Monitor) monitorCelerLedgerIntendSettle() {
	_, err := m.ethMonitor.Monitor(
		&monitor.Config{
			EventName:  string(IntendSettle),
			Contract:   m.ledgerContract,
			StartBlock: m.ethMonitor.GetCurrentBlockNumber(),
		},
		func(cb monitor.CallbackID, eLog ethtypes.Log) {
			log.Infof("Catch event IntendSettle, tx hash: %x", eLog.TxHash)
			event := NewEvent(IntendSettle, eLog)
			dberr := m.dbSet(GetPullerKey(eLog), event.MustMarshal())
			if dberr != nil {
				log.Errorln("db Set err", dberr)
			}
			dberr = m.dbSet(GetGuardKey(eLog), event.MustMarshal())
			if dberr != nil {
				log.Errorln("db Set err", dberr)
			}
		})
	if err != nil {
		log.Fatal(err)
	}
}

func (m *Monitor) monitorDPoSDelegate() {
	_, err := m.ethMonitor.Monitor(
		&monitor.Config{
			EventName:  string(Delegate),
			Contract:   m.dposContract,
			StartBlock: m.ethMonitor.GetCurrentBlockNumber(),
		},
		func(cb monitor.CallbackID, eLog ethtypes.Log) {
			log.Infof("Catch event Delegate, tx hash: %x", eLog.TxHash)
			delegate, perr := m.ethClient.DPoS.ParseDelegate(eLog)
			if perr != nil {
				log.Errorln("parse event err", perr)
				return
			}
			m.handleDPoSDelegate(delegate)
		})
	if err != nil {
		log.Fatal(err)
	}
}

func (m *Monitor) handleDPoSDelegate(delegate *mainchain.DPoSDelegate) {
	if delegate.Candidate != m.ethClient.Address {
		log.Tracef("Ignore delegate from delegator %x to candidate %x", delegate.Delegator, delegate.Candidate)
		return
	}

	log.Infof("Handle new delegate from delegator %x to candidate %x, new stake %s, pool %s",
		delegate.Delegator, delegate.Candidate, delegate.NewStake.String(), delegate.StakingPool.String())
	m.syncDelegator(delegate.Candidate, delegate.Delegator)

	if m.isValidator {
		m.syncValidator(delegate.Candidate)
	} else {
		m.claimValidatorOnMainchain()
	}
}

func (m *Monitor) claimValidatorOnMainchain() {
	candidate, err := m.ethClient.DPoS.GetCandidateInfo(&bind.CallOpts{}, m.ethClient.Address)
	if err != nil {
		log.Errorln("GetCandidateInfo err", err)
		return
	}
	if candidate.StakingPool.Cmp(candidate.MinSelfStake) == -1 {
		log.Debug("Not enough stake to become validator")
		return
	}

	minStake, err := m.ethClient.DPoS.GetMinStakingPool(&bind.CallOpts{})
	if err != nil {
		log.Errorln("GetMinStakingPool err", err)
		return
	}
	if candidate.StakingPool.Cmp(minStake) == -1 {
		log.Debug("Not enough stake to become validator")
		return
	}

	_, err = m.ethClient.Transactor.Transact(
		nil,
		func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*ethtypes.Transaction, error) {
			return m.ethClient.DPoS.ClaimValidator(opts)
		},
	)
	if err != nil {
		log.Errorln("ClaimValidator tx err", err)
		return
	}
	log.Infof("Claimed validator %x on mainchain", m.ethClient.Address)
}
