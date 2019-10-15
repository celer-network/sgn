package e2e

import (
	"context"
	"flag"
	"math/big"

	"github.com/celer-network/cChannel-eth-go/deploy"
	"github.com/celer-network/cChannel-eth-go/ethpool"
	"github.com/celer-network/cChannel-eth-go/ledger"
	log "github.com/celer-network/sgn/goceler-copy/clog"
	"github.com/celer-network/sgn/goceler-copy/common"
	"github.com/celer-network/sgn/goceler-copy/ctype"
	"github.com/celer-network/sgn/goceler-copy/utils"
	"github.com/celer-network/sgn/mainchain"
	tf "github.com/celer-network/sgn/testing"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// SetupMainchain deploy contracts, and do setups
// return profile, guardAddr, tokenAddrErc20
func SetupMainchain(appMap map[string]ctype.Addr) (*common.CProfile, string, string) {
	flag.Parse()
	conn, err := ethclient.Dial(outRootDir + "chaindata/geth.ipc")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum: %v", err)
	}
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
	ledgerContract, err := ledger.NewCelerLedger(channelAddrBundle.CelerLedgerAddr, conn)
	if err != nil {
		log.Fatal(err)
	}
	tx, err := ledgerContract.DisableBalanceLimits(etherBaseAuth)
	if err != nil {
		log.Fatalf("Failed disable channel deposit limits: %v", err)
	}
	receipt, err := utils.WaitMined(ctx, conn, tx, 0)
	if err != nil {
		log.Fatal(err)
	}
	chkTxStatus(receipt.Status, "Disable balance limit")

	// Deposit into EthPool (used for openChannel)
	ethPoolContract, err := ethpool.NewEthPool(channelAddrBundle.EthPoolAddr, conn)
	if err != nil {
		log.Fatal(err)
	}
	amt := new(big.Int)
	amt.SetString("1000000000000000000000000", 10) // 1000000 ETH
	etherBaseAuth.Value = amt
	tx, err = ethPoolContract.Deposit(etherBaseAuth, clientAddr)
	if err != nil {
		log.Fatalf("Failed to deposit into ethpool: %v", err)
	}
	receipt, err = utils.WaitMined(ctx, conn, tx, 0)
	if err != nil {
		log.Fatal(err)
	}
	etherBaseAuth.Value = big.NewInt(0)
	chkTxStatus(receipt.Status, "Deposit to ethpool")

	// Approve transferFrom of eth from ethpool for celerLedger
	tx, err = ethPoolContract.Approve(clientAuth, channelAddrBundle.CelerLedgerAddr, amt)
	if err != nil {
		log.Fatalf("Failed to approve transferFrom of ETH for celerLedger: %v", err)
	}
	receipt, err = utils.WaitMined(ctx, conn, tx, 0)
	if err != nil {
		log.Fatal(err)
	}
	chkTxStatus(receipt.Status, "Approve ethpool for ledger")

	// Deploy sample ERC20 contract (MOON)
	initAmt := new(big.Int)
	initAmt.SetString("500000000000000000000000000000000000000000000", 10)
	erc20Addr, tx, erc20, err := mainchain.DeployERC20(etherBaseAuth, conn, initAmt, "Moon", 18, "MOON")
	if err != nil {
		log.Fatalf("Failed to deploy ERC20: %v", err)
	}
	receipt, err = utils.WaitMined(ctx, conn, tx, 0)
	if err != nil {
		log.Fatal(err)
	}
	chkTxStatus(receipt.Status, "Deploy ERC20 "+erc20Addr.Hex())

	// Transfer ERC20 to etherbase and client
	moonAmt := new(big.Int)
	moonAmt.SetString("500000000000000000000000000000", 10)
	addrs := []ethcommon.Address{etherBaseAddr, clientAddr}
	for _, addr := range addrs {
		tx, err = erc20.Transfer(etherBaseAuth, addr, moonAmt)
		if err != nil {
			log.Fatalf("Failed to send MOON: %v", err)
		}
		utils.WaitMined(ctx, conn, tx, 0)
	}
	log.Infof("Sent MOON to etherbase and client")

	// Approve transferFrom of MOON for celerLedger
	tx, err = erc20.Approve(clientAuth, channelAddrBundle.CelerLedgerAddr, moonAmt)
	if err != nil {
		log.Fatalf("Failed to approve transferFrom of MOON for celerLedger: %v", err)
	}
	utils.WaitMined(ctx, conn, tx, 0)
	log.Infof("MOON transferFrom approved for celerLedger")

	// Deploy SGN Guard contract
	blameTimeout := big.NewInt(50)
	minValidatorNum := big.NewInt(1)
	minStakingPool := big.NewInt(100)
	sidechainGoLiveTimeout := big.NewInt(0)
	guardAddr, tx, _, err := mainchain.DeployGuard(etherBaseAuth, conn, erc20Addr, blameTimeout, minValidatorNum, minStakingPool, sidechainGoLiveTimeout)
	if err != nil {
		log.Fatalf("Failed to deploy Guard contract: %v", err)
	}
	receipt, err = utils.WaitMined(ctx, conn, tx, 0)
	if err != nil {
		log.Fatal(err)
	}
	chkTxStatus(receipt.Status, "Deploy Guard "+guardAddr.Hex())
	appMap["Guard"] = guardAddr

	// Deposit into EthPool client2 (used for openChannel)
	client2Addr := ctype.Hex2Addr(client2AddrStr)
	err = tf.FundAddr("100000000000000000000", []*ctype.Addr{&client2Addr})
	if err != nil {
		log.Fatalln("Failed to fund client 2", err)
	}
	client2PrivKey, _ := crypto.HexToECDSA(client2Priv)
	client2Auth := bind.NewKeyedTransactor(client2PrivKey)
	client2Auth.GasPrice = price
	amt.SetString("1000000000000000000000000", 10) // 1000000 ETH
	etherBaseAuth.Value = amt
	tx, err = ethPoolContract.Deposit(etherBaseAuth, client2Addr)
	if err != nil {
		log.Fatalf("Failed to deposit client2 into ethpool: %v", err)
	}
	receipt, err = utils.WaitMined(ctx, conn, tx, 0)
	if err != nil {
		log.Fatal(err)
	}
	etherBaseAuth.Value = big.NewInt(0)
	chkTxStatus(receipt.Status, "Deposit to ethpool client2")

	// Approve transferFrom of eth from ethpool for celerLedger
	tx, err = ethPoolContract.Approve(client2Auth, channelAddrBundle.CelerLedgerAddr, amt)
	if err != nil {
		log.Fatalf("Failed to approve client2 transferFrom of ETH for celerLedger: %v", err)
	}
	receipt, err = utils.WaitMined(ctx, conn, tx, 0)
	if err != nil {
		log.Fatal(err)
	}
	chkTxStatus(receipt.Status, "Approve ethpool for ledger, client2")

	// output json file
	p := &common.CProfile{
		// hardcoded values
		ETHInstance:     ethGateway,
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
