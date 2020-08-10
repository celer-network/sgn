package monitor

import (
	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/transactor"
	"github.com/celer-network/sgn/x/guard"
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
	var keys, vals [][]byte
	m.lock.RLock()
	iterator, err := m.db.Iterator(PullerKeyPrefix, storetypes.PrefixEndBytes(PullerKeyPrefix))
	if err != nil {
		log.Errorln("Create db iterator err", err)
		return
	}
	for ; iterator.Valid(); iterator.Next() {
		keys = append(keys, iterator.Key())
		vals = append(vals, iterator.Value())
	}
	iterator.Close()
	m.lock.RUnlock()

	for i, key := range keys {
		event := NewEventFromBytes(vals[i])
		log.Infoln("Process puller event", event.Name, "at mainchain block", event.Log.BlockNumber)
		err = m.dbDelete(key)
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
			e.Raw = event.Log
			m.syncIntendSettle(e)
		case *mainchain.CelerLedgerIntendWithdraw:
			e.Raw = event.Log
			m.syncIntendWithdrawChannel(e)
		}
	}
}

func (m *Monitor) syncDPoSValidatorChange(validatorChange *mainchain.DPoSValidatorChange) {
	log.Infof("puller queue process validator change %x", validatorChange.EthAddr)
	m.syncValidator(validatorChange.EthAddr)
}

func (m *Monitor) syncDPoSIntendWithdraw(intendWithdraw *mainchain.DPoSIntendWithdraw) {
	log.Infof("puller queue process intend withdraw %x", intendWithdraw.Candidate)
	m.syncValidator(intendWithdraw.Candidate)
	m.syncDelegator(intendWithdraw.Candidate, intendWithdraw.Delegator)
}

func (m *Monitor) syncDPoSCandidateUnbonded(candidateUnbonded *mainchain.DPoSCandidateUnbonded) {
	log.Infof("puller queue process candidate unbonded %x", candidateUnbonded.Candidate)
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
	requests, seqs := m.getGuardRequests(intendSettle.ChannelId)
	for i, request := range requests {
		m.triggerGuard(request, intendSettle.Raw, seqs[i], guard.ChanStatus_Settling)
	}
}

func (m *Monitor) syncIntendWithdrawChannel(intendWithdrawChannel *mainchain.CelerLedgerIntendWithdraw) {
	log.Infof("Sync intendWithdrawChannel %x, tx hash %x", intendWithdrawChannel.ChannelId, intendWithdrawChannel.Raw.TxHash)
	requests, seqs := m.getGuardRequests(intendWithdrawChannel.ChannelId)
	for i, request := range requests {
		m.triggerGuard(request, intendWithdrawChannel.Raw, seqs[i], guard.ChanStatus_Withdrawing)
	}
}

func (m *Monitor) triggerGuard(request *guard.Request, rawLog ethtypes.Log, seq uint64, guardStatus guard.ChanStatus) {
	if request.Status != guard.ChanStatus_Idle {
		log.Infoln("The guard state is not idle, current state", request.Status)
		return
	}
	trigger := guard.NewGuardTrigger(
		mainchain.Bytes2Cid(request.ChannelId),
		mainchain.Hex2Addr(request.SimplexReceiver),
		rawLog.TxHash,
		rawLog.BlockNumber,
		seq,
		guardStatus)
	syncData := m.operator.CliCtx.Codec.MustMarshalBinaryBare(trigger)
	msg := sync.NewMsgSubmitChange(sync.GuardTrigger, syncData, m.operator.Key.GetAddress())
	log.Infof("submit change tx: trigger guard request %s", trigger)
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

	vt := staking.Validator{
		Description: staking.Description{
			Identity: address.Hex(),
		},
		Tokens:     sdk.NewIntFromBigInt(ci.StakingPool).QuoRaw(common.TokenDec),
		Status:     mainchain.ParseStatus(ci),
		Commission: commission,
	}

	candidate, err := validator.CLIQueryCandidate(m.operator.CliCtx, validator.RouterKey, address.Hex())
	if err != nil {
		log.Errorln("sidechain query candidate err:", err)
		return
	}
	v, err := validator.CLIQueryValidator(
		m.operator.CliCtx, staking.RouterKey, candidate.Operator.String())
	if err == nil {
		if vt.Status.Equal(v.Status) && vt.Tokens.Equal(v.Tokens) &&
			vt.Commission.CommissionRates.Rate.Equal(v.Commission.CommissionRates.Rate) {
			log.Infof("no need to sync updated validator %x", address)
			return
		}
	}

	if m.ethClient.Address == address {
		pk, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, viper.GetString(common.FlagSgnPubKey))
		if err != nil {
			log.Errorln("GetConsPubKeyBech32 err:", err)
			return
		}

		vt.ConsPubKey = pk
	}

	validatorData := m.operator.CliCtx.Codec.MustMarshalBinaryBare(vt)
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
	logmsg := ""
	for _, transactor := range transactors {
		logmsg += transactor.String() + " "
	}
	log.Infoln("set transactors", logmsg)
	m.operator.AddTxMsg(setTransactorsMsg)
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
