package monitor

import (
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

		log.Infoln("process mainchain event", event.Name, "at mainchain block", event.Log.BlockNumber)
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
		log.Infoln("process puller event", event.Name)
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
		log.Infoln("process pusher event", event.Name)
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
		log.Infoln("process penalty event", event.nonce)
		m.db.Delete(iterator.Key())
		m.processPenalty(event)
	}
}

func (m *EthMonitor) processInitializeCandidate(initializeCandidate *mainchain.GuardInitializeCandidate) {
	log.Infof("Push initializeCandidate of %s to transactor msgQueue", initializeCandidate.Candidate.String())

	msg := validator.NewMsgInitializeCandidate(initializeCandidate.Candidate.String(), m.transactor.Key.GetAddress())
	m.transactor.BroadcastTx(msg)
}

func (m *EthMonitor) processIntendSettle(intendSettle *mainchain.CelerLedgerIntendSettle, latestBlockNum uint64) {
	log.Infof("Process IntendSettle %x", intendSettle.ChannelId)
	channelId := intendSettle.ChannelId[:]
	request, err := m.getRequest(channelId)
	if err != nil {
		log.Errorln("Query request err", err)
		return
	}

	if request.GuardTxHash != "" {
		log.Errorln("Request has been fulfilled")
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
	log.Infof("IntendSettle tx detail %+v", tx)

	triggerTxHash := intendSettle.Raw.TxHash.Hex()
	msg := subscribe.NewMsgGuardProof(channelId, triggerTxHash, tx.Hash().Hex(), m.transactor.Key.GetAddress())
	m.transactor.BroadcastTx(msg)
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
