// Copyright 2018 Celer Network

package e2e

import (
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/ctype"
	"github.com/celer-network/sgn/flags"
	tf "github.com/celer-network/sgn/testing"
	"github.com/celer-network/sgn/testing/log"
	"github.com/spf13/viper"
)

type SGNParams struct {
	blameTimeout           *big.Int
	minValidatorNum        *big.Int
	minStakingPool         *big.Int
	sidechainGoLiveTimeout *big.Int
	startGateway           bool
}

// used by setup_onchain and tests
var (
	etherBaseAddr = ctype.Hex2Addr(etherBaseAddrStr)
	client0Addr   = ctype.Hex2Addr(client0AddrStr)
	client1Addr   = ctype.Hex2Addr(client1AddrStr)
)

// runtime variables, will be initialized by TestMain
var (
	// root dir with ending / for all files, outRootDirPrefix + epoch seconds
	// due to testframework etc in a different testing package, we have to define
	// same var in testframework.go and expose a set api
	outRootDir    string
	envDir        = "../../testing/env"
	E2eProfile    *common.CProfile
	GuardAddr     string
	MockCelerAddr string
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
	cmdInit := exec.Command("geth", "--datadir", chainDataDir, "init", envDir+"/mainchain_genesis.json")
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

func updateSGNConfig() {
	log.Infoln("Updating SGN's config.json")

	viper.SetConfigFile("../../config.json")
	err := viper.ReadInConfig()
	tf.ChkErr(err, "failed to read config")
	viper.Set(flags.FlagEthWS, "ws://127.0.0.1:8546")
	viper.Set(flags.FlagEthGuardAddress, GuardAddr)
	viper.Set(flags.FlagEthLedgerAddress, E2eProfile.LedgerAddr)
	viper.WriteConfig()
}

// startSidechain starts sgn sidechain with the data in test/data
func startSidechain(rootDir, testName string) (*os.Process, error) {
	cmd := exec.Command("make", "update-test-data")
	// set cmd.Dir under repo root path
	cmd.Dir, _ = filepath.Abs("../..")
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	cmd = exec.Command("sgn", "start")
	cmd.Dir, _ = filepath.Abs("../..")
	logFname := rootDir + "sgn_" + testName + ".log"
	logF, _ := os.Create(logFname)
	cmd.Stderr = logF
	cmd.Stdout = logF
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	log.Infoln("sgn pid:", cmd.Process.Pid)
	go func() {
		if err := cmd.Wait(); err != nil {
			log.Errorf("sgn process for [%s] failed: %v", testName, err)
			// os.Exit(1) // sgn is expected to be killed after each test case
		}
	}()
	return cmd.Process, nil
}

func startGateway(rootDir, testName string) (*os.Process, error) {
	cmd := exec.Command("sgncli", "gateway")
	cmd.Dir, _ = filepath.Abs("../..")
	logFname := rootDir + "gateway_" + testName + ".log"
	logF, _ := os.Create(logFname)
	cmd.Stderr = logF
	cmd.Stdout = logF
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	log.Infoln("gateway pid:", cmd.Process.Pid)
	go func() {
		if err := cmd.Wait(); err != nil {
			log.Errorf("gateway process for [%s] failed: %v", testName, err)
			// os.Exit(1) // gateway is expected to be killed after each test case
		}
	}()
	return cmd.Process, nil
}

func setupNewSGNEnv(sgnParams *SGNParams, testName string) []tf.Killable {
	if sgnParams == nil {
		sgnParams = &SGNParams{
			blameTimeout:           big.NewInt(50),
			minValidatorNum:        big.NewInt(1),
			minStakingPool:         big.NewInt(100),
			sidechainGoLiveTimeout: big.NewInt(0),
		}
	}

	GuardAddr = deployGuardContract(sgnParams)

	updateSGNConfig()
	sgnProc, err := startSidechain(outRootDir, testName)
	tf.ChkErr(err, "start sidechain")
	tf.SetupEthClient()
	tf.SetupTransactor()

	killable := []tf.Killable{sgnProc}
	if sgnParams.startGateway {
		gatewayProc, err := startGateway(outRootDir, testName)
		tf.ChkErr(err, "start gateway")
		killable = append(killable, gatewayProc)
	}

	return killable
}

func sleep(second time.Duration) {
	time.Sleep(second * time.Second)
}

func sleepWithLog(second time.Duration, waitFor string) {
	log.Infof("Sleep %d seconds for %s", second, waitFor)
	sleep(second)
}
