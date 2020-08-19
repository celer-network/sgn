package monitor

import (
	"fmt"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/x/guard"
	"github.com/celer-network/sgn/x/sync"
	"github.com/celer-network/sgn/x/validator"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
		logmsg := fmt.Sprintf("Process puller event %s at mainchain block %d", event.Name, event.Log.BlockNumber)
		err = m.dbDelete(key)
		if err != nil {
			log.Errorf("%s. db Delete err: %s", logmsg, err)
			continue
		}

		switch e := event.ParseEvent(m.EthClient).(type) {
		case *mainchain.DPoSValidatorChange:
			m.syncDPoSValidatorChange(e, logmsg)
		case *mainchain.DPoSIntendWithdraw:
			m.syncDPoSIntendWithdraw(e, logmsg)
		case *mainchain.DPoSCandidateUnbonded:
			m.syncDPoSCandidateUnbonded(e, logmsg)
		case *mainchain.DPoSConfirmParamProposal:
			m.syncConfirmParamProposal(e, logmsg)
		case *mainchain.DPoSUpdateCommissionRate:
			m.syncUpdateCommissionRate(e, logmsg)
		case *mainchain.SGNUpdateSidechainAddr:
			m.syncUpdateSidechainAddr(e, logmsg)
		case *mainchain.SGNAddSubscriptionBalance:
			m.syncSGNAddSubscriptionBalance(e, logmsg)
		case *mainchain.CelerLedgerIntendSettle:
			e.Raw = event.Log
			m.syncIntendSettle(e, logmsg)
		case *mainchain.CelerLedgerIntendWithdraw:
			e.Raw = event.Log
			m.syncIntendWithdrawChannel(e, logmsg)
		}
	}
}

func (m *Monitor) syncDPoSValidatorChange(validatorChange *mainchain.DPoSValidatorChange, logmsg string) {
	log.Infof("%s. validator change %x type %d", logmsg, validatorChange.EthAddr, validatorChange.ChangeType)
	m.SyncValidator(validatorChange.EthAddr)
}

func (m *Monitor) syncDPoSIntendWithdraw(intendWithdraw *mainchain.DPoSIntendWithdraw, logmsg string) {
	log.Infof("%s. intend withdraw candidate %x delegator %x amount %s",
		logmsg, intendWithdraw.Candidate, intendWithdraw.Delegator, intendWithdraw.WithdrawAmount)
	m.SyncValidator(intendWithdraw.Candidate)
	m.SyncDelegator(intendWithdraw.Candidate, intendWithdraw.Delegator)
}

func (m *Monitor) syncDPoSCandidateUnbonded(candidateUnbonded *mainchain.DPoSCandidateUnbonded, logmsg string) {
	log.Infof("%s. candidate unbonded %x", logmsg, candidateUnbonded.Candidate)
	m.SyncValidator(candidateUnbonded.Candidate)
}

func (m *Monitor) syncUpdateCommissionRate(commission *mainchain.DPoSUpdateCommissionRate, logmsg string) {
	log.Infof("%s. commission update %x, %s", logmsg, commission.Candidate, commission.NewRate)
	m.SyncValidator(commission.Candidate)
}

func (m *Monitor) syncUpdateSidechainAddr(sidechainAddr *mainchain.SGNUpdateSidechainAddr, logmsg string) {
	log.Infof("%s. sidechainAddr update %x, %s",
		logmsg, sidechainAddr.Candidate, sdk.AccAddress(sidechainAddr.NewSidechainAddr.Bytes()))
	m.SyncUpdateSidechainAddr(sidechainAddr.Candidate)
}

