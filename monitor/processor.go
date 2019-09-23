package monitor

import (
	"log"

	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/x/subscribe"
	"github.com/celer-network/sgn/x/validator"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func (m *EthMonitor) processQueue() {
	m.processPullerQueue()

	latestBlock, err := m.getLatestBlock()
	if err != nil {
		log.Printf("Query latestBlock err", err)
		return
	}

	m.processEventQueue(latestBlock.Number)
	m.processPusherQueue(latestBlock.Number)
}

func (m *EthMonitor) processEventQueue(latestBlockNum uint64) {
	for m.eventQueue.Len() > 0 {
		e := m.eventQueue.Front().(Event)
		if latestBlockNum < e.log.BlockNumber+common.ConfirmationCount {
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

func (m *EthMonitor) processPusherQueue(latestBlockNum uint64) {
	for pusherLen := m.pusherQueue.Len(); pusherLen > 0; pusherLen-- {
		switch event := m.pusherQueue.PopFront().(type) {
		case *mainchain.CelerLedgerIntendSettle:
			m.processIntendSettle(event, latestBlockNum)
		}
	}
}

func (m *EthMonitor) processInitializeCandidate(initializeCandidate *mainchain.GuardInitializeCandidate) {
	log.Printf("Process InitializeCandidate", initializeCandidate.Candidate)

	candidateInfo, err := m.ethClient.Guard.GetCandidateInfo(&bind.CallOpts{}, initializeCandidate.Candidate)
	if err != nil {
		log.Printf("Query candidate info err", err)
		return
	}

	_, err = m.getAccount(candidateInfo.SidechainAddr)
	if err == nil {
		log.Printf("Skip initilization for existing account")
		return
	}

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

	if !m.isRequestHandler(request, latestBlockNum, intendSettle.Raw.BlockNumber) {
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
