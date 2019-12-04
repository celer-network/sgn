// Setup mainchain and sgn sidechain etc for e2e tests
package multinode

import (
	"context"
	"flag"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/celer-network/cChannel-eth-go/deploy"
	"github.com/celer-network/cChannel-eth-go/ethpool"
	"github.com/celer-network/cChannel-eth-go/ledger"
	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	tf "github.com/celer-network/sgn/testing"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
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
	e2eProfile    *CProfile
	guardAddr     mainchain.Addr
	mockCelerAddr mainchain.Addr
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
	sleepWithLog(5, "geth start")

	// TODO: can remove the fund distribution in genesis file?
	// log.Infoln("first fund client0Addr 100 ETH")
	// err := tf.FundAddr("100000000000000000000", []*mainchain.Addr{&client0Addr})
	// tf.ChkErr(err, "fund server")
	log.Infoln("set up mainchain")
	e2eProfile, mockCelerAddr = setupMainchain()

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
		log.Errorln("Tests failed. 🚧🚧🚧 Geth and sgn nodes are still running for debug. 🚧🚧🚧Run make localnet-down to stop it")
		os.Exit(ret)
	}
}

// setupMainchain deploy contracts, and do setups
// return profile, tokenAddrErc20
func setupMainchain() (*CProfile, mainchain.Addr) {
	conn, err := ethclient.Dial("ws://127.0.0.1:8546")
	tf.ChkErr(err, "failed to connect to the Ethereum")

	tf.LogBlkNum(conn)
	bal, _ := conn.BalanceAt(context.Background(), etherBaseAddr, nil)
	log.Infoln("balance is: ", bal)

	ethbasePrivKey, _ := crypto.HexToECDSA(etherBasePriv)
	etherBaseAuth := bind.NewKeyedTransactor(ethbasePrivKey)
	price := big.NewInt(2e9) // 2Gwei
	etherBaseAuth.GasPrice = price
	etherBaseAuth.GasLimit = 7000000

	client0PrivKey, _ := crypto.HexToECDSA(client0Priv)
	client0Auth := bind.NewKeyedTransactor(client0PrivKey)
	client0Auth.GasPrice = price

	ctx := context.Background()
	channelAddrBundle := deploy.DeployAll(etherBaseAuth, conn, ctx, 0)

	// Disable channel deposit limit
	tf.LogBlkNum(conn)
	ledgerContract, err := ledger.NewCelerLedger(channelAddrBundle.CelerLedgerAddr, conn)
	tf.ChkErr(err, "failed to NewCelerLedger")
	tx, err := ledgerContract.DisableBalanceLimits(etherBaseAuth)
	tf.ChkErr(err, "failed disable channel deposit limits")
	tf.WaitMinedWithChk(ctx, conn, tx, 0, "Disable balance limit")

	// Deposit into EthPool (used for openChannel)
	tf.LogBlkNum(conn)
	ethPoolContract, err := ethpool.NewEthPool(channelAddrBundle.EthPoolAddr, conn)
	tf.ChkErr(err, "failed to NewEthPool")
	amt := new(big.Int)
	amt.SetString("1000000000000000000000000", 10) // 1000000 ETH
	etherBaseAuth.Value = amt
	tx, err = ethPoolContract.Deposit(etherBaseAuth, client0Addr)
	tf.ChkErr(err, "failed to deposit into ethpool")
	tf.WaitMinedWithChk(ctx, conn, tx, 0, "Deposit to ethpool")
	etherBaseAuth.Value = big.NewInt(0)

	// Approve transferFrom of eth from ethpool for celerLedger
	tf.LogBlkNum(conn)
	tx, err = ethPoolContract.Approve(client0Auth, channelAddrBundle.CelerLedgerAddr, amt)
	tf.ChkErr(err, "failed to approve transferFrom of ETH for celerLedger")
	tf.WaitMinedWithChk(ctx, conn, tx, 0, "Approve ethpool for ledger")

	// Deploy sample ERC20 contract (CELR)
	tf.LogBlkNum(conn)
	initAmt := new(big.Int)
	initAmt.SetString("500000000000000000000000000000000000000000000", 10)
	erc20Addr, tx, erc20, err := mainchain.DeployERC20(etherBaseAuth, conn, initAmt, "Celer", 18, "CELR")
	tf.ChkErr(err, "failed to deploy ERC20")
	tf.WaitMinedWithChk(ctx, conn, tx, 0, "Deploy ERC20 "+erc20Addr.Hex())

	// Transfer ERC20 to etherbase and client0
	tf.LogBlkNum(conn)
	celrAmt := new(big.Int)
	celrAmt.SetString("500000000000000000000000000000", 10)
	addrs := []mainchain.Addr{etherBaseAddr, client0Addr}
	for _, addr := range addrs {
		tx, err = erc20.Transfer(etherBaseAuth, addr, celrAmt)
		tf.ChkErr(err, "failed to send CELR")
		mainchain.WaitMined(ctx, conn, tx, 0)
	}
	log.Infof("Sent CELR to etherbase and client0")

	// Approve transferFrom of CELR for celerLedger
	tf.LogBlkNum(conn)
	tx, err = erc20.Approve(client0Auth, channelAddrBundle.CelerLedgerAddr, celrAmt)
	tf.ChkErr(err, "failed to approve transferFrom of CELR for celerLedger")
	mainchain.WaitMined(ctx, conn, tx, 0)
	log.Infof("CELR transferFrom approved for celerLedger")

	// output json file
	p := &CProfile{
		// hardcoded values
		ETHInstance:     tf.EthInstance,
		SvrETHAddr:      client0AddrStr,
		SvrRPC:          "localhost:10000",
		ChainId:         883,
		PollingInterval: 1,
		DisputeTimeout:  10,
		// deployed addresses
		WalletAddr:       mainchain.Addr2Hex(channelAddrBundle.CelerWalletAddr),
		LedgerAddr:       mainchain.Addr2Hex(channelAddrBundle.CelerLedgerAddr),
		VirtResolverAddr: mainchain.Addr2Hex(channelAddrBundle.VirtResolverAddr),
		EthPoolAddr:      mainchain.Addr2Hex(channelAddrBundle.EthPoolAddr),
		PayResolverAddr:  mainchain.Addr2Hex(channelAddrBundle.PayResolverAddr),
		PayRegistryAddr:  mainchain.Addr2Hex(channelAddrBundle.PayRegistryAddr),
	}
	return p, erc20Addr
}
