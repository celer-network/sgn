package multinode

import (
	"os"
	"testing"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	tf "github.com/celer-network/sgn/testing"
	"github.com/celer-network/sgn/x/global"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
)

func setUpQueryLatestBlock() mainchain.Addr {
	log.Infoln("set up new sgn env")
	guardAddr := setupNewSGNEnv()
	sleepWithLog(15, "sgn syncing")

	return guardAddr
}

func TestE2EQueryLatestBlock(t *testing.T) {
	setUpQueryLatestBlock()

	t.Run("e2e-queryLatestBlock", func(t *testing.T) {
		t.Run("queryLatestBlockTest", queryLatestBlockTest)
	})
}

func queryLatestBlockTest(t *testing.T) {
	log.Info("=====================================================================")
	log.Info("======================== Test queryLatestBlock ===========================")

	conn, err := ethclient.Dial(tf.EthInstance)
	if err != nil {
		os.Exit(1)
	}

	blockSGN, err := global.CLIQueryLatestBlock(tf.Transactor.CliCtx, global.RouterKey)
	tf.ChkTestErr(t, err, "failed to query latest synced block on sgn")
	log.Infof("Latest block number on SGN is %d", blockSGN.Number)

	blkNumMain, err := tf.GetLatestBlkNum(conn)
	tf.ChkTestErr(t, err, "failed to query latest synced block on mainchain")
	log.Infof("Latest block number on mainchain is %d", blkNumMain)

	assert.GreaterOrEqual(t, blkNumMain.Uint64(), blockSGN.Number, "blkNumMain should be greater than or equal to blockSGN.Number")
}
