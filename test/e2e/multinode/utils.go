// Setup mainchain and sgn sidechain etc for e2e tests
package multinode

import (
	"context"
	"fmt"
	"math/big"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	tf "github.com/celer-network/sgn/testing"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
)

type SGNParams struct {
	blameTimeout           *big.Int
	minValidatorNum        *big.Int
	minStakingPool         *big.Int
	sidechainGoLiveTimeout *big.Int
	startGateway           bool
}

func setupNewSGNEnv() mainchain.Addr {
	// deploy guard contract
	sgnParams := &SGNParams{
		blameTimeout:           big.NewInt(50),
		minValidatorNum:        big.NewInt(1),
		minStakingPool:         big.NewInt(100),
		sidechainGoLiveTimeout: big.NewInt(0),
	}
	conn, err := ethclient.Dial(tf.EthInstance)
	tf.ChkErr(err, "failed to connect to the Ethereum")
	ctx := context.Background()
	ethbasePrivKey, _ := crypto.HexToECDSA(etherBasePriv)
	etherBaseAuth := bind.NewKeyedTransactor(ethbasePrivKey)
	price := big.NewInt(2e9) // 2Gwei
	etherBaseAuth.GasPrice = price
	etherBaseAuth.GasLimit = 7000000
	guardAddr, tx, _, err := mainchain.DeployGuard(etherBaseAuth, conn, e2eProfile.CelrAddr, sgnParams.blameTimeout, sgnParams.minValidatorNum, sgnParams.minStakingPool, sgnParams.sidechainGoLiveTimeout)
	e2eProfile.GuardAddr = guardAddr
	tf.ChkErr(err, "failed to deploy Guard contract")
	tf.WaitMinedWithChk(ctx, conn, tx, 0, "Deploy Guard "+guardAddr.Hex())

	// make prepare-sgn-data
	repoRoot, _ := filepath.Abs("../../..")
	cmd := exec.Command("make", "prepare-sgn-data")
	cmd.Dir = repoRoot
	if err := cmd.Run(); err != nil {
		log.Error(err)
	}

	// update SGN config
	log.Infoln("Updating SGN's config.json")
	for i := 0; /* 3 nodes */ i < 3; i++ {
		configPath := fmt.Sprintf("../../../docker-volumes/node%d/config.json", i)
		viper.SetConfigFile(configPath)
		err := viper.ReadInConfig()
		tf.ChkErr(err, "failed to read config")
		viper.Set(common.FlagEthGuardAddress, e2eProfile.GuardAddr.String())
		viper.Set(common.FlagEthLedgerAddress, e2eProfile.LedgerAddr)
		viper.WriteConfig()
	}

	// update config.json
	// TODO: better config.json solution
	viper.SetConfigFile("../../../config.json")
	err = viper.ReadInConfig()
	tf.ChkErr(err, "failed to read config")
	viper.Set(common.FlagEthGuardAddress, e2eProfile.GuardAddr.String())
	viper.Set(common.FlagEthLedgerAddress, e2eProfile.LedgerAddr)
	viper.Set(common.FlagEthWS, "ws://127.0.0.1:8546")
	sgnCliHome, _ := filepath.Abs("../../../docker-volumes/node0/sgncli")
	viper.Set(common.FlagSgnCLIHome, sgnCliHome)
	sgnNodeHome, _ := filepath.Abs("../../../docker-volumes/node0/sgn")
	viper.Set(common.FlagSgnNodeHome, sgnNodeHome)
	clientKeystore, err := filepath.Abs("../../keys/client0.json")
	tf.ChkErr(err, "get client keystore path")
	viper.Set(common.FlagEthKeystore, clientKeystore)
	// TODO: set operator, transactors
	viper.WriteConfig()

	// set up eth client and transactor
	ks_path, _ := filepath.Abs("../../keys/client0.json")
	tf.SetupEthClient(ks_path)
	tf.SetupTransactor()

	// make localnet-start-nodes
	cmd = exec.Command("make", "localnet-start-nodes")
	cmd.Dir = repoRoot
	if err := cmd.Run(); err != nil {
		log.Error(err)
	}

	return guardAddr
}

func sleep(second time.Duration) {
	time.Sleep(second * time.Second)
}

func sleepWithLog(second time.Duration, waitFor string) {
	log.Infof("Sleep %d seconds for %s", second, waitFor)
	sleep(second)
}

func sleepBlocksWithLog(count time.Duration, waitFor string) {
	sleepWithLog(count*sgnBlockInterval, waitFor)
}