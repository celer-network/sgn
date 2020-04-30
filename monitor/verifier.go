package monitor

import (
	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/x/global"
	"github.com/celer-network/sgn/x/subscribe"
	"github.com/celer-network/sgn/x/sync"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func (m *EthMonitor) verifyChange(change sync.Change) bool {
	switch change.Type {
	case sync.SyncBlock:
		return m.verifySyncBlock(change.Data)
	case sync.Subscribe:
		return m.verifySubscribe(change.Data)
	default:
		return false
	}
}

func (m *EthMonitor) verifySyncBlock(data []byte) bool {
	var block global.Block
	m.operator.CliCtx.Codec.MustUnmarshalBinaryBare(data, &block)

	syncedBlock, err := m.getLatestBlock()
	if err != nil {
		return false
	}

	log.Infoln("Verify SyncBlock", block, syncedBlock)
	return block.Number <= m.blkNum.Uint64() && block.Number > syncedBlock.Number
}

func (m *EthMonitor) verifySubscribe(data []byte) bool {
	var subscription subscribe.Subscription
	m.operator.CliCtx.Codec.MustUnmarshalBinaryBare(data, &subscription)
	log.Infoln("Verify subscription", subscription)

	deposit, err := m.ethClient.SGN.SubscriptionDeposits(
		&bind.CallOpts{},
		mainchain.Hex2Addr(subscription.EthAddress))
	if err != nil {
		log.Errorf("Failed to query subscription desposit: %s", err)
		return false
	}

	return subscription.Deposit.BigInt().Cmp(deposit) == 0
}
