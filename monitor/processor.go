package monitor

import (
	"context"
	"math/big"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/proto/chain"
	"github.com/celer-network/sgn/x/slash"
	"github.com/celer-network/sgn/x/subscribe"
	"github.com/celer-network/sgn/x/sync"
	"github.com/celer-network/sgn/x/validator"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/protobuf/proto"
	"github.com/spf13/viper"
)

const (
	maxPunishRetry = 5
)

func (m *EthMonitor) processQueue() {
	m.processEventQueue()
	m.processPullerQueue()
	m.processPusherQueue()
	m.processPenaltyQueue()
}

func (m *EthMonitor) processPullerQueue() {
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
		m.db.Delete(iterator.Key())

		switch e := event.ParseEvent(m.ethClient).(type) {
		case *mainchain.SGNUpdateSidechainAddr:
			m.syncUpdateSidechainAddr(e)
		case *mainchain.CelerLedgerIntendSettle:
			m.syncIntendSettle(e)
		case *mainchain.CelerLedgerIntendWithdraw:
			m.syncIntendWithdrawChannel(e)
		}
	}
}

func (m *EthMonitor) processPusherQueue() {
	iterator, err := m.db.Iterator(PusherKeyPrefix, storetypes.PrefixEndBytes(PusherKeyPrefix))
	if err != nil {
		log.Errorln("Create db iterator err", err)
		return
	}
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		event := NewEventFromBytes(iterator.Value())
		log.Infoln("Process pusher event", event.Name)
		m.db.Delete(iterator.Key())

		switch e := event.ParseEvent(m.ethClient).(type) {
		case *mainchain.CelerLedgerIntendSettle:
			m.guardIntendSettle(e)
		}
	}
}

func (m *EthMonitor) processPenaltyQueue() {
	if !m.isPusher() {
		return
	}

	iterator, err := m.db.Iterator(PenaltyKeyPrefix, storetypes.PrefixEndBytes(PenaltyKeyPrefix))
	if err != nil {
		log.Errorln("Create db iterator err", err)
		return
	}
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		event := NewPenaltyEventFromBytes(iterator.Value())
		m.db.Delete(iterator.Key())
		m.submitPenalty(event)
	}
}

// TODO: need to handle update after first initialization
func (m *EthMonitor) syncUpdateSidechainAddr(updateSidechainAddr *mainchain.SGNUpdateSidechainAddr) {
	_, err := validator.CLIQueryCandidate(m.operator.CliCtx, validator.RouterKey, mainchain.Addr2Hex(updateSidechainAddr.Candidate))
	if err == nil {
		log.Infof("The sidechain address of candidate %x has been updated", updateSidechainAddr.Candidate)
		return
	}

	log.Infof("Add UpdateSidechainAddr of %x to transactor msgQueue", updateSidechainAddr.Candidate)
	sidechainAddr, err := m.ethClient.SGN.SidechainAddrMap(&bind.CallOpts{}, updateSidechainAddr.Candidate)
	if err != nil {
		log.Errorln("Query sidechain ddress error:", err)
		return
	}

	candidate := validator.NewCandidate(updateSidechainAddr.Candidate.Hex(), sdk.AccAddress(sidechainAddr))
	candidateData := m.operator.CliCtx.Codec.MustMarshalBinaryBare(candidate)
	msg := sync.NewMsgSubmitChange(sync.UpdateSidechainAddr, candidateData, m.operator.Key.GetAddress())
	m.operator.AddTxMsg(msg)
}

func (m *EthMonitor) syncIntendSettle(intendSettle *mainchain.CelerLedgerIntendSettle) {
	log.Infof("Sync IntendSettle %x, tx hash %x", intendSettle.ChannelId, intendSettle.Raw.TxHash)
	requests := m.getRequests(intendSettle.ChannelId)
	for _, request := range requests {
		m.triggerGuard(request, intendSettle.Raw)
	}
}

func (m *EthMonitor) syncIntendWithdrawChannel(intendWithdrawChannel *mainchain.CelerLedgerIntendWithdraw) {
	log.Infof("Sync intendWithdrawChannel %x, tx hash %x", intendWithdrawChannel.ChannelId, intendWithdrawChannel.Raw.TxHash)
	requests := m.getRequests(intendWithdrawChannel.ChannelId)
	for _, request := range requests {
		m.triggerGuard(request, intendWithdrawChannel.Raw)
	}
}

func (m *EthMonitor) triggerGuard(request subscribe.Request, rawLog ethtypes.Log) {
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
	msg := sync.NewMsgSubmitChange(sync.IntendSettle, requestData, m.operator.Key.GetAddress())
	m.operator.AddTxMsg(msg)
}

func (m *EthMonitor) guardIntendSettle(intendSettle *mainchain.CelerLedgerIntendSettle) {
	log.Infof("Guard IntendSettle %x, tx hash %x", intendSettle.ChannelId, intendSettle.Raw.TxHash)
	requests := m.getRequests(intendSettle.ChannelId)
	for _, request := range requests {
		m.guardRequest(request, intendSettle.Raw, IntendSettle)
	}
}

func (m *EthMonitor) guardIntendWithdrawChannel(intendWithdrawChannel *mainchain.CelerLedgerIntendWithdraw) {
	log.Infof("Guard intendWithdrawChannel %x, tx hash %x", intendWithdrawChannel.ChannelId, intendWithdrawChannel.Raw.TxHash)
	requests := m.getRequests(intendWithdrawChannel.ChannelId)
	for _, request := range requests {
		m.guardRequest(request, intendWithdrawChannel.Raw, IntendWithdrawChannel)
	}
}

