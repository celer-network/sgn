package monitor

import (
	"math/big"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/transactor"
	"github.com/celer-network/sgn/x/global"
	"github.com/celer-network/sgn/x/slash"
	"github.com/celer-network/sgn/x/sync"
	"github.com/celer-network/sgn/x/validator"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/spf13/viper"
)

func (m *EthMonitor) processEventQueue(secureBlockNum uint64) {
	iterator, err := m.db.Iterator(EventKeyPrefix, storetypes.PrefixEndBytes(EventKeyPrefix))
	if err != nil {
		log.Errorln("Create db iterator err", err)
		return
	}
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		event := NewEventFromBytes(iterator.Value())
		if secureBlockNum < event.Log.BlockNumber {
			continue
		}

		log.Infoln("Process mainchain event", event.Name, "at mainchain block", event.Log.BlockNumber)
		m.db.Delete(iterator.Key())

		switch e := event.ParseEvent(m.ethClient).(type) {
		case *mainchain.DPoSDelegate:
			m.handleDelegate(e)
		case *mainchain.DPoSValidatorChange:
			m.handleValidatorChange(e)
		case *mainchain.DPoSIntendWithdraw:
			m.handleIntendWithdraw(e)
		}
	}
}

func (m *EthMonitor) handleNewBlock(blkNum *big.Int) {
	log.Infoln("Catch new mainchain block", blkNum)
	if !m.isPuller() {
		return
	}

	log.Infof("Add MsgSyncBlock %d to transactor msgQueue", blkNum)
	block := global.NewBlock(blkNum.Uint64())
	blockData := m.blockSyncer.CliCtx.Codec.MustMarshalBinaryBare(block)
	msg := sync.NewMsgSubmitChange(sync.SyncBlock, blockData, m.blockSyncer.Key.GetAddress())
	m.blockSyncer.AddTxMsg(msg)
}

func (m *EthMonitor) handleDelegate(delegate *mainchain.DPoSDelegate) {
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
		m.claimValidatorOnMainchain()
	}
}

func (m *EthMonitor) handleValidatorChange(validatorChange *mainchain.DPoSValidatorChange) {
	log.Infof("New validator change %x type %d", validatorChange.EthAddr, validatorChange.ChangeType)
	isAddValidator := validatorChange.ChangeType == mainchain.AddValidator
	doSync := m.isPuller() && !isAddValidator

	if validatorChange.EthAddr == m.ethClient.Address {
		m.isValidator = isAddValidator
		doSync = true

		if m.isValidator {
			m.setTransactors()
		}
	}

	if doSync {
		m.syncValidator(validatorChange.EthAddr)
	}
}

func (m *EthMonitor) handleIntendWithdraw(intendWithdraw *mainchain.DPoSIntendWithdraw) {
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

func (m *EthMonitor) handlePenalty(penaltyEvent PenaltyEvent) {
	penalty, err := slash.CLIQueryPenalty(m.operator.CliCtx, slash.StoreKey, penaltyEvent.Nonce)
	if err != nil {
		log.Errorf("Query penalty %d err %s", penaltyEvent.Nonce, err)
		return
	}
	log.Infof("New penalty to %s, reason %s, nonce %d", penalty.ValidatorAddr, penalty.Reason, penaltyEvent.Nonce)

	sig, err := m.ethClient.SignMessage(penalty.PenaltyProtoBytes)
	if err != nil {
		log.Errorln("SignMessage err", err)
		return
	}

	msg := slash.NewMsgSignPenalty(penaltyEvent.Nonce, sig, m.operator.Key.GetAddress())
	m.operator.AddTxMsg(msg)
}

func (m *EthMonitor) claimValidatorOnMainchain() {
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

	_, err = m.ethClient.DPoS.ClaimValidator(m.ethClient.Auth)
	if err != nil {
		log.Errorln("ClaimValidator tx err", err)
		return
	}
	log.Infof("Claimed validator %x on mainchain", m.ethClient.Address)
}

func (m *EthMonitor) setTransactors() {
	log.Infoln("Set transactor")
	transactors, err := transactor.ParseTransactorAddrs(viper.GetStringSlice(common.FlagSgnTransactors))
	if err != nil {
		log.Errorln("parse transactors err", err)
		return
	}
	setTransactorsMsg := validator.NewMsgSetTransactors(
		mainchain.Addr2Hex(m.ethClient.Address),
		transactors,
		m.operator.Key.GetAddress(),
	)
	m.operator.AddTxMsg(setTransactorsMsg)
}

func (m *EthMonitor) syncValidator(address mainchain.Addr) {
	log.Infof("SyncValidator %x", address)
	ci, err := m.ethClient.DPoS.GetCandidateInfo(&bind.CallOpts{}, address)
	if err != nil {
		log.Errorln("Failed to query candidate info:", err)
		return
	}

	validator := staking.Validator{
		Description: staking.Description{
			Identity: address.Hex(),
		},
		Tokens: sdk.NewIntFromBigInt(ci.StakingPool).QuoRaw(common.TokenDec),
		Status: mainchain.ParseStatus(ci),
	}

	if m.ethClient.Address == address {
		pk, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, viper.GetString(common.FlagSgnPubKey))
		if err != nil {
			log.Errorln("GetConsPubKeyBech32 err:", err)
			return
		}

		validator.ConsPubKey = pk
	}

	validatorData := m.operator.CliCtx.Codec.MustMarshalBinaryBare(validator)
	msg := sync.NewMsgSubmitChange(sync.SyncValidator, validatorData, m.operator.Key.GetAddress())
	m.operator.AddTxMsg(msg)
}

func (m *EthMonitor) syncDelegator(candidatorAddr, delegatorAddr mainchain.Addr) {
	log.Infof("SyncDelegator candidate: %x, delegator: %x", candidatorAddr, delegatorAddr)

	di, err := m.ethClient.DPoS.GetDelegatorInfo(&bind.CallOpts{}, candidatorAddr, delegatorAddr)
	if err != nil {
		log.Errorf("Failed to query delegator info: %s", err)
		return
	}

	delegator := validator.NewDelegator(mainchain.Addr2Hex(candidatorAddr), mainchain.Addr2Hex(delegatorAddr))
	delegator.DelegatedStake = sdk.NewIntFromBigInt(di.DelegatedStake)
	delegatorData := m.operator.CliCtx.Codec.MustMarshalBinaryBare(delegator)
	msg := sync.NewMsgSubmitChange(sync.SyncDelegator, delegatorData, m.operator.Key.GetAddress())
	m.operator.AddTxMsg(msg)
}
