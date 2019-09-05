package monitor

import (
	"log"

	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/x/subscribe"
)

func (m *EthMonitor) processQueue() {
	pusher := m.getPusher()
	if !pusher.ValidatorAddr.Equals(m.transactor.Key.GetAddress()) {
		return
	}

	for m.intendSettleQueue.Len() != 0 {
		m.handleIntendSettle(m.intendSettleQueue.PopFront().(*mainchain.CelerLedgerIntendSettle))
	}
}

func (m *EthMonitor) processIntendSettle(intendSettle *mainchain.CelerLedgerIntendSettle) {
	log.Printf("Process intend settle", intendSettle.ChannelId)
	channelId := intendSettle.ChannelId[:]
	request, err := subscribe.CLIQueryRequest(m.cdc, m.transactor.CliCtx, subscribe.StoreKey, channelId)
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
	_, err = m.transactor.BroadcastTx(msg)
	if err != nil {
		log.Printf("GuardProof err", err)
		return
	}
}
