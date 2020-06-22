package monitor

import (
	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/transactor"
	"github.com/celer-network/sgn/x/subscribe"
	"github.com/celer-network/sgn/x/sync"
	"github.com/celer-network/sgn/x/validator"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/viper"
)

func (m *Monitor) processPullerQueue() {
	if !m.isPuller() {
		return
	}

	iterator, err := m.db.Iterator(PullerKeyPrefix, storetypes.PrefixEndBytes(PullerKeyPrefix))
	if err != nil {
		log.Errorln("Create db iterator err", err)
		return
	}
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		event := NewEventFromBytes(iterator.Value())
		log.Infoln("Process puller event", event.Name, "at mainchain block", event.Log.BlockNumber)
		err = m.db.Delete(iterator.Key())
		if err != nil {
			log.Errorln("db Delete err", err)
			continue
		}

		switch e := event.ParseEvent(m.ethClient).(type) {
		case *mainchain.DPoSValidatorChange:
			m.syncDPoSValidatorChange(e)
		case *mainchain.DPoSIntendWithdraw:
			m.syncDPoSIntendWithdraw(e)
		case *mainchain.DPoSCandidateUnbonded:
			m.syncDPoSCandidateUnbonded(e)
		case *mainchain.DPoSConfirmParamProposal:
			m.syncConfirmParamProposal(e)
		case *mainchain.SGNUpdateSidechainAddr:
			m.syncUpdateSidechainAddr(e)
		case *mainchain.CelerLedgerIntendSettle:
			m.syncIntendSettle(e)
		case *mainchain.CelerLedgerIntendWithdraw:
			m.syncIntendWithdrawChannel(e)
		}
	}
}

func (m *Monitor) syncDPoSValidatorChange(validatorChange *mainchain.DPoSValidatorChange) {
	log.Infof("New validator change %x type %d", validatorChange.EthAddr, validatorChange.ChangeType)
	m.syncValidator(validatorChange.EthAddr)
}

func (m *Monitor) syncDPoSIntendWithdraw(intendWithdraw *mainchain.DPoSIntendWithdraw) {
	log.Infof("New intend withdraw %x", intendWithdraw.Candidate)
	m.syncValidator(intendWithdraw.Candidate)
}

func (m *Monitor) syncDPoSCandidateUnbonded(candidateUnbonded *mainchain.DPoSCandidateUnbonded) {
	log.Infof("New candidate unbonded %x", candidateUnbonded.Candidate)
	m.syncValidator(candidateUnbonded.Candidate)
}

func (m *Monitor) syncConfirmParamProposal(confirmParamProposal *mainchain.DPoSConfirmParamProposal) {
	paramChange := common.NewParamChange(sdk.NewIntFromBigInt(confirmParamProposal.Record), sdk.NewIntFromBigInt(confirmParamProposal.NewValue))
	paramChangeData := m.operator.CliCtx.Codec.MustMarshalBinaryBare(paramChange)
	msg := sync.NewMsgSubmitChange(sync.ConfirmParamProposal, paramChangeData, m.operator.Key.GetAddress())
	log.Infof("submit change tx: confirm param proposal Record %v, NewValue %v", confirmParamProposal.Record, confirmParamProposal.NewValue)
	m.operator.AddTxMsg(msg)
}

func (m *Monitor) syncUpdateSidechainAddr(updateSidechainAddr *mainchain.SGNUpdateSidechainAddr) {
	sidechainAddr, err := m.ethClient.SGN.SidechainAddrMap(&bind.CallOpts{}, updateSidechainAddr.Candidate)
	if err != nil {
		log.Errorln("Query sidechain address error:", err)
		return
	}

	c, err := validator.CLIQueryCandidate(m.operator.CliCtx, validator.RouterKey, mainchain.Addr2Hex(updateSidechainAddr.Candidate))
	if err == nil && sdk.AccAddress(sidechainAddr).Equals(c.Operator) {
		log.Infof("The sidechain address of candidate %x has been updated", updateSidechainAddr.Candidate)
		return
	}

	candidate := validator.NewCandidate(updateSidechainAddr.Candidate.Hex(), sdk.AccAddress(sidechainAddr))
	candidateData := m.operator.CliCtx.Codec.MustMarshalBinaryBare(candidate)
	msg := sync.NewMsgSubmitChange(sync.UpdateSidechainAddr, candidateData, m.operator.Key.GetAddress())
	log.Infof("submit change tx: update sidechain addr for candidate %s %s", candidate.EthAddress, candidate.Operator.String())
	m.operator.AddTxMsg(msg)
}

