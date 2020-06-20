package monitor

import (
	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/x/subscribe"
	"github.com/celer-network/sgn/x/sync"
	"github.com/celer-network/sgn/x/validator"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
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
		if m.secureBlkNum < event.Log.BlockNumber {
			continue
		}

		log.Infoln("Process puller event", event.Name, "at mainchain block", event.Log.BlockNumber)
		err = m.db.Delete(iterator.Key())
		if err != nil {
			log.Errorln("db Delete err", err)
			continue
		}

		switch e := event.ParseEvent(m.ethClient).(type) {
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

func (m *Monitor) syncConfirmParamProposal(confirmParamProposal *mainchain.DPoSConfirmParamProposal) {
	paramChange := common.NewParamChange(sdk.NewIntFromBigInt(confirmParamProposal.Record), sdk.NewIntFromBigInt(confirmParamProposal.NewValue))
	paramChangeData := m.operator.CliCtx.Codec.MustMarshalBinaryBare(paramChange)
	msg := sync.NewMsgSubmitChange(sync.ConfirmParamProposal, paramChangeData, m.operator.Key.GetAddress())
	log.Infof("submit change tx: confirm param proposal Record %v, NewValue %v", confirmParamProposal.Record, confirmParamProposal.NewValue)
	m.operator.AddTxMsg(msg)
}

func (m *Monitor) syncUpdateSidechainAddr(updateSidechainAddr *mainchain.SGNUpdateSidechainAddr) {
	sidechainAddr, err := m.ethClient.SGN.SidechainAddrMap(&bind.CallOpts{
		BlockNumber: sdk.NewIntFromUint64(m.secureBlkNum).BigInt(),
	}, updateSidechainAddr.Candidate)
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

	disputeTimeout, err := m.ethClient.Ledger.GetDisputeTimeout(&bind.CallOpts{
		BlockNumber: sdk.NewIntFromUint64(m.secureBlkNum).BigInt(),
	}, mainchain.Bytes2Cid(request.ChannelId))
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
