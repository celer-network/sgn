// Setup mainchain and sgn sidechain etc for e2e tests
package e2e

import (
	"flag"
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	tf "github.com/celer-network/sgn/testing"
)

// CProfile struct is based on github.com/goCeler/common (commit ID: d7335ae321b67150d92de18f6589f1d1fd8b0910)
// CProfile contains configurations for CelerClient/OSP
type CProfile struct {
	ETHInstance        string `json:"ethInstance"`
	SvrETHAddr         string `json:"svrEthAddr"`
	WalletAddr         string `json:"walletAddr"`
	LedgerAddr         string `json:"ledgerAddr"`
	VirtResolverAddr   string `json:"virtResolverAddr"`
	EthPoolAddr        string `json:"ethPoolAddr"`
	PayResolverAddr    string `json:"payResolverAddr"`
	PayRegistryAddr    string `json:"payRegistryAddr"`
	RouterRegistryAddr string `json:"routerRegistryAddr"`
	SvrRPC             string `json:"svrRpc"`
	SelfRPC            string `json:"selfRpc,omitempty"`
	StoreDir           string `json:"storeDir,omitempty"`
	StoreSql           string `json:"storeSql,omitempty"`
	WebPort            string `json:"webPort,omitempty"`
	WsOrigin           string `json:"wsOrigin,omitempty"`
	ChainId            int64  `json:"chainId"`
	BlockDelayNum      uint64 `json:"blockDelayNum"`
	IsOSP              bool   `json:"isOsp,omitempty"`
	ListenOnChain      bool   `json:"listenOnChain,omitempty"`
	PollingInterval    uint64 `json:"pollingInterval"`
	DisputeTimeout     uint64 `json:"disputeTimeout"`
}

// used by setup_onchain and tests
var (
	etherBaseAddr = mainchain.Hex2Addr(etherBaseAddrStr)
	client0Addr   = mainchain.Hex2Addr(client0AddrStr)
	client1Addr   = mainchain.Hex2Addr(client1AddrStr)
)

// runtime variables, will be initialized by TestMain
var (
	// root dir with ending / for all files, outRootDirPrefix + epoch seconds
	// due to testframework etc in a different testing package, we have to define
	// same var in testframework.go and expose a set api
	outRootDir    string
	envDir        = "../../testing/env"
	e2eProfile    *CProfile
	celrContract  *mainchain.ERC20
	guardAddr     mainchain.Addr
	mockCelerAddr mainchain.Addr
)

// TestMain handles common setup (start mainchain, deploy, start sidechain etc)
// and teardown. Test specific setup should be done in TestXxx
func TestMain(m *testing.M) {
	flag.Parse()
	log.EnableColor()
	common.EnableLogLongFile()

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
	err = tf.FundAddr("100000000000000000000", []*mainchain.Addr{&client0Addr})
	tf.ChkErr(err, "fund server")
	e2eProfile, mockCelerAddr = setupMainchain()

	// make install sgn and sgncli
	err = installSgn()
	tf.ChkErr(err, "installing sgn and sgncli")

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
