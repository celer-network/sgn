package e2e

import (
	"os"
	"testing"

	tf "github.com/celer-network/sgn/testing"
	"github.com/celer-network/sgn/testing/log"
	"github.com/celer-network/sgn/x/global"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
)

func setUpQueryLatestBlock() []tf.Killable {
	res := setupNewSGNEnv(nil, "query_latest_block")
	sleepWithLog(20, "sgn syncing")

	return res
}

func TestE2EQueryLatestBlock(t *testing.T) {
	toKill := setUpQueryLatestBlock()
	defer tf.TearDown(toKill)

	t.Run("e2e-queryLatestBlock", func(t *testing.T) {
		t.Run("queryLatestBlockTest", queryLatestBlockTest)
	})
}

func queryLatestBlockTest(t *testing.T) {
	// t.Parallel()

	log.Info("=====================================================================")
	log.Info("======================== Test queryLatestBlock ===========================")

	conn, err := ethclient.Dial(tf.EthInstance)
	if err != nil {
		os.Exit(1)
	}

	blockSGN, err := global.CLIQueryLatestBlock(tf.Transactor.CliCtx, global.RouterKey)
	tf.ChkErr(err, "failed to query latest synced block on sgn")
	log.Infof("Latest block number on SGN is %d", blockSGN.Number)

	blkNumMain, err := tf.GetLatestBlkNum(conn)
	tf.ChkErr(err, "failed to query latest synced block on mainchain")
	log.Infof("Latest block number on mainchain is %d", blkNumMain)

	assert.GreaterOrEqual(t, blkNumMain.Uint64(), blockSGN.Number, "blkNumMain should be greater than or equal to blockSGN.Number")
}