func (m *Monitor) syncIntendSettle(intendSettle *mainchain.CelerLedgerIntendSettle) {
	log.Infof("Sync IntendSettle %x, tx hash %x", intendSettle.ChannelId, intendSettle.Raw.TxHash)
	requests := m.getRequests(intendSettle.ChannelId)
	for _, request := range requests {
		m.triggerGuard(request, intendSettle.Raw)
	}
}

func (m *Monitor) syncIntendWithdrawChannel(intendWithdrawChannel *mainchain.CelerLedgerIntendWithdraw) {
	log.Infof("Sync intendWithdrawChannel %x, tx hash %x", intendWithdrawChannel.ChannelId, intendWithdrawChannel.Raw.TxHash)
	requests := m.getRequests(intendWithdrawChannel.ChannelId)
	for _, request := range requests {
		m.triggerGuard(request, intendWithdrawChannel.Raw)
	}
}

func (m *Monitor) triggerGuard(request subscribe.Request, rawLog ethtypes.Log) {
	if request.TriggerTxHash != "" {
		log.Infoln("The intendSettle event has been synced on sgn")
		return
	}

	disputeTimeout, err := m.ethClient.Ledger.GetDisputeTimeout(&bind.CallOpts{}, mainchain.Bytes2Cid(request.ChannelId))
	if err != nil {
		log.Errorln("GetDisputeTimeout err:", err)
		return
	}

	request.DisputeTimeout = disputeTimeout.Uint64()
	request.TriggerTxHash = rawLog.TxHash.Hex()
	request.TriggerTxBlkNum = rawLog.BlockNumber
	requestData := m.operator.CliCtx.Codec.MustMarshalBinaryBare(request)
	msg := sync.NewMsgSubmitChange(sync.TriggerGuard, requestData, m.operator.Key.GetAddress())
	log.Infof("submit change tx: trigger guard request %s", request)
	m.operator.AddTxMsg(msg)
}

func (m *Monitor) syncValidator(address mainchain.Addr) {
	ci, err := m.ethClient.DPoS.GetCandidateInfo(&bind.CallOpts{}, address)
	if err != nil {
		log.Errorln("Failed to query candidate info:", err)
		return
	}

	commission, err := common.NewCommission(m.ethClient, ci.CommissionRate)
	if err != nil {
		log.Errorln("Failed to create new commission:", err)
		return
	}

	validator := staking.Validator{
		Description: staking.Description{
			Identity: address.Hex(),
		},
		Tokens:     sdk.NewIntFromBigInt(ci.StakingPool).QuoRaw(common.TokenDec),
		Status:     mainchain.ParseStatus(ci),
		Commission: commission,
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
	log.Infof("submit change tx: sync validator %x", address)
	m.operator.AddTxMsg(msg)
}

func (m *Monitor) setTransactors() {
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
	log.Infoln("set transactors", transactors)
	m.operator.AddTxMsg(setTransactorsMsg)
}

func (m *Monitor) handleDPoSDelegate(delegate *mainchain.DPoSDelegate) {
	if delegate.Candidate != m.ethClient.Address {
		log.Tracef("Ignore delegate from delegator %x to candidate %x", delegate.Delegator, delegate.Candidate)
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

func (m *Monitor) syncDelegator(candidatorAddr, delegatorAddr mainchain.Addr) {
	di, err := m.ethClient.DPoS.GetDelegatorInfo(&bind.CallOpts{}, candidatorAddr, delegatorAddr)
	if err != nil {
		log.Errorf("Failed to query delegator info: %s", err)
		return
	}

	delegator := validator.NewDelegator(mainchain.Addr2Hex(candidatorAddr), mainchain.Addr2Hex(delegatorAddr))
	delegator.DelegatedStake = sdk.NewIntFromBigInt(di.DelegatedStake)
	delegatorData := m.operator.CliCtx.Codec.MustMarshalBinaryBare(delegator)
	msg := sync.NewMsgSubmitChange(sync.SyncDelegator, delegatorData, m.operator.Key.GetAddress())
	log.Infof("submit change tx: sync delegator %x candidate %x stake %s", delegatorAddr, candidatorAddr, delegator.DelegatedStake)
	m.operator.AddTxMsg(msg)
}
