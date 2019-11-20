package monitor

import (
	"log"
	"math/big"

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
		log.Printf("Query secureBlockNum err", err)
		return
	}

	iterator := m.db.Iterator(EventKeyPrefix, storetypes.PrefixEndBytes(EventKeyPrefix))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		event := NewEventFromBytes(iterator.Value())
		if secureBlockNum < event.Log.BlockNumber {
			continue
		}

		log.Printf("process mainchain event", event.Name, event.Log.BlockNumber)
		m.db.Delete(iterator.Key())

		switch e := event.ParseEvent(m.ethClient).(type) {
		case *mainchain.GuardInitializeCandidate:
			e.Raw = event.Log
			m.handleInitializeCandidate(e)
		case *mainchain.GuardDelegate:
			m.handleDelegate(e)
		case *mainchain.GuardValidatorChange:
			m.handleValidatorChange(e)
		case *mainchain.GuardIntendWithdraw:
			m.handleIntendWithdraw(e)
		case *mainchain.CelerLedgerIntendSettle:
			m.handleIntendSettle(e)
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
		log.Printf("process puller event", event.Name)
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
		log.Printf("Query latestBlock err", err)
		return
	}

	iterator := m.db.Iterator(PusherKeyPrefix, storetypes.PrefixEndBytes(PusherKeyPrefix))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		event := NewEventFromBytes(iterator.Value())
		log.Printf("process pusher event", event.Name)
		m.db.Delete(iterator.Key())

		switch e := event.ParseEvent(m.ethClient).(type) {
		case *mainchain.CelerLedgerIntendSettle:
			m.processIntendSettle(e, latestBlock.Number)
		}
	}
}

func (m *EthMonitor) processPenaltyQueue() {
	iterator := m.db.Iterator(PenaltyKeyPrefix, storetypes.PrefixEndBytes(PenaltyKeyPrefix))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		event := NewPenaltyEventFromBytes(iterator.Value())
		log.Printf("process penalty event", event.nonce)
		m.db.Delete(iterator.Key())
		m.processPenalty(event)
	}
}

func (m *EthMonitor) processInitializeCandidate(initializeCandidate *mainchain.GuardInitializeCandidate) {
	log.Printf("Push MsgInitializeCandidate of %s to Transactor's msgQueue for broadcast", initializeCandidate.Candidate.String())

	msg := validator.NewMsgInitializeCandidate(initializeCandidate.Candidate.String(), m.transactor.Key.GetAddress())
	m.transactor.BroadcastTx(msg)
}

func (m *EthMonitor) processIntendSettle(intendSettle *mainchain.CelerLedgerIntendSettle, latestBlockNum uint64) {
	log.Printf("Process IntendSettle", intendSettle.ChannelId)
	channelId := intendSettle.ChannelId[:]
	request, err := m.getRequest(channelId)
	if err != nil {
		log.Printf("Query request err", err)
		return
	}

	if request.GuardTxHash != "" {
		log.Printf("Request has been fulfilled")
		return
	}

	if !m.isRequestGuard(request, latestBlockNum, intendSettle.Raw.BlockNumber) {
		event := NewEvent(IntendSettle, intendSettle.Raw)
		m.db.Set(GetPusherKey(intendSettle.Raw), event.MustMarshal())
		return
	}

	var signedSimplexState chain.SignedSimplexState
	err = protobuf.Unmarshal(request.SignedSimplexStateBytes, &signedSimplexState)
	if err != nil {
		log.Print("Unmarshal SignedSimplexState error: ", err)
		return
	}
	signedSimplexStateArrayBytes, err := protobuf.Marshal(&chain.SignedSimplexStateArray{
		SignedSimplexStates: []*chain.SignedSimplexState{&signedSimplexState},
	})
	if err != nil {
		log.Print("Marshal signedSimplexStateArrayBytes error: ", err)
		return
	}
	// TODO: use snapshotStates instead of intendSettle here? (need to update cChannel contract first)
	tx, err := m.ethClient.Ledger.IntendSettle(m.ethClient.Auth, signedSimplexStateArrayBytes)
	if err != nil {
		log.Printf("intendSettle err", err)
		return
	}
	log.Printf("IntendSettle tx detail", tx)

	triggerTxHash := intendSettle.Raw.TxHash.Hex()
	msg := subscribe.NewMsgGuardProof(channelId, triggerTxHash, tx.Hash().Hex(), m.transactor.Key.GetAddress())
	m.transactor.BroadcastTx(msg)
}

func (m *EthMonitor) processPenalty(penaltyEvent PenaltyEvent) {
	log.Printf("Process Penalty", penaltyEvent.nonce)

	used, err := m.ethClient.Guard.UsedPenaltyNonce(&bind.CallOpts{}, big.NewInt(int64(penaltyEvent.nonce)))
	if err != nil {
		log.Printf("get usedPenaltyNonce err", err)
		return
	}

	if used {
		return
	}

	penaltyRequest, err := slash.CLIQueryPenaltyRequest(m.transactor.CliCtx, slash.StoreKey, penaltyEvent.nonce)
	if err != nil {
		log.Printf("QueryPenaltyRequest err", err)
		return
	}

	tx, err := m.ethClient.Guard.Punish(m.ethClient.Auth, penaltyRequest)
	if err != nil {
		log.Printf("Punish err", err)
		m.db.Set(GetPenaltyKey(penaltyEvent.nonce), penaltyEvent.MustMarshal())
		return
	}

	log.Printf("Punish tx detail", tx)
}
