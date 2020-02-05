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
	tc "github.com/celer-network/sgn/test/common"
	"github.com/spf13/viper"
)

func setupNewSGNEnv(sgnParams *tc.SGNParams) {
	log.Infoln("Deploy guard contract")
	if sgnParams == nil {
		sgnParams = &tc.SGNParams{
			BlameTimeout:           big.NewInt(50),
			MinValidatorNum:        big.NewInt(1),
			MinStakingPool:         big.NewInt(100),
			SidechainGoLiveTimeout: big.NewInt(0),
			CelrAddr:               tc.E2eProfile.CelrAddr,
			MaxValidatorNum:        big.NewInt(11),
		}
	}
	tc.E2eProfile.GuardAddr = tc.DeployGuardContract(sgnParams)

	log.Infoln("make localnet-down-nodes")
	cmd := exec.Command("make", "localnet-down-nodes")
	repoRoot, _ := filepath.Abs("../../..")
	cmd.Dir = repoRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	tc.ChkErr(err, "Failed to make localnet-down-nodes")

	log.Infoln("make prepare-sgn-data")
	cmd = exec.Command("make", "prepare-sgn-data")
	cmd.Dir = repoRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	tc.ChkErr(err, "Failed to make prepare-sgn-data")

	log.Infoln("Updating config files of SGN nodes")
	for i := 0; i < 3; i++ {
		configPath := fmt.Sprintf("../../../docker-volumes/node%d/config.json", i)
		viper.SetConfigFile(configPath)
		err = viper.ReadInConfig()
		tc.ChkErr(err, "Failed to read config")
		viper.Set(common.FlagEthGuardAddress, tc.E2eProfile.GuardAddr)
		viper.Set(common.FlagEthLedgerAddress, tc.E2eProfile.LedgerAddr)
		err = viper.WriteConfig()
		tc.ChkErr(err, "Failed to write config")
	}

	log.Infoln("SetContracts")
	err = tc.DefaultTestEthClient.SetContracts(tc.E2eProfile.GuardAddr.String(), tc.E2eProfile.LedgerAddr.String())
	tc.ChkErr(err, "Failed to SetContracts")

	log.Infoln("make localnet-up-nodes")
	cmd = exec.Command("make", "localnet-up-nodes")
	cmd.Dir = repoRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	tc.ChkErr(err, "Failed to make localnet-up-nodes")
}

func shutdownNode(node uint) {
	log.Infoln("Shutdown node", node)
	cmd := exec.Command("docker-compose", "stop", fmt.Sprintf("sgnnode%d", node))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Error(err)
	}
}

func turnOffMonitor(node uint) {
	log.Infoln("Turn off node monitor", node)

	configPath := fmt.Sprintf("../../../docker-volumes/node%d/config.json", node)
	viper.SetConfigFile(configPath)
	err := viper.ReadInConfig()
	tc.ChkErr(err, "Failed to read config")
	viper.Set(common.FlagStartMonitor, false)
	err = viper.WriteConfig()
	tc.ChkErr(err, "Failed to write config")
	viper.Set(common.FlagStartMonitor, true)

	cmd := exec.Command("docker-compose", "restart", fmt.Sprintf("sgnnode%d", node))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Error(err)
	}
}
