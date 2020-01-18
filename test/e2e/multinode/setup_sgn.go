// Setup mainchain and sgn sidechain etc for e2e tests
package multinode

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	tf "github.com/celer-network/sgn/testing"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/spf13/viper"
)

func setupNewSGNEnv(sgnParams *tf.SGNParams) {
	log.Infoln("Deploy guard contract")
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

	log.Infoln("make localnet-down-nodes")
	cmd := exec.Command("make", "localnet-down-nodes")
	repoRoot, _ := filepath.Abs("../../..")
	cmd.Dir = repoRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Error(err)
	}

	log.Infoln("make prepare-sgn-data")
	cmd = exec.Command("make", "prepare-sgn-data")
	cmd.Dir = repoRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Error(err)
	}

	log.Infoln("Updating config files of SGN nodes")
	for i := 0; i < 3; i++ {
		configPath := fmt.Sprintf("../../../docker-volumes/node%d/config.json", i)
		viper.SetConfigFile(configPath)
		err := viper.ReadInConfig()
		tf.ChkErr(err, "Failed to read config")
		viper.Set(common.FlagEthGuardAddress, tf.E2eProfile.GuardAddr)
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
	tf.ChkErr(err, "Failed to read config")
	viper.Set(common.FlagStartMonitor, false)
	viper.WriteConfig()

	cmd := exec.Command("docker-compose", "restart", fmt.Sprintf("sgnnode%d", node))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Error(err)
	}
}

func addValidators(ethkss []string, ethpps []string, sgnops []string, amts []*big.Int) {
	for i := 0; i < len(ethkss); i++ {
		log.Infoln("Adding validator", i)
		err := addValidator(ethkss[i], ethpps[i], sgnops[i], amts[i])
		tf.ChkErr(err, "Failed to add validator")
	}
}

func addValidator(ethks string, ethpp string, sgnop string, amt *big.Int) error {
	// get auth
	keystoreBytes, err := ioutil.ReadFile(ethks)
	if err != nil {
		return err
	}
	key, err := keystore.DecryptKey(keystoreBytes, ethpp)
	if err != nil {
		return err
	}
	auth, err := bind.NewTransactor(strings.NewReader(string(keystoreBytes)), ethpp)
	if err != nil {
		return err
	}

	// get sgnAddr
	sgnAddr, err := sdk.AccAddressFromBech32(sgnop)
	if err != nil {
		return err
	}

	err = tf.AddValidator(tf.E2eProfile.CelrContract, tf.E2eProfile.GuardAddr, auth, key.Address, sgnAddr, amt)
	if err != nil {
		return err
	}

	return nil
}
