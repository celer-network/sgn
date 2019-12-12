// Setup mainchain and sgn sidechain etc for e2e tests
package multinode

import (
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	tf "github.com/celer-network/sgn/testing"
)

// TestMain handles common setup (start mainchain, deploy, start sidechain etc)
// and teardown. Test specific setup should be done in TestXxx
func TestMain(m *testing.M) {
	flag.Parse()
	log.EnableColor()
	common.EnableLogLongFile()

	repoRoot, _ := filepath.Abs("../../..")

	log.Infoln("make localnet-down")
	cmd := exec.Command("make", "localnet-down")
	cmd.Dir = repoRoot
	if err := cmd.Run(); err != nil {
		log.Error(err)
	}

	log.Infoln("build dockers, get geth, build sgn binary")
	cmd = exec.Command("make", "prepare-docker-env")
	cmd.Dir = repoRoot
	if err := cmd.Run(); err != nil {
		log.Error(err)
	}

	log.Infoln("start geth container")
	cmd = exec.Command("make", "localnet-start-geth")
	cmd.Dir = repoRoot
	if err := cmd.Run(); err != nil {
		log.Error(err)
	}
	tf.SleepWithLog(5, "geth start")

	log.Infoln("first fund client0Addr 100 ETH")
	err := tf.FundAddr("1"+strings.Repeat("0", 20), []*mainchain.Addr{&tf.Client0Addr})
	tf.ChkErr(err, "fund client0")
	log.Infoln("set up mainchain")
	tf.E2eProfile = setupMainchain()

	log.Infoln("run all e2e tests")
	ret := m.Run()

	if ret == 0 {
		log.Infoln("All tests passed! 🎉🎉🎉")
		log.Infoln("Tearing down all containers...")
		cmd = exec.Command("make", "localnet-down")
		cmd.Dir = repoRoot
		if err := cmd.Run(); err != nil {
			log.Error(err)
		}
		os.Exit(0)
	} else {
		log.Errorln("Tests failed. 🚧🚧🚧 Geth and sgn containers are still running for debug. 🚧🚧🚧 Run \"make localnet-down\" to stop them")
		os.Exit(ret)
	}
}

// setupMainchain deploy contracts, and do setups
// return profile, tokenAddrErc20
func setupMainchain() *tf.TestProfile {
	ethClient := tf.EthClient
	err := ethClient.SetupClient(tf.EthInstance)
	tf.ChkErr(err, "failed to connect to the Ethereum")
	err = ethClient.SetupAuth("../../keys/client0.json", "")
	tf.ChkErr(err, "failed to create auth")

	ledgerAddr := tf.DeployLedgerContract()

	// Deploy sample ERC20 contract (CELR)
	tf.LogBlkNum(ethClient.Client)
	erc20Addr, erc20 := tf.DeployERC20Contract()

	return &tf.TestProfile{
		// hardcoded values
		DisputeTimeout: 10,
		// deployed addresses
		LedgerAddr:   ledgerAddr,
		CelrAddr:     erc20Addr,
		CelrContract: erc20,
	}
}
