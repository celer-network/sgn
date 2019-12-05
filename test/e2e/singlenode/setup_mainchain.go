package singlenode

import (
	"context"
	"math/big"
	"os"
	"os/exec"
	"strings"
	"path/filepath"

	"github.com/celer-network/cChannel-eth-go/deploy"
	"github.com/celer-network/cChannel-eth-go/ledger"
	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	tf "github.com/celer-network/sgn/testing"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// start process to handle eth rpc, and fund etherbase and server account
func startMainchain() (*os.Process, error) {
	log.Infoln("outRootDir", outRootDir, "envDir", envDir)
	chainDataDir := outRootDir + "mainchaindata"
	logFname := outRootDir + "mainchain.log"
	if err := os.MkdirAll(chainDataDir, os.ModePerm); err != nil {
		return nil, err
	}

	// geth init
	cmdInit := exec.Command("geth", "--datadir", chainDataDir, "init", "mainchain_genesis.json")
	// set cmd.Dir because relative files are under testing/env
	cmdInit.Dir, _ = filepath.Abs(envDir)
	if err := cmdInit.Run(); err != nil {
		return nil, err
	}

	// actually run geth, blocking. set syncmode full to avoid bloom mem cache by fast sync
	cmd := exec.Command("geth", "--networkid", "883", "--cache", "256", "--nousb", "--syncmode", "full", "--nodiscover", "--maxpeers", "0",
		"--netrestrict", "127.0.0.1/8", "--datadir", chainDataDir, "--keystore", "keystore", "--targetgaslimit", "8000000",
		"--ws", "--wsaddr", "localhost", "--wsport", "8546", "--wsapi", "admin,debug,eth,miner,net,personal,shh,txpool,web3",
		"--mine", "--allow-insecure-unlock", "--unlock", "0", "--password", "empty_password.txt", "--rpc", "--rpccorsdomain", "*",
		"--rpcaddr", "localhost", "--rpcport", "8545", "--rpcapi", "admin,debug,eth,miner,net,personal,shh,txpool,web3")
	cmd.Dir = cmdInit.Dir

	logF, _ := os.Create(logFname)
	cmd.Stderr = logF
	cmd.Stdout = logF
	log.Infoln("ready to start cmd")
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	log.Infoln("geth pid:", cmd.Process.Pid)
	// in case geth exits with non-zero, exit test early
	// if geth is killed by ethProc.Signal, it exits w/ 0
	go func() {
		if err := cmd.Wait(); err != nil {
			log.Errorln("geth process failed:", err)
			os.Exit(1)
		}
	}()
	return cmd.Process, nil
}

// setupMainchain deploy contracts, and do setups
// return profile, tokenAddrErc20
func setupMainchain() *TestProfile {
	conn, err := ethclient.Dial(outRootDir + "mainchaindata/geth.ipc")
	tf.ChkErr(err, "failed to connect to the Ethereum")
	ethbasePrivKey, _ := crypto.HexToECDSA(etherBasePriv)
	etherBaseAuth := bind.NewKeyedTransactor(ethbasePrivKey)

	client0PrivKey, _ := crypto.HexToECDSA(client0Priv)
	client0Auth := bind.NewKeyedTransactor(client0PrivKey)

	ctx := context.Background()
	channelAddrBundle := deploy.DeployAll(etherBaseAuth, conn, ctx, 0)

	// Disable channel deposit limit
	tf.LogBlkNum(conn)
	ledgerContract, err := ledger.NewCelerLedger(channelAddrBundle.CelerLedgerAddr, conn)
	tf.ChkErr(err, "failed to NewCelerLedger")
	tx, err := ledgerContract.DisableBalanceLimits(etherBaseAuth)
	tf.ChkErr(err, "failed disable channel deposit limits")
	tf.WaitMinedWithChk(ctx, conn, tx, 0, "Disable balance limit")

	// Deploy sample ERC20 contract (CELR)
	tf.LogBlkNum(conn)
	initAmt := new(big.Int)
	initAmt.SetString("5" + strings.Repeat("0", 44), 10)
	erc20Addr, tx, erc20, err := mainchain.DeployERC20(etherBaseAuth, conn, initAmt, "Celer", 18, "CELR")
	tf.ChkErr(err, "failed to deploy ERC20")
	tf.WaitMinedWithChk(ctx, conn, tx, 0, "Deploy ERC20 "+erc20Addr.Hex())

	// Transfer ERC20 to etherbase and client0
	tf.LogBlkNum(conn)
	celrAmt := new(big.Int)
	celrAmt.SetString("5" + strings.Repeat("0", 29), 10)
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

	return &TestProfile{
		// hardcoded values
		DisputeTimeout: 10,
		// deployed addresses
		LedgerAddr:   channelAddrBundle.CelerLedgerAddr,
		CelrAddr:     erc20Addr,
		CelrContract: erc20,
	}
}

func deployGuardContract(sgnParams *SGNParams) mainchain.Addr {
	conn, err := ethclient.Dial(tf.EthInstance)
	tf.ChkErr(err, "failed to connect to the Ethereum")

	ctx := context.Background()
	ethbasePrivKey, _ := crypto.HexToECDSA(etherBasePriv)
	etherBaseAuth := bind.NewKeyedTransactor(ethbasePrivKey)

	guardAddr, tx, _, err := mainchain.DeployGuard(etherBaseAuth, conn, e2eProfile.CelrAddr, sgnParams.blameTimeout, sgnParams.minValidatorNum, sgnParams.minStakingPool, sgnParams.sidechainGoLiveTimeout)
	tf.ChkErr(err, "failed to deploy Guard contract")
	tf.WaitMinedWithChk(ctx, conn, tx, 0, "Deploy Guard "+guardAddr.Hex())

	return guardAddr
}
