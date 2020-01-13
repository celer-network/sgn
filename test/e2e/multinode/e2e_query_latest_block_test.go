package multinode

import (
	"context"
	"math/big"
	"testing"

	"github.com/celer-network/goutils/log"
	tf "github.com/celer-network/sgn/testing"
	"github.com/celer-network/sgn/x/global"
	"github.com/stretchr/testify/assert"
)

func setUpQueryLatestBlock() {
	log.Infoln("Set up new sgn env")
	setupNewSGNEnv(nil)
	amts := []*big.Int{big.NewInt(1000000000000000000), big.NewInt(1000000000000000000), big.NewInt(1000000000000000000)}
	addValidators(ethKeystores[:], ethKeystorePps[:], sgnOperators[:], amts)
	tf.SleepWithLog(10, "sgn syncing")
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

	conn := tf.DefaultTestEthClient.Client

	transactor := tf.NewTransactor(
		t,
		sgnCLIHomes[1],
		sgnChainID,
		sgnNodeURIs[1],
		sgnTransactors[1],
		sgnPassphrase,
		sgnGasPrice,
	)

	blockSGN, err := global.CLIQueryLatestBlock(transactor.CliCtx, global.RouterKey)
	tf.ChkTestErr(t, err, "failed to query latest synced block on sgn")
	log.Infof("Latest block number on SGN is %d", blockSGN.Number)

	header, err := conn.HeaderByNumber(context.Background(), nil)
	tf.ChkTestErr(t, err, "failed to query latest synced block on mainchain")
	log.Infof("Latest block number on mainchain is %d", header.Number)

	assert.GreaterOrEqual(t, header.Number.Uint64(), blockSGN.Number, "blkNumMain should be greater than or equal to blockSGN.Number")
	assert.Greater(t, blockSGN.Number, uint64(0), "blockSGN.Number should be larger than 0")
}
