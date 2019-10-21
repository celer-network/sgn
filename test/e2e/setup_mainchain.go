package e2e

import (
	"context"
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
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// SetupMainchain deploy contracts, and do setups
// return profile, tokenAddrErc20
func SetupMainchain() (*common.CProfile, string) {
	conn, err := ethclient.Dial(outRootDir + "chaindata/geth.ipc")
	tf.ChkErr(err, "failed to connect to the Ethereum")
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

	// Deploy sample ERC20 contract (MOON)
	tf.LogBlkNum(conn)
	initAmt := new(big.Int)
	initAmt.SetString("500000000000000000000000000000000000000000000", 10)
	erc20Addr, tx, erc20, err := mainchain.DeployERC20(etherBaseAuth, conn, initAmt, "Moon", 18, "MOON")
	tf.ChkErr(err, "failed to deploy ERC20")
	tf.WaitMinedWithChk(ctx, conn, tx, 0, "Deploy ERC20 "+erc20Addr.Hex())

	// Transfer ERC20 to etherbase and client0
	tf.LogBlkNum(conn)
	moonAmt := new(big.Int)
	moonAmt.SetString("500000000000000000000000000000", 10)
	addrs := []ethcommon.Address{etherBaseAddr, client0Addr}
	for _, addr := range addrs {
		tx, err = erc20.Transfer(etherBaseAuth, addr, moonAmt)
		tf.ChkErr(err, "failed to send MOON")
		utils.WaitMined(ctx, conn, tx, 0)
	}
	log.Infof("Sent MOON to etherbase and client0")

	// Approve transferFrom of MOON for celerLedger
	tf.LogBlkNum(conn)
	tx, err = erc20.Approve(client0Auth, channelAddrBundle.CelerLedgerAddr, moonAmt)
	tf.ChkErr(err, "failed to approve transferFrom of MOON for celerLedger")
	utils.WaitMined(ctx, conn, tx, 0)
	log.Infof("MOON transferFrom approved for celerLedger")

	// Deposit into EthPool client1 (used for openChannel)
	tf.LogBlkNum(conn)
	err = tf.FundAddr("100000000000000000000", []*ctype.Addr{&client1Addr})
	tf.ChkErr(err, "failed to fund client 1")
	client1PrivKey, _ := crypto.HexToECDSA(client1Priv)
	client1Auth := bind.NewKeyedTransactor(client1PrivKey)
	client1Auth.GasPrice = price
	amt.SetString("1000000000000000000000000", 10) // 1000000 ETH
	etherBaseAuth.Value = amt
	tx, err = ethPoolContract.Deposit(etherBaseAuth, client1Addr)
	tf.ChkErr(err, "failed to deposit client1 into ethpool")
	tf.WaitMinedWithChk(ctx, conn, tx, 0, "Deposit to ethpool client1")
	etherBaseAuth.Value = big.NewInt(0)

	// Approve transferFrom of eth from ethpool for celerLedger
	tf.LogBlkNum(conn)
	tx, err = ethPoolContract.Approve(client1Auth, channelAddrBundle.CelerLedgerAddr, amt)
	tf.ChkErr(err, "failed to approve client1 transferFrom of ETH for celerLedger")
	tf.WaitMinedWithChk(ctx, conn, tx, 0, "Approve ethpool for ledger, client1")

	// output json file
	p := &common.CProfile{
		// hardcoded values
		ETHInstance:     tf.EthInstance,
		SvrETHAddr:      client0AddrStr,
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
	return p, ctype.Addr2Hex(erc20Addr)
}

func DeployGuardContract(ctx context.Context, auth *bind.TransactOpts, conn *ethclient.Client, erc20Addr ethcommon.Address,
	blameTimeout, minValidatorNum, minStakingPool, sidechainGoLiveTimeout *big.Int) string {
	guardAddr, tx, _, err := mainchain.DeployGuard(auth, conn, erc20Addr, blameTimeout, minValidatorNum, minStakingPool, sidechainGoLiveTimeout)
	tf.ChkErr(err, "failed to deploy Guard contract")
	tf.WaitMinedWithChk(ctx, conn, tx, 0, "Deploy Guard "+guardAddr.Hex())

	return ctype.Addr2Hex(guardAddr)
}