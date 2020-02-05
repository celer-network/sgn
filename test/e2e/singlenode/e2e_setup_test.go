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
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	tc "github.com/celer-network/sgn/test/common"
)

var (
	CLIHome = os.ExpandEnv("$HOME/.sgncli")

	// root dir with ending / for all files, OutRootDirPrefix + epoch seconds
	// due to testframework etc in a different testing package, we have to define
	// same var in testframework.go and expose a set api
	outRootDir string
)

// TestMain handles common setup (start mainchain, deploy, start sidechain etc)
// and teardown. Test specific setup should be done in TestXxx
func TestMain(m *testing.M) {
	flag.Parse()
	log.EnableColor()
	common.EnableLogLongFile()

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
	// first fund client0Addr 100 ETH
	err = tc.FundAddrsETH("1"+strings.Repeat("0", 20), []mainchain.Addr{mainchain.Hex2Addr(tc.EthAddresses[0])})
	tc.ChkErr(err, "fund server")
	tc.SetupDefaultTestEthClient()
	tc.SetupE2eProfile()

	// make install sgn and sgncli
	err = installSgn()
	tc.ChkErr(err, "installing sgn and sgncli")

	// run all e2e tests
	ret := m.Run()

	ethProc.Signal(syscall.SIGTERM)
	if ret == 0 {
		log.Infoln("All tests passed! 🎉🎉🎉")
		os.RemoveAll(outRootDir)
		os.Exit(0)
	} else {
		log.Errorln("Tests failed. 🚧🚧🚧 Geth still running for debug. 🚧🚧🚧", "Run kill", ethProc.Pid, "to stop it")
		os.Exit(ret)
	}
}
