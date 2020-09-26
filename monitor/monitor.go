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
	*Operator
	db              dbm.DB
	ethMonitor      *monitor.Service
	dposContract    monitor.Contract
	sgnContract     monitor.Contract
	ledgerContract  monitor.Contract
	verifiedChanges *bigcache.BigCache
	sidechainAcct   sdk.AccAddress
	bonded          bool
	executeSlash    bool
	bootstrapped    bool // SGN has bootstrapped with at least one bonded validator on the mainchain contract
	lock            sync.RWMutex
}

func NewMonitor(operator *Operator, db dbm.DB) {
	monitorDb := dbm.NewPrefixDB(db, []byte(prefixMonitor))
	dal := newWatcherDAL(monitorDb)
	watchService := watcher.NewWatchService(operator.EthClient.Client, dal, viper.GetUint64(common.FlagEthPollInterval))
	if watchService == nil {
		log.Fatalln("Cannot create watch service")
	}

	blkDelay := viper.GetUint64(common.FlagEthBlockDelay)
	ethMonitor := monitor.NewService(watchService, blkDelay, true /* enabled */)
	ethMonitor.Init()

	dposCandidateInfo, err := operator.EthClient.DPoS.GetCandidateInfo(&bind.CallOpts{}, operator.EthClient.Address)
	if err != nil {
		log.Fatalln("GetCandidateInfo err", err)
	}

	valnum, err := operator.EthClient.DPoS.GetValidatorNum(&bind.CallOpts{})
	if err != nil {
		log.Fatalln("GetValidatorNum err", err)
	}

	dposContract := NewMonitorContractInfo(operator.EthClient.DPoSAddress, mainchain.DPoSABI)
	sgnContract := NewMonitorContractInfo(operator.EthClient.SGNAddress, mainchain.SGNABI)
	ledgerContract := NewMonitorContractInfo(operator.EthClient.LedgerAddress, mainchain.CelerLedgerABI)

	verifiedChanges, err := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
	if err != nil {
		log.Fatalln("NewBigCache err", err)
	}

	m := Monitor{
		Operator:        operator,
		db:              db,
		ethMonitor:      ethMonitor,
		dposContract:    dposContract,
		sgnContract:     sgnContract,
		ledgerContract:  ledgerContract,
		verifiedChanges: verifiedChanges,
		bonded:          mainchain.IsBonded(dposCandidateInfo),
		bootstrapped:    valnum.Uint64() > 0,
		executeSlash:    viper.GetBool(common.FlagSgnExecuteSlash),
	}
	m.sidechainAcct, err = sdk.AccAddressFromBech32(viper.GetString(common.FlagSgnValidatorAccount))
	if err != nil {
		log.Fatalln("Sidechain acct error")
	}

	go m.processQueues()

	go m.monitorDPoSValidatorChange()
	go m.monitorDPoSUpdateDelegatedStake()
	go m.monitorDPoSCandidateUnbonded()
	go m.monitorDPoSConfirmParamProposal()
	go m.monitorDPoSUpdateCommissionRate()
	go m.monitorSGNUpdateSidechainAddr()
	go m.monitorSGNAddSubscriptionBalance()
	go m.monitorCelerLedgerIntendSettle()
	go m.monitorCelerLedgerIntendWithdraw()

	go m.monitorSidechainCreateValidator()
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
			Reset:         true,
			CheckInterval: getEventCheckInterval(UpdateSidechainAddr),
		},
		func(cb monitor.CallbackID, eLog ethtypes.Log) {
			log.Infof("Catch event UpdateSidechainAddr, tx hash: %x", eLog.TxHash)
			event := NewEvent(UpdateSidechainAddr, eLog)
			dberr := m.dbSet(GetPullerKey(eLog), event.MustMarshal())
			if dberr != nil {
				log.Errorln("db Set err", dberr)
			}
			if !m.isBonded() {
				e, perr := m.EthClient.SGN.ParseUpdateSidechainAddr(eLog)
				if perr != nil {
					log.Errorln("parse event err", perr)
					return
				}
				if e.Candidate == m.EthClient.Address && m.shouldClaimValidator() {
					m.claimValidatorOnMainchain()
				}
			}
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}

