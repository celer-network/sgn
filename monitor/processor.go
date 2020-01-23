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
	protobuf "github.com/golang/protobuf/proto"
)

func (m *EthMonitor) processQueue() {
	m.processPullerQueue()
	m.processEventQueue()
	m.processPusherQueue()
	m.processPenaltyQueue()
}

func (m *EthMonitor) processEventQueue() {
	secureBlockNum, err := m.getSecureBlockNum()
	if err != nil {
		log.Errorln("Query secureBlockNum err", err)
		return
	}

	iterator := m.db.Iterator(EventKeyPrefix, storetypes.PrefixEndBytes(EventKeyPrefix))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		event := NewEventFromBytes(iterator.Value())
		if secureBlockNum < event.Log.BlockNumber {
			continue
		}

		log.Infoln("Process mainchain event", event.Name, "at mainchain block", event.Log.BlockNumber)
		m.db.Delete(iterator.Key())

		switch e := event.ParseEvent(m.ethClient).(type) {
		case *mainchain.GuardInitializeCandidate:
			m.handleInitializeCandidate(e)
		case *mainchain.GuardDelegate:
			m.handleDelegate(e)
		case *mainchain.GuardValidatorChange:
			m.handleValidatorChange(e)
		case *mainchain.GuardIntendWithdraw:
			m.handleIntendWithdraw(e)
		}
	}
}

func (m *EthMonitor) processPullerQueue() {
	if !m.isPuller() {
		return
	}

	iterator := m.db.Iterator(PullerKeyPrefix, storetypes.PrefixEndBytes(PullerKeyPrefix))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		event := NewEventFromBytes(iterator.Value())
		log.Infoln("Process puller event", event.Name)
		m.db.Delete(iterator.Key())

		switch e := event.ParseEvent(m.ethClient).(type) {
		case *mainchain.GuardInitializeCandidate:
			m.processInitializeCandidate(e)
		}
	}
}

func (m *EthMonitor) processPusherQueue() {
	latestBlock, err := m.getLatestBlock()
	if err != nil {
		log.Errorln("Query latestBlock err", err)
		return
	}

	iterator := m.db.Iterator(PusherKeyPrefix, storetypes.PrefixEndBytes(PusherKeyPrefix))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		event := NewEventFromBytes(iterator.Value())
		log.Infoln("Process pusher event", event.Name)
		m.db.Delete(iterator.Key())

		switch e := event.ParseEvent(m.ethClient).(type) {
		case *mainchain.CelerLedgerIntendSettle:
			go m.processIntendSettle(e, latestBlock.Number)
		}
	}
}

func (m *EthMonitor) processPenaltyQueue() {
	iterator := m.db.Iterator(PenaltyKeyPrefix, storetypes.PrefixEndBytes(PenaltyKeyPrefix))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		event := NewPenaltyEventFromBytes(iterator.Value())
		m.db.Delete(iterator.Key())
		m.processPenalty(event)
	}
}

func (m *EthMonitor) processInitializeCandidate(initializeCandidate *mainchain.GuardInitializeCandidate) {
	log.Infof("Add InitializeCandidate of %x to transactor msgQueue", initializeCandidate.Candidate)

	msg := validator.NewMsgInitializeCandidate(
		mainchain.Addr2Hex(initializeCandidate.Candidate), m.transactor.Key.GetAddress())
	m.transactor.AddTxMsg(msg)
}

func (m *EthMonitor) processIntendSettle(intendSettle *mainchain.CelerLedgerIntendSettle, latestBlockNum uint64) {
	log.Infof("Process IntendSettle %x, tx hash %x", intendSettle.ChannelId, intendSettle.Raw.TxHash)
	channelId := intendSettle.ChannelId[:]
	addresses, seqNums, err := m.ethClient.Ledger.GetStateSeqNumMap(&bind.CallOpts{}, intendSettle.ChannelId)
	if err != nil {
		log.Errorln("Query StateSeqNumMap err", err)
		return
	}

	for i := 0; i < 2; i++ {
		peerFrom := mainchain.Addr2Hex(addresses[i])
		request, err := m.getRequest(channelId, mainchain.Addr2Hex(addresses[i]))
		if err != nil {
			log.Errorln("Query request err", err)
			return
		}

		if request.GuardTxHash != "" {
			log.Errorln("Request has been fulfilled")
			return
		}

		if seqNums[request.PeerFromIndex].Uint64() >= request.SeqNum {
			log.Infoln("Ignore the intendSettle event with a larger seqNum")
			return
		}

		if !m.isRequestGuard(request, latestBlockNum, intendSettle.Raw.BlockNumber) {
			log.Infof("Not valid guard at current mainchain block")
			event := NewEvent(IntendSettle, intendSettle.Raw)
			m.db.Set(GetPusherKey(intendSettle.Raw), event.MustMarshal())
			return
		}

		var signedSimplexState chain.SignedSimplexState
		err = protobuf.Unmarshal(request.SignedSimplexStateBytes, &signedSimplexState)
		if err != nil {
			log.Errorln("Unmarshal SignedSimplexState error:", err)
			return
		}

		signedSimplexStateArrayBytes, err := protobuf.Marshal(&chain.SignedSimplexStateArray{
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

		log.Infof("IntendSettle tx hash %x", tx.Hash())
		// TODO: 1) bockDelay, 2) may need a better way than wait mined,
		mainchain.WaitMined(context.Background(), m.ethClient.Client, tx, 2)

		log.Infof("Add MsgGuardProof %x to transactor msgQueue", tx.Hash())
		msg := subscribe.NewMsgGuardProof(channelId, peerFrom, intendSettle.Raw.TxHash.Hex(), tx.Hash().Hex(), m.transactor.Key.GetAddress())
		m.transactor.AddTxMsg(msg)
	}

}

func (m *EthMonitor) processPenalty(penaltyEvent PenaltyEvent) {
	log.Infoln("Process Penalty", penaltyEvent.nonce)

	used, err := m.ethClient.Guard.UsedPenaltyNonce(&bind.CallOpts{}, big.NewInt(int64(penaltyEvent.nonce)))
	if err != nil {
		log.Errorln("get usedPenaltyNonce err", err)
		return
	}

	if used {
		return
	}

	penaltyRequest, err := slash.CLIQueryPenaltyRequest(m.transactor.CliCtx, slash.StoreKey, penaltyEvent.nonce)
	if err != nil {
		log.Errorln("QueryPenaltyRequest err", err)
		return
	}

	tx, err := m.ethClient.Guard.Punish(m.ethClient.Auth, penaltyRequest)
	if err != nil {
		log.Errorln("Punish err", err)
		m.db.Set(GetPenaltyKey(penaltyEvent.nonce), penaltyEvent.MustMarshal())
		return
	}

	log.Infoln("Punish tx detail", tx)
}
