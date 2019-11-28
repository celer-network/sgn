package main

import (
	"context"
	"math/big"
	"time"

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
	"github.com/spf13/viper"
)

type SGNParams struct {
	blameTimeout           *big.Int
	minValidatorNum        *big.Int
	minStakingPool         *big.Int
	sidechainGoLiveTimeout *big.Int
	startGateway           bool
}

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

const (
	// outPathPrefix is the path prefix for all output from multinode (incl. chain data, binaries etc)
	// the code will append epoch second to this and create the folder
	// the folder will be deleted after test ends successfully
	outRootDirPrefix = "/tmp/celer_multinode_"

	// etherbase and osp addr/priv key in hex
	etherBaseAddrStr  = "b5bb8b7f6f1883e0c01ffb8697024532e6f3238c"
	etherBasePriv     = "69ef4da8204644e354d759ca93b94361474259f63caac6e12d7d0abcca0063f8"
	client0AddrStr    = "6a6d2a97da1c453a4e099e8054865a0a59728863"
	client0Priv       = "a7c9fa8bcd45a86fdb5f30fecf88337f20185b0c526088f2b8e0f726cad12857"
	client0SGNAddrStr = "cosmos1ddvpnk98da5hgzz8lf5y82gnsrhvu3jd3cukpp"
	client1AddrStr    = "ba756d65a1a03f07d205749f35e2406e4a8522ad"
	client1Priv       = "c2ff7d4ce25f7448de00e21bbbb7b884bb8dc0ca642031642863e78a35cb933d"

	maxBlockDiff     = 2 // defined in sidechain's genesis file
	blockDelay       = 2
	sgnBlockInterval = 1
	defaultTimeout   = 60 * time.Second
)

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
	e2eProfile    *CProfile
	celrContract  *mainchain.ERC20
	guardAddr     mainchain.Addr
	mockCelerAddr mainchain.Addr
)

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

func deployGuardContract() mainchain.Addr {
	sgnParams := &SGNParams{
		blameTimeout:           big.NewInt(50),
		minValidatorNum:        big.NewInt(1),
		minStakingPool:         big.NewInt(100),
		sidechainGoLiveTimeout: big.NewInt(0),
	}

	conn, err := ethclient.Dial(tf.EthInstance)
	tf.ChkErr(err, "failed to connect to the Ethereum")

	ctx := context.Background()
	ethbasePrivKey, _ := crypto.HexToECDSA(etherBasePriv)
	etherBaseAuth := bind.NewKeyedTransactor(ethbasePrivKey)
	price := big.NewInt(2e9) // 2Gwei
	etherBaseAuth.GasPrice = price
	etherBaseAuth.GasLimit = 7000000

	guardAddr, tx, _, err := mainchain.DeployGuard(etherBaseAuth, conn, mockCelerAddr, sgnParams.blameTimeout, sgnParams.minValidatorNum, sgnParams.minStakingPool, sgnParams.sidechainGoLiveTimeout)
	tf.ChkErr(err, "failed to deploy Guard contract")
	tf.WaitMinedWithChk(ctx, conn, tx, 0, "Deploy Guard "+guardAddr.Hex())

	return guardAddr
}

func updateSGNConfig() {
	log.Infoln("Updating SGN's config.json")

	viper.SetConfigFile("/Users/cliu/repos/celer/sgn/docker-volumes/node0/config.json")
	err := viper.ReadInConfig()
	tf.ChkErr(err, "failed to read config")
	viper.Set(common.FlagEthGuardAddress, guardAddr.String())
	viper.Set(common.FlagEthLedgerAddress, e2eProfile.LedgerAddr)
	viper.WriteConfig()

	viper.SetConfigFile("/Users/cliu/repos/celer/sgn/docker-volumes/node0/config.json")
	err = viper.ReadInConfig()
	tf.ChkErr(err, "failed to read config")
	viper.Set(common.FlagEthGuardAddress, guardAddr.String())
	viper.Set(common.FlagEthLedgerAddress, e2eProfile.LedgerAddr)
	viper.WriteConfig()
}

func main() {
	e2eProfile, mockCelerAddr = setupMainchain()

	guardAddr = deployGuardContract()

	updateSGNConfig()
}