func (m *Monitor) monitorSGNAddSubscriptionBalance() {
	_, err := m.ethMonitor.Monitor(
		&monitor.Config{
			EventName:     string(AddSubscriptionBalance),
			Contract:      m.sgnContract,
			StartBlock:    m.getCurrentBlockNumber(),
			Reset:         true,
			CheckInterval: getEventCheckInterval(AddSubscriptionBalance),
		},
		func(cb monitor.CallbackID, eLog ethtypes.Log) {
			log.Infof("Catch event AddSubscriptionBalance, tx hash: %x", eLog.TxHash)
			event := NewEvent(AddSubscriptionBalance, eLog)
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

func (m *Monitor) monitorDPoSCandidateUnbonded() {
	_, err := m.ethMonitor.Monitor(
		&monitor.Config{
			EventName:     string(CandidateUnbonded),
			Contract:      m.dposContract,
			StartBlock:    m.getCurrentBlockNumber(),
			Reset:         true,
			CheckInterval: getEventCheckInterval(CandidateUnbonded),
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
			EventName:     string(ConfirmParamProposal),
			Contract:      m.dposContract,
			StartBlock:    m.getCurrentBlockNumber(),
			Reset:         true,
			CheckInterval: getEventCheckInterval(ConfirmParamProposal),
		},
		func(cb monitor.CallbackID, eLog ethtypes.Log) {
			log.Infof("Catch event ConfirmParamProposal, tx hash: %x, blknum: %d", eLog.TxHash, eLog.BlockNumber)
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

func (m *Monitor) monitorDPoSUpdateCommissionRate() {
	_, err := m.ethMonitor.Monitor(
		&monitor.Config{
			EventName:     string(UpdateCommissionRate),
			Contract:      m.dposContract,
			StartBlock:    m.getCurrentBlockNumber(),
			Reset:         true,
			CheckInterval: getEventCheckInterval(UpdateCommissionRate),
		},
		func(cb monitor.CallbackID, eLog ethtypes.Log) {
			log.Infof("Catch event UpdateCommissionRate, tx hash: %x, blknum: %d", eLog.TxHash, eLog.BlockNumber)
			event := NewEvent(UpdateCommissionRate, eLog)
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
			EventName:     string(ValidatorChange),
			Contract:      m.dposContract,
			StartBlock:    m.getCurrentBlockNumber(),
			Reset:         true,
			CheckInterval: getEventCheckInterval(ValidatorChange),
		},
		func(cb monitor.CallbackID, eLog ethtypes.Log) {
			logmsg := fmt.Sprintf("Catch event ValidatorChange, tx hash: %x, blknum: %d", eLog.TxHash, eLog.BlockNumber)
			validatorChange, perr := m.EthClient.DPoS.ParseValidatorChange(eLog)
			if perr != nil {
				log.Errorf("%s. parse event err: %s", logmsg, perr)
				return
			}
			if validatorChange.ChangeType == mainchain.AddValidator {
				m.setBootstrapped()
				// self init sync if add validator
				if validatorChange.EthAddr == m.EthClient.Address {
					log.Infof("%s. Init my own validator.", logmsg)
					m.setBonded()
					go m.selfSyncValidator()
				} else {
					log.Infof("%s, addValidator addr: %x, ", logmsg, validatorChange.EthAddr)
				}
			} else {
				// self only put removal event to puller queue
				log.Infof("%s, removeValidator addr: %x, ", logmsg, validatorChange.EthAddr)
				if validatorChange.EthAddr == m.EthClient.Address {
					m.setUnbonded()
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

func (m *Monitor) monitorDPoSUpdateDelegatedStake() {
	_, err := m.ethMonitor.Monitor(
		&monitor.Config{
			EventName:     string(UpdateDelegatedStake),
			Contract:      m.dposContract,
			StartBlock:    m.getCurrentBlockNumber(),
			Reset:         true,
			CheckInterval: getEventCheckInterval(UpdateDelegatedStake),
		},
		func(cb monitor.CallbackID, eLog ethtypes.Log) {
			log.Infof("Catch event UpdateDelegatedStake, tx hash: %x, blknum: %d", eLog.TxHash, eLog.BlockNumber)
			event := NewEvent(UpdateDelegatedStake, eLog)
			dberr := m.dbSet(GetPullerKey(eLog), event.MustMarshal())
			if dberr != nil {
				log.Errorln("db Set err", dberr)
			}
			if !m.isBonded() {
				e, perr := m.EthClient.DPoS.ParseUpdateDelegatedStake(eLog)
				if perr != nil {
					log.Errorln("parse event err", perr)
					return
				}
				if e.Candidate == m.EthClient.Address && m.shouldClaimValidator() {
					m.claimValidatorOnMainchain()
				}
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
			Reset:         true,
			CheckInterval: getEventCheckInterval(IntendSettle),
		},
		func(cb monitor.CallbackID, eLog ethtypes.Log) {
			log.Infof("Catch event IntendSettle, tx hash: %x, blknum: %d", eLog.TxHash, eLog.BlockNumber)
			err := m.dbSet(GetPullerKey(eLog), NewEvent(IntendSettle, eLog).MustMarshal())
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
			Reset:         true,
			CheckInterval: getEventCheckInterval(IntendWithdrawChannel),
		},
		func(cb monitor.CallbackID, eLog ethtypes.Log) {
			log.Infof("Catch event IntendWithdrawChannel, tx hash: %x, blknum: %d", eLog.TxHash, eLog.BlockNumber)
			event := NewEvent(IntendWithdrawChannel, eLog).MustMarshal()
			dberr := m.dbSet(GetPullerKey(eLog), event)
			if dberr != nil {
				log.Errorln("db Set err", dberr)
			}
			m.setGuardEvent(eLog, ChanInfoState_CaughtWithdraw)
		})
	if err != nil {
		log.Fatal(err)
	}
}

const (
	selfSyncRetryNum         int = 5
	selfSyncRetryIntervalSec int = 60
)

func (m *Monitor) selfSyncValidator() {
	var i int
	for i = 1; i < selfSyncRetryNum; i++ {
		updated := m.SyncValidator(m.EthClient.Address)
		if updated {
			return
		}
		time.Sleep(time.Duration(selfSyncRetryIntervalSec) * time.Second)
	}
	log.Warn("self validator not synced yet")
}

func (m *Monitor) shouldClaimValidator() bool {
	candidate, err := m.EthClient.DPoS.GetCandidateInfo(&bind.CallOpts{}, m.EthClient.Address)
	if err != nil {
		log.Errorln("GetCandidateInfo err", err)
		return false
	}

	if !candidate.Initialized {
		log.Debug("Candidate not initialized on mainchain")
		return false
	}

	if mainchain.IsBonded(candidate) {
		log.Debug("Already bonded on mainchain")
		return false
	}

	minStake, err := m.EthClient.DPoS.GetMinStakingPool(&bind.CallOpts{})
	if err != nil {
		log.Errorln("GetMinStakingPool err", err)
		return false
	}
	if candidate.StakingPool.Cmp(minStake) == -1 {
		log.Debugf("Not enough stake to become a validator, my pool: %s, current min pool: %s", candidate.StakingPool, minStake)
		return false
	}

	delegator, err := m.EthClient.DPoS.GetDelegatorInfo(&bind.CallOpts{}, m.EthClient.Address, m.EthClient.Address)
	if err != nil {
		log.Errorln("GetDelegatorInfo err", err)
		return false
	}
	if delegator.DelegatedStake.Cmp(candidate.MinSelfStake) == -1 {
		log.Debugf("Not enough self-delegate stake, current: %s, require: %s", delegator.DelegatedStake, candidate.MinSelfStake)
		return false
	}

	minStakeInPool, err := m.EthClient.DPoS.GetUIntValue(&bind.CallOpts{}, big.NewInt(mainchain.MinStakeInPool))
	if err != nil {
		log.Errorln("Get MinStakeInPool param err", err)
		return false
	}
	if candidate.StakingPool.Cmp(minStakeInPool) == -1 {
		log.Debugf("Not enough stake to become a validator, my pool: %s, required min pool: %s", candidate.StakingPool, minStakeInPool)
		return false
	}

	sidechainAddr, err := m.EthClient.SGN.SidechainAddrMap(&bind.CallOpts{}, m.EthClient.Address)
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
	_, err := m.EthClient.Transactor.Transact(
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
			return m.EthClient.DPoS.ClaimValidator(opts)
		},
	)
	if err != nil {
		log.Errorln("ClaimValidator tx err", err)
		return
	}
	log.Infof("Claimed validator %x on mainchain", m.EthClient.Address)
}
