package monitor

import (
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
		log.Infoln("Process puller event", event.Name, "at mainchain block", event.Log.BlockNumber)
		err = m.dbDelete(key)
		if err != nil {
			log.Errorln("db Delete err", err)
			continue
		}

		switch e := event.ParseEvent(m.EthClient).(type) {
		case *mainchain.DPoSValidatorChange:
			m.syncDPoSValidatorChange(e)
		case *mainchain.DPoSIntendWithdraw:
			m.syncDPoSIntendWithdraw(e)
		case *mainchain.DPoSCandidateUnbonded:
			m.syncDPoSCandidateUnbonded(e)
		case *mainchain.DPoSConfirmParamProposal:
			m.syncConfirmParamProposal(e)
		case *mainchain.DPoSUpdateCommissionRate:
			m.syncUpdateCommissionRate(e)
		case *mainchain.SGNUpdateSidechainAddr:
			m.SyncUpdateSidechainAddr(e.Candidate)
		case *mainchain.SGNAddSubscriptionBalance:
			m.syncSGNAddSubscriptionBalance(e)
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
	m.SyncValidator(validatorChange.EthAddr)
}

func (m *Monitor) syncDPoSIntendWithdraw(intendWithdraw *mainchain.DPoSIntendWithdraw) {
	log.Infof("puller queue process intend withdraw %x", intendWithdraw.Candidate)
	m.SyncValidator(intendWithdraw.Candidate)
	m.SyncDelegator(intendWithdraw.Candidate, intendWithdraw.Delegator)
}

func (m *Monitor) syncDPoSCandidateUnbonded(candidateUnbonded *mainchain.DPoSCandidateUnbonded) {
	log.Infof("puller queue process candidate unbonded %x", candidateUnbonded.Candidate)
	m.SyncValidator(candidateUnbonded.Candidate)
}

func (m *Monitor) syncUpdateCommissionRate(commission *mainchain.DPoSUpdateCommissionRate) {
	log.Infof("puller queue process commission update %x, %s", commission.Candidate, commission.NewRate)
	m.SyncValidator(commission.Candidate)
}

func (m *Monitor) syncConfirmParamProposal(confirmParamProposal *mainchain.DPoSConfirmParamProposal) {
	paramChange := common.NewParamChange(sdk.NewIntFromBigInt(confirmParamProposal.Record), sdk.NewIntFromBigInt(confirmParamProposal.NewValue))
	paramChangeData := m.Transactor.CliCtx.Codec.MustMarshalBinaryBare(paramChange)
	msg := sync.NewMsgSubmitChange(sync.ConfirmParamProposal, paramChangeData, m.Transactor.Key.GetAddress())
	log.Infof("submit change tx: confirm param proposal Record %v, NewValue %v", confirmParamProposal.Record, confirmParamProposal.NewValue)
	m.Transactor.AddTxMsg(msg)
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

func (m *Monitor) syncSGNAddSubscriptionBalance(event *mainchain.SGNAddSubscriptionBalance) {
	transactor := m.Transactor
	consumer := event.Consumer
	consumerEthAddress := consumer.Hex()
	amount := event.Amount
	amountInt := sdk.NewIntFromBigInt(amount)
	subscription, err := guard.CLIQuerySubscription(transactor.CliCtx, guard.RouterKey, consumerEthAddress)
	if err == nil {
		if subscription.Deposit.Equal(amountInt) {
			log.Infof("Subscription already updated for %s, amount %s", consumerEthAddress, amount)
			return
		}
	}
	subscription = guard.NewSubscription(consumer.Hex())
	subscription.Deposit = sdk.NewIntFromBigInt(amount)
	subscriptionData := transactor.CliCtx.Codec.MustMarshalBinaryBare(subscription)
	msg := sync.NewMsgSubmitChange(sync.Subscribe, subscriptionData, m.Transactor.Key.GetAddress())
	log.Infof("Submit change tx: subscribe ethAddress %s, amount %s, mainchain tx hash %x", consumerEthAddress, amount, event.Raw.TxHash)
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
	msg := sync.NewMsgSubmitChange(sync.GuardTrigger, syncData, m.Transactor.Key.GetAddress())
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