func (m *EthMonitor) guardRequest(request subscribe.Request, rawLog ethtypes.Log, eventName EventName) {
	log.Infoln("Guard request", request)
	if request.GuardTxHash != "" {
		log.Errorln("Request has been fulfilled")
		return
	}

	if !m.isRequestGuard(request, rawLog.BlockNumber) {
		log.Infof("Not valid guard at current mainchain block")
		event := NewEvent(IntendSettle, rawLog)
		m.db.Set(GetPusherKey(rawLog), event.MustMarshal())
		return
	}

	var signedSimplexState chain.SignedSimplexState
	err := proto.Unmarshal(request.SignedSimplexStateBytes, &signedSimplexState)
	if err != nil {
		log.Errorln("Unmarshal SignedSimplexState error:", err)
		return
	}

	signedSimplexStateArrayBytes, err := proto.Marshal(&chain.SignedSimplexStateArray{
		SignedSimplexStates: []*chain.SignedSimplexState{&signedSimplexState},
	})
	if err != nil {
		log.Errorln("Marshal signedSimplexStateArrayBytes error:", err)
		return
	}

	// TODO: use snapshotStates instead of intendSettle here? (need to update cChannel contract first)
	var tx *ethtypes.Transaction
	switch eventName {
	case IntendWithdrawChannel:
		tx, err = m.ethClient.Ledger.SnapshotStates(m.ethClient.Auth, signedSimplexStateArrayBytes)
	case IntendSettle:
		tx, err = m.ethClient.Ledger.IntendSettle(m.ethClient.Auth, signedSimplexStateArrayBytes)
	default:
		log.Errorln("Invalid eventName", eventName)
		return
	}

	if err != nil {
		log.Errorln("intendSettle/snapshotStates err", err)
		return
	}

	// TODO: 1) bockDelay, 2) may need a better way than wait mined,
	res, err := mainchain.WaitMined(context.Background(), m.ethClient.Client, tx, viper.GetUint64(common.FlagEthConfirmCount))
	if err != nil {
		log.Errorln("intendSettle WaitMined err", err, tx.Hash().Hex())
		return
	}
	if res.Status != ethtypes.ReceiptStatusSuccessful {
		log.Errorln("intendSettle failed", tx.Hash().Hex())
		return
	}

	log.Infof("Add MsgGuardProof %x to transactor msgQueue", tx.Hash())
	request.GuardTxHash = tx.Hash().Hex()
	request.GuardTxBlkNum = res.BlockNumber.Uint64()
	request.GuardSender = mainchain.Addr2Hex(m.ethClient.Address)
	requestData := m.operator.CliCtx.Codec.MustMarshalBinaryBare(request)
	msg := sync.NewMsgSubmitChange(sync.GuardProof, requestData, m.operator.Key.GetAddress())
	m.operator.AddTxMsg(msg)
}

func (m *EthMonitor) submitPenalty(penaltyEvent PenaltyEvent) {
	log.Infoln("Process Penalty", penaltyEvent.Nonce)

	used, err := m.ethClient.DPoS.UsedPenaltyNonce(&bind.CallOpts{}, big.NewInt(int64(penaltyEvent.Nonce)))
	if err != nil {
		log.Errorln("Get usedPenaltyNonce err", err)
		return
	}

	if used {
		log.Infof("Penalty %d has been used", penaltyEvent.Nonce)
		return
	}

	penaltyRequest, err := slash.CLIQueryPenaltyRequest(m.operator.CliCtx, slash.StoreKey, penaltyEvent.Nonce)
	if err != nil {
		log.Errorln("QueryPenaltyRequest err", err)
		return
	}

	tx, err := m.ethClient.DPoS.Punish(m.ethClient.Auth, penaltyRequest)
	if err != nil {
		if penaltyEvent.RetryCount < maxPunishRetry {
			penaltyEvent.RetryCount = penaltyEvent.RetryCount + 1
			m.db.Set(GetPenaltyKey(penaltyEvent.Nonce), penaltyEvent.MustMarshal())
			return
		}
		log.Errorln("Punish err", err)
		return
	}
	log.Infoln("Punish tx submitted", tx.Hash().Hex())

	go m.waitPunishMined(tx)
}

func (m *EthMonitor) waitPunishMined(tx *ethtypes.Transaction) {
	// TODO: blockdelay
	res, err := mainchain.WaitMined(context.Background(), m.ethClient.Client, tx, 2)
	if err != nil {
		log.Errorln("Punish tx WaitMined err", err, tx.Hash().Hex())
		return
	}
	if res.Status != ethtypes.ReceiptStatusSuccessful {
		log.Errorln("Punish tx failed", tx.Hash().Hex())
		return
	}
	log.Infoln("Punish tx mined", tx.Hash().Hex())
}

func (m *EthMonitor) getRequests(cid [32]byte) (requests []subscribe.Request) {
	channelId := cid[:]
	addresses, seqNums, err := m.ethClient.Ledger.GetStateSeqNumMap(&bind.CallOpts{}, cid)
	if err != nil {
		log.Errorln("Query StateSeqNumMap err", err)
		return
	}

	for _, addr := range addresses {
		peerFrom := mainchain.Addr2Hex(addr)
		request, err := m.getRequest(channelId, peerFrom)
		if err != nil {
			continue
		}

		if seqNums[request.PeerFromIndex].Uint64() >= request.SeqNum {
			log.Infoln("Ignore the intendSettle event with an equal or larger seqNum")
			continue
		}

		requests = append(requests, request)
	}

	return requests
}
