package monitor

import (
	"context"
	"math/big"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/proto/chain"
	"github.com/celer-network/sgn/x/slash"
	"github.com/celer-network/sgn/x/subscribe"
	"github.com/celer-network/sgn/x/validator"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/protobuf/proto"
)

func (m *EthMonitor) processQueue() {
	secureBlockNum, err := m.getSecureBlockNum()
	if err != nil {
		// Retry once
		log.Errorln("Query secureBlockNum err", err)
		return
	}

	m.processEventQueue(secureBlockNum)
	m.processPullerQueue(secureBlockNum)
	m.processPusherQueue()
	m.processPenaltyQueue()
}

func (m *EthMonitor) processPullerQueue(secureBlockNum uint64) {
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
		if secureBlockNum < event.Log.BlockNumber {
			continue
		}

		log.Infoln("Process puller event", event.Name)
		m.db.Delete(iterator.Key())

		switch e := event.ParseEvent(m.ethClient).(type) {
		case *mainchain.GuardInitializeCandidate:
			m.syncInitializeCandidate(e)
		case *mainchain.CelerLedgerIntendSettle:
			m.syncIntendSettle(e)
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

func (m *EthMonitor) syncInitializeCandidate(initializeCandidate *mainchain.GuardInitializeCandidate) {
	_, err := validator.CLIQueryCandidate(m.operator.CliCtx, validator.RouterKey, mainchain.Addr2Hex(initializeCandidate.Candidate))
	if err == nil {
		log.Infof("Candidate %x has been initialized", initializeCandidate.Candidate)
		return
	}

	log.Infof("Add InitializeCandidate of %x to transactor msgQueue", initializeCandidate.Candidate)
	msg := validator.NewMsgInitializeCandidate(
		mainchain.Addr2Hex(initializeCandidate.Candidate), m.operator.Key.GetAddress())
	m.operator.AddTxMsg(msg)
}

func (m *EthMonitor) syncIntendSettle(intendSettle *mainchain.CelerLedgerIntendSettle) {
	log.Infof("Sync IntendSettle %x, tx hash %x", intendSettle.ChannelId, intendSettle.Raw.TxHash)
	requests := m.processIntendSettle(intendSettle)
	for _, request := range requests {
		if request.TriggerTxHash != "" {
			log.Infoln("The intendSettle event has been synced on sgn")
			return
		}

		msg := subscribe.NewMsgIntendSettle(request.ChannelId, request.GetPeerAddress(), intendSettle.Raw.TxHash.Hex(), m.operator.Key.GetAddress())
		m.operator.AddTxMsg(msg)
	}
}

func (m *EthMonitor) guardIntendSettle(intendSettle *mainchain.CelerLedgerIntendSettle) {
	log.Infof("Guard IntendSettle %x, tx hash %x", intendSettle.ChannelId, intendSettle.Raw.TxHash)
	requests := m.processIntendSettle(intendSettle)
	for _, request := range requests {
		if request.GuardTxHash != "" {
			log.Errorln("Request has been fulfilled")
			return
		}

		if !m.isRequestGuard(request, intendSettle.Raw.BlockNumber) {
			log.Infof("Not valid guard at current mainchain block")
			event := NewEvent(IntendSettle, intendSettle.Raw)
			m.db.Set(GetPusherKey(intendSettle.Raw), event.MustMarshal())
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
		tx, err := m.ethClient.Ledger.IntendSettle(m.ethClient.Auth, signedSimplexStateArrayBytes)
		if err != nil {
			log.Errorln("intendSettle err", err)
			return
		}

		// TODO: 1) bockDelay, 2) may need a better way than wait mined,
		res, err := mainchain.WaitMined(context.Background(), m.ethClient.Client, tx, 2)
		if err != nil {
			log.Errorln("intendSettle WaitMined err", err, tx.Hash().Hex())
			return
		}
		if res.Status != ethtypes.ReceiptStatusSuccessful {
			log.Errorln("intendSettle failed", tx.Hash().Hex())
			return
		}

		log.Infof("Add MsgGuardProof %x to transactor msgQueue", tx.Hash())
		msg := subscribe.NewMsgGuardProof(request.ChannelId, request.GetPeerAddress(), tx.Hash().Hex(), m.operator.Key.GetAddress())
		m.operator.AddTxMsg(msg)
	}
}

func (m *EthMonitor) submitPenalty(penaltyEvent PenaltyEvent) {
	log.Infoln("Process Penalty", penaltyEvent.Nonce)

	used, err := m.ethClient.Guard.UsedPenaltyNonce(&bind.CallOpts{}, big.NewInt(int64(penaltyEvent.Nonce)))
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

	tx, err := m.ethClient.Guard.Punish(m.ethClient.Auth, penaltyRequest)
	if err != nil {
		log.Errorln("Punish err", err)
		m.db.Set(GetPenaltyKey(penaltyEvent.Nonce), penaltyEvent.MustMarshal())
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

func (m *EthMonitor) processIntendSettle(intendSettle *mainchain.CelerLedgerIntendSettle) (requests []subscribe.Request) {
	channelId := intendSettle.ChannelId[:]
	addresses, seqNums, err := m.ethClient.Ledger.GetStateSeqNumMap(&bind.CallOpts{}, intendSettle.ChannelId)
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
