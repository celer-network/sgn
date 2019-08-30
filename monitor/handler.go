package monitor

import (
	"log"

	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/x/global"
	"github.com/celer-network/sgn/x/subscribe"
	"github.com/ethereum/go-ethereum/core/types"
)

func (m *EthMonitor) handleNewBlock(header *types.Header) {
	msg := global.NewMsgSyncBlock(header.Number.Uint64(), m.transactor.Key.GetAddress())
	_, err := m.transactor.BroadcastTx(msg)
	if err != nil {
		log.Printf("SyncBlock err", err)
	}
}

func (m *EthMonitor) handleIntendSettle(intendSettle *mainchain.CelerLedgerIntendSettle) {
	request, err := subscribe.CLIQueryRequest(m.cdc, m.transactor.CliCtx, subscribe.StoreKey, intendSettle.ChannelId[:])
	if err != nil {
		log.Printf("Query request err", err)
		return
	}

	if intendSettle.SeqNums[request.PeerFromIndex].Uint64() >= request.SeqNum {
		log.Printf("Ignore the intendSettle event due to larger seqNum")
		return
	}

	tx, err := m.ethClient.Ledger.IntendSettle(m.ethClient.Auth, request.SignedSimplexStateBytes)
	if err != nil {
		log.Printf("Tx err", err)
		return
	}
	log.Printf("Tx detail", tx)
}
