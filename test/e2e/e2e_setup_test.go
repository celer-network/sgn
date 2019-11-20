// Setup mainchain and sgn sidechain etc for e2e tests
package e2e

import (
	"flag"
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"

	log "github.com/celer-network/sgn/clog"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/ctype"
	"github.com/celer-network/sgn/mainchain"
	tf "github.com/celer-network/sgn/testing"
)

// used by setup_onchain and tests
var (
	etherBaseAddr = ctype.Hex2Addr(etherBaseAddrStr)
	client0Addr   = ctype.Hex2Addr(client0AddrStr)
	client1Addr   = ctype.Hex2Addr(client1AddrStr)
)

// runtime variables, will be initialized by TestMain
var (
	// root dir with ending / for all files, outRootDirPrefix + epoch seconds
	// due to testframework etc in a different testing package, we have to define
	// same var in testframework.go and expose a set api
	outRootDir    string
	envDir        = "../../testing/env"
	e2eProfile    *common.CProfile
	celrContract  *mainchain.ERC20
	guardAddr     ctype.Addr
	mockCelerAddr ctype.Addr
)

// TestMain handles common setup (start mainchain, deploy, start sidechain etc)
// and teardown. Test specific setup should be done in TestXxx
func TestMain(m *testing.M) {
	flag.Parse()
	log.EnableColor()
	log.EnableLongFile()
	// mkdir out root
	tf.SetEnvDir(envDir)
	outRootDir = fmt.Sprintf("%s%d/", outRootDirPrefix, time.Now().Unix())
	err := os.MkdirAll(outRootDir, os.ModePerm)
	tf.ChkErr(err, "creating root dir")
	fmt.Println("Using folder:", outRootDir)
	// set testing pkg level path
	// start geth, not waiting for it to be fully ready. also watch geth proc
	// if geth exits with non-zero, os.Exit(1)
	ethProc, err := startMainchain()
	tf.ChkErr(err, "starting chain")
	sleep(2)

	// set up mainchain: deploy contracts and fund ethpool etc
	// first fund client0Addr 100 ETH
	err = tf.FundAddr("100000000000000000000", []*ctype.Addr{&client0Addr})
	tf.ChkErr(err, "fund server")
	e2eProfile, mockCelerAddr = setupMainchain()

	// run all e2e tests
	ret := m.Run()

	ethProc.Signal(syscall.SIGTERM)
	if ret == 0 {
		fmt.Println("All tests passed! 🎉🎉🎉")
		os.RemoveAll(outRootDir)
		os.Exit(0)
	} else {
		fmt.Println("Tests failed. 🚧🚧🚧 Geth still running for debug. 🚧🚧🚧", "Run kill", ethProc.Pid, "to stop it")
		os.Exit(ret)
	}
}
