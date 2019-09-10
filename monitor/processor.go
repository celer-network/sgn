package monitor

import (
	"log"

	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/x/subscribe"
)

func (m *EthMonitor) processQueue() {
	m.processEventQueue()
	m.processIntendSettleQueue()
}

func (m *EthMonitor) processEventQueue() {
	latestBlock, err := m.getLatestBlock()
	if err != nil {
		log.Printf("Query latestBlock err", err)
		return
	}

	for m.eventQueue.Len() > 0 {
		e := m.eventQueue.Front().(Event)
		if latestBlock.Number < e.log.BlockNumber+common.ConfirmationCount {
			return
		}

		switch event := e.event.(type) {
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

func (m *EthMonitor) processIntendSettleQueue() {
	if !m.isPusher() {
		return
	}

	for m.intendSettleQueue.Len() != 0 {
		m.processIntendSettle(m.intendSettleQueue.PopFront().(*mainchain.CelerLedgerIntendSettle))
	}
}

func (m *EthMonitor) processIntendSettle(intendSettle *mainchain.CelerLedgerIntendSettle) {
	log.Printf("Process intend settle", intendSettle.ChannelId)
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

	tx, err := m.ethClient.Ledger.IntendSettle(m.ethClient.Auth, request.SignedSimplexStateBytes)
	if err != nil {
		log.Printf("intendSettle err", err)
		return
	}
	log.Printf("IntendSettle tx detail", tx)

	msg := subscribe.NewMsgGuardProof(channelId, tx.Hash().Hex(), m.transactor.Key.GetAddress())
	m.transactor.BroadcastTx(msg)
}
