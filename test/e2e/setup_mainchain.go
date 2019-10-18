package e2e

import (
	"context"
	"flag"
	"math/big"

	"github.com/celer-network/cChannel-eth-go/deploy"
	"github.com/celer-network/cChannel-eth-go/ethpool"
	"github.com/celer-network/cChannel-eth-go/ledger"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/ctype"
	"github.com/celer-network/sgn/mainchain"
	tf "github.com/celer-network/sgn/testing"
	"github.com/celer-network/sgn/testing/log"
	"github.com/celer-network/sgn/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// SetupMainchain deploy contracts, and do setups
// return profile, guardAddr, tokenAddrErc20
func SetupMainchain() (*common.CProfile, string, string) {
	flag.Parse()
	conn, err := ethclient.Dial(outRootDir + "chaindata/geth.ipc")
	chkErr(err, "failed to connect to the Ethereum")
	ethbasePrivKey, _ := crypto.HexToECDSA(etherBasePriv)
	etherBaseAuth := bind.NewKeyedTransactor(ethbasePrivKey)
	price := big.NewInt(2e9) // 2Gwei
	etherBaseAuth.GasPrice = price
	etherBaseAuth.GasLimit = 7000000

	clientPrivKey, _ := crypto.HexToECDSA(clientPriv)
	clientAuth := bind.NewKeyedTransactor(clientPrivKey)
	clientAuth.GasPrice = price

	ctx := context.Background()
	channelAddrBundle := deploy.DeployAll(etherBaseAuth, conn, ctx, 0)

	// Disable channel deposit limit
	logBlkNum(conn)
	ledgerContract, err := ledger.NewCelerLedger(channelAddrBundle.CelerLedgerAddr, conn)
	chkErr(err, "failed to NewCelerLedger")
	tx, err := ledgerContract.DisableBalanceLimits(etherBaseAuth)
	chkErr(err, "failed disable channel deposit limits")
	waitMinedWithChk(ctx, conn, tx, 0, "Disable balance limit")

	// Deposit into EthPool (used for openChannel)
	logBlkNum(conn)
	ethPoolContract, err := ethpool.NewEthPool(channelAddrBundle.EthPoolAddr, conn)
	chkErr(err, "failed to NewEthPool")
	amt := new(big.Int)
	amt.SetString("1000000000000000000000000", 10) // 1000000 ETH
	etherBaseAuth.Value = amt
	tx, err = ethPoolContract.Deposit(etherBaseAuth, clientAddr)
	chkErr(err, "failed to deposit into ethpool")
	waitMinedWithChk(ctx, conn, tx, 0, "Deposit to ethpool")
	etherBaseAuth.Value = big.NewInt(0)

	// Approve transferFrom of eth from ethpool for celerLedger
	logBlkNum(conn)
	tx, err = ethPoolContract.Approve(clientAuth, channelAddrBundle.CelerLedgerAddr, amt)
	chkErr(err, "failed to approve transferFrom of ETH for celerLedger")
	waitMinedWithChk(ctx, conn, tx, 0, "Approve ethpool for ledger")

	// Deploy sample ERC20 contract (MOON)
	logBlkNum(conn)
	initAmt := new(big.Int)
	initAmt.SetString("500000000000000000000000000000000000000000000", 10)
	erc20Addr, tx, erc20, err := mainchain.DeployERC20(etherBaseAuth, conn, initAmt, "Moon", 18, "MOON")
	chkErr(err, "failed to deploy ERC20")
	waitMinedWithChk(ctx, conn, tx, 0, "Deploy ERC20 "+erc20Addr.Hex())

	// Transfer ERC20 to etherbase and client
	logBlkNum(conn)
	moonAmt := new(big.Int)
	moonAmt.SetString("500000000000000000000000000000", 10)
	addrs := []ethcommon.Address{etherBaseAddr, clientAddr}
	for _, addr := range addrs {
		tx, err = erc20.Transfer(etherBaseAuth, addr, moonAmt)
		chkErr(err, "failed to send MOON")
		utils.WaitMined(ctx, conn, tx, 0)
	}
	log.Infof("Sent MOON to etherbase and client")

	// Approve transferFrom of MOON for celerLedger
	logBlkNum(conn)
	tx, err = erc20.Approve(clientAuth, channelAddrBundle.CelerLedgerAddr, moonAmt)
	chkErr(err, "failed to approve transferFrom of MOON for celerLedger")
	utils.WaitMined(ctx, conn, tx, 0)
	log.Infof("MOON transferFrom approved for celerLedger")

	// Deploy SGN Guard contract
	logBlkNum(conn)
	blameTimeout := big.NewInt(50)
	minValidatorNum := big.NewInt(1)
	minStakingPool := big.NewInt(100)
	sidechainGoLiveTimeout := big.NewInt(0)
	guardAddr, tx, _, err := mainchain.DeployGuard(etherBaseAuth, conn, erc20Addr, blameTimeout, minValidatorNum, minStakingPool, sidechainGoLiveTimeout)
	chkErr(err, "failed to deploy Guard contract")
	waitMinedWithChk(ctx, conn, tx, 0, "Deploy Guard "+guardAddr.Hex())

	// Deposit into EthPool client2 (used for openChannel)
	logBlkNum(conn)
	client2Addr := ctype.Hex2Addr(client2AddrStr)
	err = tf.FundAddr("100000000000000000000", []*ctype.Addr{&client2Addr})
	chkErr(err, "failed to fund client 2")
	client2PrivKey, _ := crypto.HexToECDSA(client2Priv)
	client2Auth := bind.NewKeyedTransactor(client2PrivKey)
	client2Auth.GasPrice = price
	amt.SetString("1000000000000000000000000", 10) // 1000000 ETH
	etherBaseAuth.Value = amt
	tx, err = ethPoolContract.Deposit(etherBaseAuth, client2Addr)
	chkErr(err, "failed to deposit client2 into ethpool")
	waitMinedWithChk(ctx, conn, tx, 0, "Deposit to ethpool client2")
	etherBaseAuth.Value = big.NewInt(0)

	// Approve transferFrom of eth from ethpool for celerLedger
	logBlkNum(conn)
	tx, err = ethPoolContract.Approve(client2Auth, channelAddrBundle.CelerLedgerAddr, amt)
	chkErr(err, "failed to approve client2 transferFrom of ETH for celerLedger")
	waitMinedWithChk(ctx, conn, tx, 0, "Approve ethpool for ledger, client2")

	// output json file
	p := &common.CProfile{
		// hardcoded values
		ETHInstance:     tf.EthInstance,
		SvrETHAddr:      clientAddrStr,
		SvrRPC:          "localhost:10000",
		ChainId:         883,
		PollingInterval: 1,
		DisputeTimeout:  10,
		// deployed addresses
		WalletAddr:       ctype.Addr2Hex(channelAddrBundle.CelerWalletAddr),
		LedgerAddr:       ctype.Addr2Hex(channelAddrBundle.CelerLedgerAddr),
		VirtResolverAddr: ctype.Addr2Hex(channelAddrBundle.VirtResolverAddr),
		EthPoolAddr:      ctype.Addr2Hex(channelAddrBundle.EthPoolAddr),
		PayResolverAddr:  ctype.Addr2Hex(channelAddrBundle.PayResolverAddr),
		PayRegistryAddr:  ctype.Addr2Hex(channelAddrBundle.PayRegistryAddr),
	}
	return p, ctype.Addr2Hex(guardAddr), ctype.Addr2Hex(erc20Addr)
}

// if status isn't 1 (sucess), log.Fatal
func chkTxStatus(s uint64, txname string) {
	if s != 1 {
		log.Fatal(txname + " tx failed")
	}
	log.Info(txname + " tx success")
}

func logBlkNum(conn *ethclient.Client) {
	header, err := conn.HeaderByNumber(context.Background(), nil)
	chkErr(err, "failed to get HeaderByNumber")
	log.Infoln("Latest block number on mainchain: ", header.Number)
}

func waitMinedWithChk(ctx context.Context, conn *ethclient.Client,
	tx *ethtypes.Transaction, blockDelay uint64, txname string) {
	receipt, err := utils.WaitMined(ctx, conn, tx, blockDelay)
	chkErr(err, "WaitMined error")
	chkTxStatus(receipt.Status, txname)
}
