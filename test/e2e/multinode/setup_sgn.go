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

func setupNewSGNEnv(sgnParams *tf.SGNParams) {
	// deploy guard contract
	if sgnParams == nil {
		sgnParams = &tf.SGNParams{
			BlameTimeout:           big.NewInt(50),
			MinValidatorNum:        big.NewInt(1),
			MinStakingPool:         big.NewInt(100),
			SidechainGoLiveTimeout: big.NewInt(0),
			CelrAddr:               tf.E2eProfile.CelrAddr,
		}
	}
	tf.E2eProfile.GuardAddr = tf.DeployGuardContract(sgnParams)

	// make prepare-sgn-data
	repoRoot, _ := filepath.Abs("../../..")
	cmd := exec.Command("make", "prepare-sgn-data")
	cmd.Dir = repoRoot
	if err := cmd.Run(); err != nil {
		log.Error(err)
	}

	// update config files of SGN nodes
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

	tf.DefaultTestEthClient.SetContracts(tf.E2eProfile.GuardAddr.String(), tf.E2eProfile.LedgerAddr.String())

	// make localnet-start-nodes
	cmd = exec.Command("make", "localnet-start-nodes")
	cmd.Dir = repoRoot
	if err := cmd.Run(); err != nil {
		log.Error(err)
	}
}
