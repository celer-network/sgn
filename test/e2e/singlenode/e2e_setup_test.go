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
	tf "github.com/celer-network/sgn/testing"
)

type TestProfile struct {
	DisputeTimeout uint64
	LedgerAddr     mainchain.Addr
	GuardAddr      mainchain.Addr
	CelrAddr       mainchain.Addr
	CelrContract   *mainchain.ERC20
}

// used by setup_onchain and tests
var (
	client0Addr = mainchain.Hex2Addr(tf.Client0AddrStr)
	client1Addr = mainchain.Hex2Addr(tf.Client1AddrStr)
)

// runtime variables, will be initialized by TestMain
var (
	// root dir with ending / for all files, OutRootDirPrefix + epoch seconds
	// due to testframework etc in a different testing package, we have to define
	// same var in testframework.go and expose a set api
	outRootDir string
	e2eProfile *TestProfile
)

// TestMain handles common setup (start mainchain, deploy, start sidechain etc)
// and teardown. Test specific setup should be done in TestXxx
func TestMain(m *testing.M) {
	flag.Parse()
	log.EnableColor()
	common.EnableLogLongFile()

	// mkdir out root
	outRootDir = fmt.Sprintf("%s%d/", tf.OutRootDirPrefix, time.Now().Unix())
	err := os.MkdirAll(outRootDir, os.ModePerm)
	tf.ChkErr(err, "creating root dir")
	log.Infoln("Using folder:", outRootDir)
	// set testing pkg level path
	// start geth, not waiting for it to be fully ready. also watch geth proc
	// if geth exits with non-zero, os.Exit(1)
	ethProc, err := startMainchain()
	tf.ChkErr(err, "starting mainchain")
	tf.SleepWithLog(2, "starting mainchain")

	// set up mainchain: deploy contracts and fund ethpool etc
	// first fund client0Addr 100 ETH
	err = tf.FundAddr("1"+strings.Repeat("0", 20), []*mainchain.Addr{&client0Addr})
	tf.ChkErr(err, "fund server")
	e2eProfile = setupMainchain()

	// make install sgn and sgncli
	err = installSgn()
	tf.ChkErr(err, "installing sgn and sgncli")

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
