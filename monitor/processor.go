package monitor

import (
	"math/big"

	"github.com/celer-network/goutils/eth"
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
		err = m.db.Delete(iterator.Key())
		if err != nil {
			log.Errorln("db Delete err", err)
			continue
		}

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
		err = m.db.Delete(iterator.Key())
		if err != nil {
			log.Errorln("db Delete err", err)
			continue
		}
		m.submitPenalty(event)
	}
}

func (m *EthMonitor) syncConfirmParamProposal(confirmParamProposal *mainchain.DPoSConfirmParamProposal) {
	paramChange := common.NewParamChange(sdk.NewIntFromBigInt(confirmParamProposal.Record), sdk.NewIntFromBigInt(confirmParamProposal.NewValue))
	paramChangeData := m.operator.CliCtx.Codec.MustMarshalBinaryBare(paramChange)
	msg := sync.NewMsgSubmitChange(sync.ConfirmParamProposal, paramChangeData, m.operator.Key.GetAddress())
	log.Infof("submit change tx: confirm param proposal Record %v, NewValue %v", confirmParamProposal.Record, confirmParamProposal.NewValue)
	m.operator.AddTxMsg(msg)
}

func (m *EthMonitor) syncUpdateSidechainAddr(updateSidechainAddr *mainchain.SGNUpdateSidechainAddr) {
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
		err := m.db.Set(GetPusherKey(rawLog), event.MustMarshal())
		if err != nil {
			log.Errorln("db Set err", err)
		}
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
	var receipt *ethtypes.Receipt
	switch eventName {
	case IntendWithdrawChannel:
		receipt, err = m.ethClient.Transactor.TransactWaitMined(
			"SnapshotStates",
			func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*ethtypes.Transaction, error) {
				return m.ethClient.Ledger.SnapshotStates(opts, signedSimplexStateArrayBytes)
			})
	case IntendSettle:
		receipt, err = m.ethClient.Transactor.TransactWaitMined(
			"IntendSettle",
			func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*ethtypes.Transaction, error) {
				return m.ethClient.Ledger.IntendSettle(opts, signedSimplexStateArrayBytes)
			})
	default:
		log.Errorln("Invalid eventName", eventName)
		return
	}

	if err != nil {
		log.Errorln("intendSettle/snapshotStates err", err)
		return
	}

	txHash := receipt.TxHash
	log.Infof("Add MsgGuardProof %x to transactor msgQueue", txHash)
	request.GuardTxHash = txHash.Hex()
	request.GuardTxBlkNum = receipt.BlockNumber.Uint64()
	request.GuardSender = mainchain.Addr2Hex(m.ethClient.Address)
	requestData := m.operator.CliCtx.Codec.MustMarshalBinaryBare(request)
	msg := sync.NewMsgSubmitChange(sync.GuardProof, requestData, m.operator.Key.GetAddress())
	log.Infof("submit change tx: guard proof request %s", request)
	m.operator.AddTxMsg(msg)
}

func (m *EthMonitor) submitPenalty(penaltyEvent PenaltyEvent) {
	log.Infoln("Process Penalty", penaltyEvent.Nonce)

	used, err := m.ethClient.DPoS.UsedPenaltyNonce(&bind.CallOpts{
		BlockNumber: sdk.NewIntFromUint64(m.secureBlkNum).BigInt(),
	}, big.NewInt(int64(penaltyEvent.Nonce)))
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

	tx, err := m.ethClient.Transactor.Transact(
		&eth.TransactionStateHandler{
			OnMined: func(receipt *ethtypes.Receipt) {
				if receipt.Status == ethtypes.ReceiptStatusSuccessful {
					log.Infof("Punish transaction %x succeeded", receipt.TxHash)
				} else {
					log.Errorf("Punish transaction %x failed", receipt.TxHash)
				}
			},
			OnError: func(tx *ethtypes.Transaction, err error) {
				log.Errorf("Punish transaction %x err: %s", tx.Hash(), err)
			},
		},
		func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*ethtypes.Transaction, error) {
			return m.ethClient.DPoS.Punish(opts, penaltyRequest)
		},
	)
	if err != nil {
		if penaltyEvent.RetryCount < maxPunishRetry {
			penaltyEvent.RetryCount = penaltyEvent.RetryCount + 1
			err = m.db.Set(GetPenaltyKey(penaltyEvent.Nonce), penaltyEvent.MustMarshal())
			if err != nil {
				log.Errorln("db Set err", err)
			}
			return
		}
		log.Errorln("Punish err", err)
		return
	}
	log.Infoln("Punish tx submitted", tx.Hash().Hex())
}

func (m *EthMonitor) getRequests(cid [32]byte) (requests []subscribe.Request) {
	channelId := cid[:]
	addresses, seqNums, err := m.ethClient.Ledger.GetStateSeqNumMap(&bind.CallOpts{
		BlockNumber: sdk.NewIntFromUint64(m.secureBlkNum).BigInt(),
	}, cid)
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
