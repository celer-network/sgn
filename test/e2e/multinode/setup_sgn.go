// Setup mainchain and sgn sidechain etc for e2e tests
package multinode

import (
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	tf "github.com/celer-network/sgn/testing"
	"github.com/spf13/viper"
)

func setupNewSGNEnv(sgnParams *tf.SGNParams) {
	log.Infoln("deploy guard contract")
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

	log.Infoln("make prepare-sgn-data")
	repoRoot, _ := filepath.Abs("../../..")
	cmd := exec.Command("make", "prepare-sgn-data")
	cmd.Dir = repoRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Error(err)
	}

	log.Infoln("Updating config files of SGN nodes")
	for i := 0; /* 3 nodes */ i < 3; i++ {
		configPath := fmt.Sprintf("../../../docker-volumes/node%d/config.json", i)
		viper.SetConfigFile(configPath)
		err := viper.ReadInConfig()
		tf.ChkErr(err, "failed to read config")
		viper.Set(common.FlagEthGuardAddress, tf.E2eProfile.GuardAddr.String())
		viper.Set(common.FlagEthLedgerAddress, tf.E2eProfile.LedgerAddr)
		viper.WriteConfig()
	}

	log.Infoln("SetContracts")
	tf.DefaultTestEthClient.SetContracts(tf.E2eProfile.GuardAddr.String(), tf.E2eProfile.LedgerAddr.String())

	log.Infoln("make localnet-up-nodes")
	cmd = exec.Command("make", "localnet-up-nodes")
	cmd.Dir = repoRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Error(err)
	}
}
