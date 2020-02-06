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
	tc "github.com/celer-network/sgn/test/common"
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
	tc.SleepWithLog(5, "geth start")

	log.Infoln("fund each validator's ETH address 100 ETH")
	addr0 := mainchain.Hex2Addr(tc.ValEthAddrs[0])
	addr1 := mainchain.Hex2Addr(tc.ValEthAddrs[1])
	addr2 := mainchain.Hex2Addr(tc.ValEthAddrs[2])
	addr3 := mainchain.Hex2Addr(tc.ClientEthAddrs[0])
	err := tc.FundAddrsETH("1"+strings.Repeat("0", 20), []mainchain.Addr{addr0, addr1, addr2, addr3})
	tc.ChkErr(err, "fund each validator ETH")

	log.Infoln("set up mainchain")
	tc.SetupEthClients()
	tc.SetupE2eProfile()

	// fund CELR to each eth account
	log.Infoln("fund each validator 10 million CELR")
	err = tc.FundAddrsErc20(tc.E2eProfile.CelrAddr, []mainchain.Addr{addr0, addr1, addr2, addr3}, "1"+strings.Repeat("0", 25))
	tc.ChkErr(err, "fund each validator ERC20")

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
