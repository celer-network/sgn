package monitor

import (
	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/x/global"
	"github.com/celer-network/sgn/x/sync"
)

func (m *EthMonitor) verifyChange(change sync.Change) bool {
	switch change.Type {
	case sync.SyncBlock:
		return m.verifySyncBlock(change.Data)
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

	log.Infof("Verify SyncBlock", block, syncedBlock)
	return block.Number <= m.blkNum.Uint64() && block.Number > syncedBlock.Number
}
