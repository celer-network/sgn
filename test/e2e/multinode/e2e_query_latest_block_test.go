package multinode

import (
	"context"
	"os"
	"testing"

	"github.com/celer-network/goutils/log"
	tf "github.com/celer-network/sgn/testing"
	"github.com/celer-network/sgn/x/global"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
)

func setUpQueryLatestBlock() {
	log.Infoln("set up new sgn env")
	setupNewSGNEnv()
	tf.SleepWithLog(5, "sgn syncing")
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

	header, err := conn.HeaderByNumber(context.Background(), nil)
	tf.ChkTestErr(t, err, "failed to query latest synced block on mainchain")
	log.Infof("Latest block number on mainchain is %d", header.Number)

	assert.GreaterOrEqual(t, header.Number.Uint64(), blockSGN.Number, "blkNumMain should be greater than or equal to blockSGN.Number")
}
