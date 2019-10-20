package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/celer-network/sgn/ctype"
	tf "github.com/celer-network/sgn/testing"
	"github.com/celer-network/sgn/testing/log"
	"github.com/celer-network/sgn/x/global"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
)

func setUpQueryLatestBlock() []tf.Killable {
	// TODO: duplicate code in SetupMainchain(), need to put these in a function
	ctx := context.Background()
	conn, err := ethclient.Dial(tf.EthInstance)
	tf.ChkErr(err, "failed to connect to the Ethereum")
	ethbasePrivKey, _ := crypto.HexToECDSA(etherBasePriv)
	etherBaseAuth := bind.NewKeyedTransactor(ethbasePrivKey)
	price := big.NewInt(2e9) // 2Gwei
	etherBaseAuth.GasPrice = price
	etherBaseAuth.GasLimit = 7000000

	// deploy guard contract
	tf.LogBlkNum(conn)
	blameTimeout := big.NewInt(50)
	minValidatorNum := big.NewInt(1)
	minStakingPool := big.NewInt(100)
	sidechainGoLiveTimeout := big.NewInt(0)
	GuardAddr = DeployGuardContract(ctx, etherBaseAuth, conn, ctype.Hex2Addr(Erc20TokenAddr), blameTimeout, minValidatorNum, minStakingPool, sidechainGoLiveTimeout)

	// update SGN config
	UpdateSGNConfig()

	// start sgn sidechain
	sgnProc, err := StartSidechainDefault(outRootDir)
	tf.ChkErr(err, "start sidechain")
	fmt.Println("Sleep for 20 seconds to let sgn be fully ready")
	sleep(20) // wait for sgn to be fully ready

	tf.SetupEthClient()
	tf.SetupTransactor()

	return []tf.Killable{sgnProc}
}

func TestE2EQueryLatestBlock(t *testing.T) {
	toKill := setUpQueryLatestBlock()
	defer tf.TearDown(toKill)

	t.Run("e2e-queryLatestBlock", func(t *testing.T) {
		t.Run("queryLatestBlockTest", queryLatestBlockTest)
	})
}

func queryLatestBlockTest(t *testing.T) {
	t.Parallel()
	queryLatestBlock(t)
}

func queryLatestBlock(t *testing.T) {
	log.Info("=====================================================================")
	log.Info("======================== Test queryLatestBlock ===========================")

	conn, err := ethclient.Dial(tf.EthInstance)
	if err != nil {
		os.Exit(1)
	}

	route := fmt.Sprintf("custom/%s/%s", global.ModuleName, global.QueryLatestBlock)
	resJson, _, err := tf.Transactor.CliCtx.Query(route)
	var result map[string]interface{}
	json.Unmarshal(resJson, &result)
	blkNumSGN := new(big.Int)
	blkNumSGN, _ = blkNumSGN.SetString(result["number"].(string), 10)
	log.Infof("Latest block number on SGN is %d", blkNumSGN)

	blkNumMain, err := tf.GetLatestBlkNum(conn)
	if err != nil {
		os.Exit(1)
	}
	log.Infof("Latest block number on mainchain is %d", blkNumMain)

	assert.Equal(t, err, nil, "The command should run successfully")
	diff := new(big.Int).Sub(blkNumMain, blkNumSGN)
	assert.GreaterOrEqual(t, big.NewInt(5).Cmp(diff), 0, "blkNumMain should be greater than or equal to blkNumSGN")
}
