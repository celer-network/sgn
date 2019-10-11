package monitor

import (
	"log"
	"math/big"

	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/x/slash"
	"github.com/celer-network/sgn/x/subscribe"
	"github.com/celer-network/sgn/x/validator"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func (m *EthMonitor) processQueue() {
	m.processPullerQueue()
	m.processEventQueue()
	m.processPusherQueue()
}

func (m *EthMonitor) processEventQueue() {
	secureBlockNum, err := m.getSecureBlockNum()
	if err != nil {
		log.Printf("Query secureBlockNum err", err)
		return
	}

	for m.eventQueue.Len() > 0 {
		e := m.eventQueue.Front().(Event)
		if secureBlockNum < e.log.BlockNumber {
			return
		}

		switch event := e.event.(type) {
		case *mainchain.GuardInitializeCandidate:
			m.handleInitializeCandidate(event)
		case *mainchain.GuardDelegate:
			m.handleDelegate(event)
		case *mainchain.GuardValidatorChange:
			m.handleValidatorChange(event)
		case *mainchain.GuardIntendWithdraw:
			m.handleIntendWithdraw(event)
		case *mainchain.CelerLedgerIntendSettle:
			m.handleIntendSettle(event)
		}

		m.eventQueue.PopFront()
	}
}

func (m *EthMonitor) processPullerQueue() {
	if !m.isPuller() {
		return
	}

	for m.pullerQueue.Len() != 0 {
		switch event := m.pullerQueue.PopFront().(type) {
		case *mainchain.GuardInitializeCandidate:
			m.processInitializeCandidate(event)
		}
	}
}

func (m *EthMonitor) processPusherQueue() {
	latestBlock, err := m.getLatestBlock()
	if err != nil {
		log.Printf("Query latestBlock err", err)
		return
	}

	for pusherLen := m.pusherQueue.Len(); pusherLen > 0; pusherLen-- {
		switch event := m.pusherQueue.PopFront().(type) {
		case *mainchain.CelerLedgerIntendSettle:
			m.processIntendSettle(event, latestBlock.Number)
		case PenaltyEvent:
			m.processPenalty(event)

		}
	}
}

func (m *EthMonitor) processInitializeCandidate(initializeCandidate *mainchain.GuardInitializeCandidate) {
	log.Printf("Process InitializeCandidate", initializeCandidate.Candidate)

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

	if request.TxHash != "" {
		log.Printf("Request has been fullfilled")
		return
	}

	if !m.isRequestGuard(request, latestBlockNum, intendSettle.Raw.BlockNumber) {
		m.pusherQueue.PushBack(intendSettle)
		return
	}

	tx, err := m.ethClient.Ledger.IntendSettle(m.ethClient.Auth, request.SignedSimplexStateBytes)
	if err != nil {
		log.Printf("intendSettle err", err)
		return
	}
	log.Printf("IntendSettle tx detail", tx)

	msg := subscribe.NewMsgGuardProof(channelId, tx.Hash().Hex(), m.transactor.Key.GetAddress())
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

	penaltyRequest, err := slash.CLIQueryPenaltyRequest(m.cdc, m.transactor.CliCtx, slash.StoreKey, penaltyEvent.nonce)
	if err != nil {
		log.Printf("QueryPenaltyRequest err", err)
		return
	}

	tx, err := m.ethClient.Guard.Punish(m.ethClient.Auth, penaltyRequest)
	if err != nil {
		log.Printf("Punish err", err)

		if !penaltyEvent.retry {
			penaltyEvent.retry = true
			m.pusherQueue.PushBack(penaltyEvent)
		}
		return
	}

	log.Printf("Punish tx detail", tx)
}
