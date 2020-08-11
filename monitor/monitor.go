package monitor

import (
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/allegro/bigcache"
	"github.com/celer-network/goutils/eth"
	"github.com/celer-network/goutils/eth/monitor"
	"github.com/celer-network/goutils/eth/watcher"
	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/transactor"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	sidechainAcct   sdk.AccAddress
	bonded          bool
	executeSlash    bool
	lock            sync.RWMutex
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
		bonded:          mainchain.IsBonded(dposCandidateInfo),
		executeSlash:    viper.GetBool(common.FlagSgnExecuteSlash),
	}
	m.sidechainAcct, err = sdk.AccAddressFromBech32(viper.GetString(common.FlagSgnOperator))
	if err != nil {
		log.Fatalln("Sidechain acct error")
	}

	go m.processQueues()

	go m.monitorDPoSDelegate()
	go m.monitorDPoSValidatorChange()
	go m.monitorDPoSIntendWithdraw()
	go m.monitorDPoSCandidateUnbonded()
	go m.monitorDPoSConfirmParamProposal()
	go m.monitorSGNUpdateSidechainAddr()
	go m.monitorCelerLedgerIntendSettle()
	go m.monitorCelerLedgerIntendWithdraw()

	go m.monitorSidechainWithdrawReward()
	if m.executeSlash {
		go m.monitorSidechainSlash()
	}
}

func (m *Monitor) processQueues() {
	ticker := time.NewTicker(time.Duration(viper.GetUint64(common.FlagEthPollInterval)) * time.Second)
	defer ticker.Stop()

	blkNum := m.getCurrentBlockNumber().Uint64()
	for {
		<-ticker.C
		newblk := m.getCurrentBlockNumber().Uint64()
		if blkNum == newblk {
			continue
		}

		blkNum = newblk

		m.processPullerQueue()
		m.processGuardQueue()
		m.verifyActiveChanges()
		if m.executeSlash {
			m.processPenaltyQueue()
		}
	}
}

