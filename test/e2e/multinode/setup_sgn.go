// Setup mainchain and sgn sidechain etc for e2e tests
package multinode

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	tc "github.com/celer-network/sgn/test/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/viper"
)

func setupNewSGNEnv(sgnParams *tc.SGNParams) {
	log.Infoln("Deploy DPoS and SGN contracts")
	if sgnParams == nil {
		sgnParams = &tc.SGNParams{
			CelrAddr:               tc.E2eProfile.CelrAddr,
			GovernProposalDeposit:  big.NewInt(1), // TODO: use a more practical value
			GovernVoteTimeout:      big.NewInt(1), // TODO: use a more practical value
			BlameTimeout:           big.NewInt(50),
			MinValidatorNum:        big.NewInt(1),
			MaxValidatorNum:        big.NewInt(11),
			MinStakingPool:         big.NewInt(100),
			IncreaseRateWaitTime:   big.NewInt(1), // TODO: use a more practical value
			SidechainGoLiveTimeout: big.NewInt(0),
		}
	}
	var tx *types.Transaction
	tx, tc.E2eProfile.DPoSAddr, tc.E2eProfile.SGNAddr = tc.DeployDPoSSGNContracts(sgnParams)
	tc.WaitMinedWithChk(context.Background(), tc.EthClient, tx, tc.BlockDelay, tc.PollingInterval, "DeployDPoSSGNContracts")

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
		viper.Set(common.FlagEthDPoSAddress, tc.E2eProfile.DPoSAddr)
		viper.Set(common.FlagEthSGNAddress, tc.E2eProfile.SGNAddr)
		viper.Set(common.FlagEthLedgerAddress, tc.E2eProfile.LedgerAddr)
		err = viper.WriteConfig()
		tc.ChkErr(err, "Failed to write config")
	}

	err = tc.SetContracts(tc.E2eProfile.DPoSAddr, tc.E2eProfile.SGNAddr, tc.E2eProfile.LedgerAddr)
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
