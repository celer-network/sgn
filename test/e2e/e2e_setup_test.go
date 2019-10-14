// Setup mainchain and sgn sidechain etc for e2e tests
package e2e

import (
	"flag"
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"

	tf "github.com/celer-network/sgn/testing"
)

// TestMain handles common setup (start mainchain, deploy, build etc)
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
	ethProc, err := StartChain()
	chkErr(err, "starting chain")

	// build binaries should take long enough for geth to be fully started
	err = buildBins(outRootDir)
	chkErr(err, "build binaries")

	// // deploy contracts and fund ethpool etc, also update appAddrMap
	// // first fund svrAddr 100 ETH
	// err = tf.FundAddr("100000000000000000000", []*ctype.Addr{&svrAddr})
	// chkErr(err, "fund server")
	// tf.E2eProfile, tokenAddrErc20 = SetupOnChain(appAddrMap)

	// // profile.json and profile2.json are multi-server single OSP profiles
	// // profile2.json is used by StartC2WithoutProxy
	// noProxyProfile = outRootDir + "profile.json"
	// saveProfile(tf.E2eProfile, noProxyProfile)
	// p2 := *tf.E2eProfile
	// p2.SvrRPC = s2Addr
	// saveProfile(&p2, outRootDir+"profile2.json")
	// // multiosp.json is for osp-to-osp test, 2nd client profile
	// e2eProfile2 := *tf.E2eProfile
	// e2eProfile2.SvrETHAddr = server2AddrStr
	// e2eProfile2.SvrRPC = "localhost:10001"
	// saveProfile(&e2eProfile2, outRootDir+"multiosp.json")

	//TODO: update rt_config.json and tokens.json

	// run all e2e tests
	ret := m.Run()

	if ret == 0 {
		fmt.Println("All tests passed! 🎉🎉🎉")
		ethProc.Signal(syscall.SIGTERM)
		os.RemoveAll(outRootDir)
		os.Exit(0)
	} else {
		// fmt.Println("Tests failed. 🚧🚧🚧 Geth still running for debug. 🚧🚧🚧", "Run kill", ethProc.Pid, "to stop it")
		os.Exit(ret)
	}
}