func (m *Monitor) syncConfirmParamProposal(confirmParamProposal *mainchain.DPoSConfirmParamProposal, logmsg string) {
	paramChange := common.NewParamChange(sdk.NewIntFromBigInt(confirmParamProposal.Record), sdk.NewIntFromBigInt(confirmParamProposal.NewValue))
	paramChangeData := m.Transactor.CliCtx.Codec.MustMarshalBinaryBare(paramChange)
	msg := m.Transactor.NewMsgSubmitChange(sync.ConfirmParamProposal, paramChangeData, m.EthClient.Client)
	log.Infof("%s. submit change tx: confirm param proposal Record %v, NewValue %v", logmsg, confirmParamProposal.Record, confirmParamProposal.NewValue)
	m.Transactor.AddTxMsg(msg)
}

func (m *Monitor) syncIntendSettle(intendSettle *mainchain.CelerLedgerIntendSettle, logmsg string) {
	log.Infof("%s. sync IntendSettle %x, tx hash %x", logmsg, intendSettle.ChannelId, intendSettle.Raw.TxHash)
	requests, seqs := m.getGuardRequests(intendSettle.ChannelId)
	for i, request := range requests {
		m.triggerGuard(request, intendSettle.Raw, seqs[i], guard.ChanStatus_Settling)
	}
}

func (m *Monitor) syncIntendWithdrawChannel(intendWithdrawChannel *mainchain.CelerLedgerIntendWithdraw, logmsg string) {
	log.Infof("%s. sync intendWithdrawChannel %x, tx hash %x", logmsg, intendWithdrawChannel.ChannelId, intendWithdrawChannel.Raw.TxHash)
	requests, seqs := m.getGuardRequests(intendWithdrawChannel.ChannelId)
	for i, request := range requests {
		m.triggerGuard(request, intendWithdrawChannel.Raw, seqs[i], guard.ChanStatus_Withdrawing)
	}
}

func (m *Monitor) syncSGNAddSubscriptionBalance(event *mainchain.SGNAddSubscriptionBalance, logmsg string) {
	transactor := m.Transactor
	consumer := event.Consumer
	consumerEthAddress := consumer.Hex()
	amount := event.Amount
	amountInt := sdk.NewIntFromBigInt(amount)
	subscription, err := guard.CLIQuerySubscription(transactor.CliCtx, guard.RouterKey, consumerEthAddress)
	if err == nil {
		if subscription.Deposit.Equal(amountInt) {
			log.Infof("%s. subscription already updated for %s, amount %s", logmsg, consumerEthAddress, amount)
			return
		}
	}
	subscription = guard.NewSubscription(consumer.Hex())
	subscription.Deposit = sdk.NewIntFromBigInt(amount)
	subscriptionData := transactor.CliCtx.Codec.MustMarshalBinaryBare(subscription)
	msg := m.Transactor.NewMsgSubmitChange(sync.Subscribe, subscriptionData, m.EthClient.Client)
	log.Infof("%s. submit change tx: subscribe ethAddress %s, amount %s, mainchain tx hash %x", logmsg, consumerEthAddress, amount, event.Raw.TxHash)
	transactor.AddTxMsg(msg)
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
	syncData := m.Transactor.CliCtx.Codec.MustMarshalBinaryBare(trigger)
	msg := m.Transactor.NewMsgSubmitChange(sync.GuardTrigger, syncData, m.EthClient.Client)
	log.Infof("submit change tx: trigger guard request %s", trigger)
	m.Transactor.AddTxMsg(msg)
}

func (m *Monitor) setTransactors() {
	transactors, err := common.ParseTransactorAddrs(viper.GetStringSlice(common.FlagSgnTransactors))
	if err != nil {
		log.Errorln("parse transactors err", err)
		return
	}
	setTransactorsMsg := validator.NewMsgSetTransactors(
		mainchain.Addr2Hex(m.EthClient.Address),
		transactors,
		m.Transactor.Key.GetAddress(),
	)
	logmsg := ""
	for _, transactor := range transactors {
		logmsg += transactor.String() + " "
	}
	log.Infoln("set transactors", logmsg)
	m.Transactor.AddTxMsg(setTransactorsMsg)
}
