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
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// SetupMainchain deploy contracts, and do setups
// return profile, tokenAddrErc20
func SetupMainchain() (*common.CProfile, string) {
	conn, err := ethclient.Dial(outRootDir + "mainchaindata/geth.ipc")
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
	addrs := []ethcommon.Address{etherBaseAddr, client0Addr}
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

func DeployGuardContract(ctx context.Context, auth *bind.TransactOpts, conn *ethclient.Client, erc20Addr ethcommon.Address, sgnParams *SGNParams) string {
	if sgnParams == nil {
		sgnParams = &SGNParams{
			blameTimeout:           big.NewInt(50),
			minValidatorNum:        big.NewInt(1),
			minStakingPool:         big.NewInt(100),
			sidechainGoLiveTimeout: big.NewInt(0),
		}
	}
	guardAddr, tx, _, err := mainchain.DeployGuard(auth, conn, erc20Addr, sgnParams.blameTimeout, sgnParams.minValidatorNum, sgnParams.minStakingPool, sgnParams.sidechainGoLiveTimeout)
	tf.ChkErr(err, "failed to deploy Guard contract")
	tf.WaitMinedWithChk(ctx, conn, tx, 0, "Deploy Guard "+guardAddr.Hex())

	return ctype.Addr2Hex(guardAddr)
}
