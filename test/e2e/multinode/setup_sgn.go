// Setup mainchain and sgn sidechain etc for e2e tests
package multinode

import (
	"fmt"
	"math/big"
	"os/exec"
	"path/filepath"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	tf "github.com/celer-network/sgn/testing"
	"github.com/spf13/viper"
)

func setupNewSGNEnv() {
	// deploy guard contract
	sgnParams := &tf.SGNParams{
		BlameTimeout:           big.NewInt(50),
		MinValidatorNum:        big.NewInt(1),
		MinStakingPool:         big.NewInt(100),
		SidechainGoLiveTimeout: big.NewInt(0),
	}

	tf.E2eProfile.GuardAddr = tf.DeployGuardContract(sgnParams)

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
		viper.Set(common.FlagEthGuardAddress, tf.E2eProfile.GuardAddr.String())
		viper.Set(common.FlagEthLedgerAddress, tf.E2eProfile.LedgerAddr)
		viper.WriteConfig()
	}

	// update config.json
	// TODO: better config.json solution
	viper.SetConfigFile("../../../config.json")
	err := viper.ReadInConfig()
	tf.ChkErr(err, "failed to read config")
	viper.Set(common.FlagEthGuardAddress, tf.E2eProfile.GuardAddr.String())
	viper.Set(common.FlagEthLedgerAddress, tf.E2eProfile.LedgerAddr)
	viper.Set(common.FlagEthWS, "ws://127.0.0.1:8546")
	clientKeystore, err := filepath.Abs("../../keys/client0.json")
	tf.ChkErr(err, "get client keystore path")
	viper.Set(common.FlagEthKeystore, clientKeystore)
	// TODO: set operator, transactors
	viper.WriteConfig()

	// set up eth client and transactor
	ks_path, _ := filepath.Abs("../../keys/client0.json")
	tf.SetupEthClient(ks_path)

	// make localnet-start-nodes
	cmd = exec.Command("make", "localnet-start-nodes")
	cmd.Dir = repoRoot
	if err := cmd.Run(); err != nil {
		log.Error(err)
	}
}
