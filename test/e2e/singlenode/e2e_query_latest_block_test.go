package singlenode

import (
	"context"
	"testing"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	tc "github.com/celer-network/sgn/test/common"
	"github.com/celer-network/sgn/x/global"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func setUpQueryLatestBlock() []tc.Killable {
	res := setupNewSGNEnv(nil, "query_latest_block")
	tc.SleepWithLog(10, "sgn syncing")

	return res
}

func TestE2EQueryLatestBlock(t *testing.T) {
	toKill := setUpQueryLatestBlock()
	defer tc.TearDown(toKill)

	t.Run("e2e-queryLatestBlock", func(t *testing.T) {
		t.Run("queryLatestBlockTest", queryLatestBlockTest)
	})
}

func queryLatestBlockTest(t *testing.T) {
	// t.Parallel()

	log.Info("=====================================================================")
	log.Info("======================== Test queryLatestBlock ===========================")

	conn := tc.DefaultTestEthClient.Client

	transactor := tc.NewTransactor(
		t,
		CLIHome,
		viper.GetString(common.FlagSgnChainID),
		viper.GetString(common.FlagSgnNodeURI),
		viper.GetStringSlice(common.FlagSgnTransactors)[1],
		viper.GetString(common.FlagSgnPassphrase),
		viper.GetString(common.FlagSgnGasPrice),
	)

	blockSGN, err := global.CLIQueryLatestBlock(transactor.CliCtx, global.RouterKey)
	tc.ChkTestErr(t, err, "failed to query latest synced block on sgn")
	log.Infof("Latest block number on SGN is %d", blockSGN.Number)

	header, err := conn.HeaderByNumber(context.Background(), nil)
	tc.ChkTestErr(t, err, "failed to query latest synced block on mainchain")
	log.Infof("Latest block number on mainchain is %d", header.Number)

	assert.GreaterOrEqual(t, header.Number.Uint64(), blockSGN.Number, "blkNumMain should be greater than or equal to blockSGN.Number")
	assert.Greater(t, blockSGN.Number, uint64(0), "blockSGN.Number should be larger than 0")
}
