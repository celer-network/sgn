// Setup mainchain and sgn sidechain etc for e2e tests
package e2e

import (
	"flag"
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/celer-network/sgn/goceler-copy/ctype"
	tf "github.com/celer-network/sgn/testing"
)

// TestMain handles common setup (start mainchain, deploy, start sidechain etc)
// and teardown. Test specific setup should be done in TestXxx
func TestMain(m *testing.M) {
	flag.Parse()
	// mkdir out root
	outRootDir = fmt.Sprintf("%s%d/", outRootDirPrefix, time.Now().Unix())
	err := os.MkdirAll(outRootDir, os.ModePerm)
	chkErr(err, "creating root dir")
	fmt.Println("Using folder:", outRootDir)
	// set testing pkg level path
	tf.SetOutRootDir(outRootDir)
	// start geth, not waiting for it to be fully ready. also watch geth proc
	// if geth exits with non-zero, os.Exit(1)
	ethProc, err := StartMainchain()
	chkErr(err, "starting chain")

	// install sgn bins
	err = installBins()
	chkErr(err, "install SGN bins")

	// set up mainchain: deploy contracts and fund ethpool etc, also update appAddrMap
	// first fund clientAddr 100 ETH
	err = tf.FundAddr("100000000000000000000", []*ctype.Addr{&clientAddr})
	chkErr(err, "fund server")
	tf.E2eProfile, tf.GuardAddr, tf.Erc20TokenAddr = SetupMainchain(appAddrMap)

	// set up sidechain (SGN)
	sgnProc, removeCmd, err := StartSidechainDefault(outRootDir)
	sleep(10) // wait for sgn to be fully ready
	chkErr(err, "start sidechain")

	// run all e2e tests
	ret := m.Run()

	ethProc.Signal(syscall.SIGTERM)
	sgnProc.Signal(syscall.SIGTERM)
	os.RemoveAll(outRootDir)
	chkErr(removeCmd.Run(), "remove sidechain directory")
	if ret == 0 {
		fmt.Println("All tests passed! 🎉🎉🎉")
		os.Exit(0)
	} else {
		fmt.Println("Tests failed. 🚧🚧🚧 Geth still running for debug. 🚧🚧🚧", "Run kill", ethProc.Pid, "to stop it")
		os.Exit(ret)
	}
}
