package monitor

import (
	"fmt"
	"strings"
	"time"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn-contract/bindings/go/sgncontracts"
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
	if !m.isSyncer() {
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

	validators := make(map[mainchain.Addr]bool)
	delegators := make(map[string]bool)
	for i, key := range keys {
		event := NewEventFromBytes(vals[i])
		logmsg := fmt.Sprintf("Process puller event %s at mainchain block %d", event.Name, event.Log.BlockNumber)
		err = m.dbDelete(key)
		if err != nil {
			log.Errorf("%s. db Delete err: %s", logmsg, err)
			continue
		}

		switch e := event.ParseEvent(m.EthClient).(type) {
		case *sgncontracts.DPoSValidatorChange:
			log.Infof("%s. validator change %x type %d", logmsg, e.EthAddr, e.ChangeType)
			validators[e.EthAddr] = true

		case *sgncontracts.DPoSUpdateDelegatedStake:
			log.Infof("%s. stake update delegator %x, candidate %x, stake %s, pool %s",
				logmsg, e.Delegator, e.Candidate, e.DelegatorStake, e.CandidatePool)
			validators[e.Candidate] = true
			delegators[getDelegatorKey(e.Candidate, e.Delegator)] = true

		case *sgncontracts.DPoSCandidateUnbonded:
			log.Infof("%s. candidate unbonded %x", logmsg, e.Candidate)
			validators[e.Candidate] = true

		case *sgncontracts.DPoSUpdateCommissionRate:
			log.Infof("%s. commission update %x, %s", logmsg, e.Candidate, e.NewRate)
			validators[e.Candidate] = true

		case *sgncontracts.SGNUpdateSidechainAddr:
			m.syncUpdateSidechainAddr(e, logmsg)

		case *sgncontracts.DPoSConfirmParamProposal:
			m.syncConfirmParamProposal(e, logmsg)

		case *sgncontracts.SGNAddSubscriptionBalance:
			m.syncSubscriptionBalance(e, logmsg)

		case *mainchain.CelerLedgerIntendSettle:
			e.Raw = event.Log
			m.syncIntendSettle(e, logmsg)

		case *mainchain.CelerLedgerIntendWithdraw:
			e.Raw = event.Log
			m.syncIntendWithdrawChannel(e, logmsg)
		}
	}

	if m.isBootstrapped() {
		for validatorAddr := range validators {
			m.SyncValidator(validatorAddr)
		}
	}
	for delegatorKey := range delegators {
		candidatorAddr := mainchain.Hex2Addr(strings.Split(delegatorKey, ":")[0])
		delegatorAddr := mainchain.Hex2Addr(strings.Split(delegatorKey, ":")[1])
		m.SyncDelegator(candidatorAddr, delegatorAddr)
	}
}

func getDelegatorKey(candidate, delegator mainchain.Addr) string {
	return mainchain.Addr2Hex(candidate) + ":" + mainchain.Addr2Hex(delegator)
}

func (m *Monitor) syncBlkNum() {
	for {
		time.Sleep(time.Minute)
		if !m.isSyncer() {
			continue
		}

		msg := sync.NewMsgSubmitChange(sync.SyncBlkNum, []byte{0}, m.EthClient.Client, m.Transactor.Key.GetAddress())
		log.Infof("submit change tx: sync maichain block number", msg.BlockNum)
		m.Transactor.AddTxMsg(msg)
	}
}

func (m *Monitor) syncUpdateSidechainAddr(sidechainAddr *sgncontracts.SGNUpdateSidechainAddr, logmsg string) {
	log.Infof("%s. sidechainAddr update %x, %s",
		logmsg, sidechainAddr.Candidate, sdk.AccAddress(sidechainAddr.NewSidechainAddr.Bytes()))
	m.SyncUpdateSidechainAddr(sidechainAddr.Candidate)
}

func (m *Monitor) syncConfirmParamProposal(confirmParamProposal *sgncontracts.DPoSConfirmParamProposal, logmsg string) {
	paramChange := common.NewParamChange(sdk.NewIntFromBigInt(confirmParamProposal.Record), sdk.NewIntFromBigInt(confirmParamProposal.NewValue))
	paramChangeData := m.Transactor.CliCtx.Codec.MustMarshalBinaryBare(paramChange)
	msg := sync.NewMsgSubmitChange(sync.ConfirmParamProposal, paramChangeData, m.EthClient.Client, m.Transactor.Key.GetAddress())
	log.Infof("%s. submit change tx: confirm param proposal Record %v, NewValue %v", logmsg, confirmParamProposal.Record, confirmParamProposal.NewValue)
	m.Transactor.AddTxMsg(msg)
}

func (m *Monitor) syncSubscriptionBalance(event *sgncontracts.SGNAddSubscriptionBalance, logmsg string) {
	m.SyncSubscriptionBalance(event.Consumer)
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
		if intendWithdrawChannel.Receiver == mainchain.Hex2Addr(request.SimplexSender) {
			m.triggerGuard(request, intendWithdrawChannel.Raw, seqs[i], guard.ChanStatus_Withdrawing)
		}
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
	syncData := m.Transactor.CliCtx.Codec.MustMarshalBinaryBare(trigger)
	msg := sync.NewMsgSubmitChange(sync.GuardTrigger, syncData, m.EthClient.Client, m.Transactor.Key.GetAddress())
	log.Infof("submit change tx: trigger guard request %s", trigger)
	m.Transactor.AddTxMsg(msg)
}

func (m *Monitor) setTransactors() {
	transactors, err := common.ParseTransactorAddrs(viper.GetStringSlice(common.FlagSgnTransactors))
	if err != nil {
		log.Errorln("parse transactors err", err)
		return
	}
	if len(transactors) == 0 {
		return
	}
	setTransactorsMsg := validator.NewMsgSetTransactors(
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