func (m *Monitor) monitorSGNUpdateSidechainAddr() {
	_, err := m.ethMonitor.Monitor(
		&monitor.Config{
			EventName:     string(UpdateSidechainAddr),
			Contract:      m.sgnContract,
			StartBlock:    m.getCurrentBlockNumber(),
			CheckInterval: eventCheckInterval(UpdateSidechainAddr),
		},
		func(cb monitor.CallbackID, eLog ethtypes.Log) {
			log.Infof("Catch event UpdateSidechainAddr, tx hash: %x", eLog.TxHash)
			event := NewEvent(UpdateSidechainAddr, eLog)
			dberr := m.dbSet(GetPullerKey(eLog.TxHash), event.MustMarshal())
			if dberr != nil {
				log.Errorln("db Set err", dberr)
			}
			if !m.isBonded() {
				e, perr := m.ethClient.SGN.ParseUpdateSidechainAddr(eLog)
				if perr != nil {
					log.Errorln("parse event err", perr)
					return
				}
				if e.Candidate == m.ethClient.Address && m.shouldClaimValidator() {
					m.claimValidatorOnMainchain()
				}
			}
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}

func (m *Monitor) monitorDPoSCandidateUnbonded() {
	_, err := m.ethMonitor.Monitor(
		&monitor.Config{
			EventName:     string(CandidateUnbonded),
			Contract:      m.dposContract,
			StartBlock:    m.getCurrentBlockNumber(),
			CheckInterval: eventCheckInterval(CandidateUnbonded),
		},
		func(cb monitor.CallbackID, eLog ethtypes.Log) {
			log.Infof("Catch event CandidateUnbonded, tx hash: %x", eLog.TxHash)
			event := NewEvent(CandidateUnbonded, eLog)
			dberr := m.dbSet(GetPullerKey(eLog.TxHash), event.MustMarshal())
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
			EventName:     string(ConfirmParamProposal),
			Contract:      m.dposContract,
			StartBlock:    m.getCurrentBlockNumber(),
			CheckInterval: eventCheckInterval(ConfirmParamProposal),
		},
		func(cb monitor.CallbackID, eLog ethtypes.Log) {
			log.Infof("Catch event ConfirmParamProposal, tx hash: %x", eLog.TxHash)
			event := NewEvent(ConfirmParamProposal, eLog)
			dberr := m.dbSet(GetPullerKey(eLog.TxHash), event.MustMarshal())
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
			EventName:     string(ValidatorChange),
			Contract:      m.dposContract,
			StartBlock:    m.getCurrentBlockNumber(),
			CheckInterval: eventCheckInterval(ValidatorChange),
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
					m.setBonded()
					m.syncValidator(validatorChange.EthAddr)
					m.setTransactors()
				}
			} else {
				// self only put removal event to puller queue
				log.Infof("%s, eth addr: %x", logmsg, validatorChange.EthAddr)
				if validatorChange.EthAddr == m.ethClient.Address {
					m.setUnbonded()
				}
				event := NewEvent(ValidatorChange, eLog)
				dberr := m.dbSet(GetPullerKey(eLog.TxHash), event.MustMarshal())
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
			EventName:     string(IntendWithdraw),
			Contract:      m.dposContract,
			StartBlock:    m.getCurrentBlockNumber(),
			CheckInterval: eventCheckInterval(IntendWithdrawDpos),
		},
		func(cb monitor.CallbackID, eLog ethtypes.Log) {
			log.Infof("Catch event IntendWithdrawDpos, tx hash: %x", eLog.TxHash)
			event := NewEvent(IntendWithdrawDpos, eLog)
			dberr := m.dbSet(GetPullerKey(eLog.TxHash), event.MustMarshal())
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
			EventName:     string(IntendSettle),
			Contract:      m.ledgerContract,
			StartBlock:    m.getCurrentBlockNumber(),
			CheckInterval: eventCheckInterval(IntendSettle),
		},
		func(cb monitor.CallbackID, eLog ethtypes.Log) {
			log.Infof("Catch event IntendSettle, tx hash: %x", eLog.TxHash)
			err := m.dbSet(GetPullerKey(eLog.TxHash), NewEvent(IntendSettle, eLog).MustMarshal())
			if err != nil {
				log.Errorln("db Set err", err)
			}
			m.setGuardEvent(eLog, ChanInfoState_CaughtSettle)
		})
	if err != nil {
		log.Fatal(err)
	}
}

func (m *Monitor) monitorCelerLedgerIntendWithdraw() {
	_, err := m.ethMonitor.Monitor(
		&monitor.Config{
			EventName:     string(IntendWithdraw),
			Contract:      m.ledgerContract,
			StartBlock:    m.getCurrentBlockNumber(),
			CheckInterval: eventCheckInterval(IntendWithdrawChannel),
		},
		func(cb monitor.CallbackID, eLog ethtypes.Log) {
			log.Infof("Catch event IntendWithdrawChannel, tx hash: %x", eLog.TxHash)
			err := m.dbSet(GetPullerKey(eLog.TxHash), NewEvent(IntendWithdrawChannel, eLog).MustMarshal())
			if err != nil {
				log.Errorln("db Set err", err)
			}
			m.setGuardEvent(eLog, ChanInfoState_CaughtWithdraw)
		})
	if err != nil {
		log.Fatal(err)
	}
}

func (m *Monitor) monitorDPoSDelegate() {
	_, err := m.ethMonitor.Monitor(
		&monitor.Config{
			EventName:     string(Delegate),
			Contract:      m.dposContract,
			StartBlock:    m.getCurrentBlockNumber(),
			CheckInterval: eventCheckInterval(Delegate),
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

	if m.isBonded() {
		m.syncValidator(delegate.Candidate)
	} else if m.shouldClaimValidator() {
		m.claimValidatorOnMainchain()
	}
}

func (m *Monitor) shouldClaimValidator() bool {
	candidate, err := m.ethClient.DPoS.GetCandidateInfo(&bind.CallOpts{}, m.ethClient.Address)
	if err != nil {
		log.Errorln("GetCandidateInfo err", err)
		return false
	}

	if !candidate.Initialized {
		log.Debug("Candidate not initialized on mainchain")
		return false
	}

	if mainchain.IsBonded(candidate) {
		log.Infoln("Already bonded on mainchain")
		return false
	}

	minStake, err := m.ethClient.DPoS.GetMinStakingPool(&bind.CallOpts{})
	if err != nil {
		log.Errorln("GetMinStakingPool err", err)
		return false
	}
	if candidate.StakingPool.Cmp(minStake) == -1 {
		log.Debugf("Not enough stake to become a validator, my pool: %s, current min pool: %s", candidate.StakingPool, minStake)
		return false
	}

	delegator, err := m.ethClient.DPoS.GetDelegatorInfo(&bind.CallOpts{}, m.ethClient.Address, m.ethClient.Address)
	if err != nil {
		log.Errorln("GetDelegatorInfo err", err)
		return false
	}
	if delegator.DelegatedStake.Cmp(candidate.MinSelfStake) == -1 {
		log.Debugf("Not enough self-delegate stake, current: %s, require: %s", delegator.DelegatedStake, candidate.MinSelfStake)
		return false
	}

	minStakeInPool, err := m.ethClient.DPoS.GetUIntValue(&bind.CallOpts{}, big.NewInt(mainchain.MinStakeInPool))
	if err != nil {
		log.Errorln("Get MinStakeInPool param err", err)
		return false
	}
	if candidate.StakingPool.Cmp(minStakeInPool) == -1 {
		log.Debugf("Not enough stake to become a validator, my pool: %s, required min pool: %s", candidate.StakingPool, minStakeInPool)
		return false
	}

	sidechainAddr, err := m.ethClient.SGN.SidechainAddrMap(&bind.CallOpts{}, m.ethClient.Address)
	if err != nil {
		log.Errorln("Query sidechain address error:", err)
		return false
	}
	if !sdk.AccAddress(sidechainAddr).Equals(m.sidechainAcct) {
		log.Debugf("sidechain address not match, %s %s", sdk.AccAddress(sidechainAddr), m.sidechainAcct)
		return false
	}

	return true
}

func (m *Monitor) claimValidatorOnMainchain() {
	_, err := m.ethClient.Transactor.Transact(
		&eth.TransactionStateHandler{
			OnMined: func(receipt *ethtypes.Receipt) {
				if receipt.Status == ethtypes.ReceiptStatusSuccessful {
					log.Infof("ClaimValidator transaction %x succeeded", receipt.TxHash)
				} else {
					log.Errorf("ClaimValidator transaction %x failed", receipt.TxHash)
				}
			},
			OnError: func(tx *ethtypes.Transaction, err error) {
				log.Errorf("ClaimValidator transaction %x err: %s", tx.Hash(), err)
			},
		},
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
