// Setup mainchain and sgn sidechain etc for e2e tests
package singlenode

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	tc "github.com/celer-network/sgn/testing/common"
)

// TestMain handles common setup (start mainchain, deploy, start sidechain etc)
// and teardown. Test specific setup should be done in TestXxx
func TestMain(m *testing.M) {
	flag.Parse()
	log.EnableColor()

	// mkdir out root
	outRootDir = fmt.Sprintf("%s%d/", tc.OutRootDirPrefix, time.Now().Unix())
	err := os.MkdirAll(outRootDir, os.ModePerm)
	tc.ChkErr(err, "creating root dir")
	log.Infoln("Using folder:", outRootDir)
	// set testing pkg level path
	// start geth, not waiting for it to be fully ready. also watch geth proc
	// if geth exits with non-zero, os.Exit(1)
	ethProc, err := startMainchain()
	tc.ChkErr(err, "starting mainchain")
	tc.SleepWithLog(2, "starting mainchain")

	// set up mainchain: deploy contracts and fund ethpool etc
	// first fund each account 100 ETH
	addrs := []mainchain.Addr{
		mainchain.Hex2Addr(tc.ValEthAddrs[0]),
		mainchain.Hex2Addr(tc.DelEthAddrs[0]),
		mainchain.Hex2Addr(tc.DelEthAddrs[1]),
		mainchain.Hex2Addr(tc.DelEthAddrs[2]),
		mainchain.Hex2Addr(tc.DelEthAddrs[3]),
		mainchain.Hex2Addr(tc.ClientEthAddrs[0]),
		mainchain.Hex2Addr(tc.ClientEthAddrs[1]),
	}
	err = tc.FundAddrsETH("1"+strings.Repeat("0", 20), addrs)
	tc.ChkErr(err, "fund eth")
	tc.SetupEthClients()
	tc.SetupE2eProfile()

	// fund CELR to each eth account
	log.Infoln("fund each addr 10 million CELR")
	err = tc.FundAddrsErc20(tc.E2eProfile.CelrAddr, addrs, "1"+strings.Repeat("0", 25))
	tc.ChkErr(err, "fund each addr ERC20")

	// make install sgn and sgncli
	err = installSgn()
	tc.ChkErr(err, "installing sgn and sgncli")

	tc.SetupSidechain()
	// run all e2e tests
	ret := m.Run()

	ethProc.Signal(syscall.SIGTERM)
	if ret == 0 {
		log.Infoln("All tests passed! 🎉🎉🎉")
		os.RemoveAll(outRootDir)
		os.Exit(0)
	} else {
		log.Errorln("Tests failed. 🔥🔥🔥")
		os.Exit(ret)
	}
}
