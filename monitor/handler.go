package monitor

import (
	"time"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/transactor"
	"github.com/celer-network/sgn/x/global"
	"github.com/celer-network/sgn/x/slash"
	"github.com/celer-network/sgn/x/validator"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/viper"
)

func (m *EthMonitor) processEventQueue() {
	secureBlockNum, err := m.getSecureBlockNum()
	if err != nil {
		log.Errorln("Query secureBlockNum err", err)
		return
	}

	iterator := m.db.Iterator(EventKeyPrefix, storetypes.PrefixEndBytes(EventKeyPrefix))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		event := NewEventFromBytes(iterator.Value())
		if secureBlockNum < event.Log.BlockNumber {
			continue
		}

		log.Infoln("Process mainchain event", event.Name, "at mainchain block", event.Log.BlockNumber)
		m.db.Delete(iterator.Key())

		switch e := event.ParseEvent(m.ethClient).(type) {
		case *mainchain.GuardDelegate:
			m.handleDelegate(e)
		case *mainchain.GuardValidatorChange:
			m.handleValidatorChange(e)
		case *mainchain.GuardIntendWithdraw:
			m.handleIntendWithdraw(e)
		}
	}
}

func (m *EthMonitor) handleNewBlock(header *types.Header) {
	log.Infoln("Catch new mainchain block", header.Number)
	if !m.isPuller() {
		return
	}

	params, err := m.getGlobalParams()
	if err != nil {
		log.Errorln("Query global params", err)
		return
	}

	time.Sleep(time.Duration(viper.GetInt64(common.FlagSgnTimeoutCommit)+params.BlkTimeDiffLower) * time.Second)

	log.Infof("Add MsgSyncBlock %d to transactor msgQueue", header.Number)
	transactor := m.tsPool.GetTransactor()
	msg := global.NewMsgSyncBlock(header.Number.Uint64(), transactor.Key.GetAddress())
	transactor.AddTxMsg(msg)
}

func (m *EthMonitor) handleDelegate(delegate *mainchain.GuardDelegate) {
	if delegate.Candidate != m.ethClient.Address {
		log.Infof("Ignore delegate from delegator %x to candidate %x", delegate.Delegator, delegate.Candidate)
		return
	}

	log.Infof("Handle new delegate from delegator %x to candidate %x, stake %s pool %s",
		delegate.Delegator, delegate.Candidate, delegate.NewStake.String(), delegate.StakingPool.String())
	m.syncDelegator(delegate.Candidate, delegate.Delegator)

	if m.isValidator {
		m.syncValidator(delegate.Candidate)
	} else {
		m.claimValidatorOnMainchain(delegate)
	}
}

func (m *EthMonitor) handleValidatorChange(validatorChange *mainchain.GuardValidatorChange) {
	log.Infof("New validator change %x type %d", validatorChange.EthAddr, validatorChange.ChangeType)
	doSync := m.isPuller()

	if validatorChange.EthAddr == m.ethClient.Address {
		m.isValidator = validatorChange.ChangeType == mainchain.AddValidator
		if m.isValidator {
			m.claimValidator()
			return
		}

		doSync = true
	}

	if doSync {
		m.syncValidator(validatorChange.EthAddr)
	}
}

func (m *EthMonitor) handleIntendWithdraw(intendWithdraw *mainchain.GuardIntendWithdraw) {
	log.Infof("New intend withdraw %x", intendWithdraw.Candidate)

	if m.isPullerOrOwner(intendWithdraw.Candidate) {
		m.syncValidator(intendWithdraw.Candidate)
	}
}

func (m *EthMonitor) handleInitiateWithdrawReward(ethAddr string) {
	log.Infoln("New initiate withdraw", ethAddr)

	reward, err := validator.CLIQueryReward(m.operator.CliCtx, validator.StoreKey, ethAddr)
	if err != nil {
		log.Errorln("Query reward err", err)
		return
	}

	sig, err := m.ethClient.SignMessage(reward.RewardProtoBytes)
	if err != nil {
		log.Errorln("SignMessage err", err)
		return
	}

	msg := validator.NewMsgSignReward(ethAddr, sig, m.operator.Key.GetAddress())
	m.operator.AddTxMsg(msg)
}

func (m *EthMonitor) handlePenalty(nonce uint64) {
	log.Infoln("New Penalty", nonce)

	penalty, err := slash.CLIQueryPenalty(m.operator.CliCtx, slash.StoreKey, nonce)
	if err != nil {
		log.Errorln("Query penalty err", err)
		return
	}

	sig, err := m.ethClient.SignMessage(penalty.PenaltyProtoBytes)
	if err != nil {
		log.Errorln("SignMessage err", err)
		return
	}

	msg := slash.NewMsgSignPenalty(nonce, sig, m.operator.Key.GetAddress())
	m.operator.AddTxMsg(msg)
}

func (m *EthMonitor) claimValidatorOnMainchain(delegate *mainchain.GuardDelegate) {
	minStake, err := m.ethClient.Guard.GetMinStakingPool(&bind.CallOpts{})
	if err != nil {
		log.Errorln("GetMinStakingPool err", err)
		return
	}

	if delegate.StakingPool.Uint64() <= minStake.Uint64() {
		log.Debug("Not enough stake to become validator")
		return
	}

	_, err = m.ethClient.Guard.ClaimValidator(m.ethClient.Auth)
	if err != nil {
		log.Errorln("ClaimValidator tx err", err)
		return
	}
	log.Infof("Claimed validator %x on mainchain", delegate.Candidate)
}

func (m *EthMonitor) claimValidator() {
	log.Infof("Claim self as a validator on sidechain, self address %x", m.ethClient.Address)
	transactors, err := transactor.ParseTransactorAddrs(viper.GetStringSlice(common.FlagSgnTransactors))
	if err != nil {
		log.Errorln("parse transactors err", err)
		return
	}

	msg := validator.NewMsgClaimValidator(
		mainchain.Addr2Hex(m.ethClient.Address),
		viper.GetString(common.FlagSgnPubKey),
		transactors,
		m.operator.Key.GetAddress(),
	)
	m.operator.AddTxMsg(msg)
	m.addTransactorsToPool()
}

func (m *EthMonitor) syncValidator(address mainchain.Addr) {
	log.Infof("SyncValidator %x", address)
	transactor := m.tsPool.GetTransactor()
	msg := validator.NewMsgSyncValidator(mainchain.Addr2Hex(address), transactor.Key.GetAddress())
	transactor.AddTxMsg(msg)
}

func (m *EthMonitor) syncDelegator(candidatorAddr, delegatorAddr mainchain.Addr) {
	log.Infof("SyncDelegator candidate: %x, delegator: %x", candidatorAddr, delegatorAddr)

	transactor := m.tsPool.GetTransactor()
	msg := validator.NewMsgSyncDelegator(
		mainchain.Addr2Hex(candidatorAddr), mainchain.Addr2Hex(delegatorAddr), transactor.Key.GetAddress())
	transactor.AddTxMsg(msg)
}
