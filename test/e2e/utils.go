// Copyright 2018 Celer Network

package e2e

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	ccommon "github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/ctype"
	"github.com/celer-network/sgn/flags"
	"github.com/spf13/viper"
)

// used by setup_onchain and tests
var (
	etherBaseAddr = ctype.Hex2Addr(etherBaseAddrStr)
	clientAddr    = ctype.Hex2Addr(clientAddrStr)
)

// runtime variables, will be initialized by TestMain
var (
	// root dir with ending / for all files, outRootDirPrefix + epoch seconds
	// due to testframework etc in a different testing package, we have to define
	// same var in testframework.go and expose a set api
	outRootDir     string
	envDir         = "../../testing/env"
	E2eProfile     *ccommon.CProfile
	GuardAddr      string
	Erc20TokenAddr string
)

// start process to handle eth rpc, and fund etherbase and server account
func StartMainchain() (*os.Process, error) {
	log.Println("outRootDir", outRootDir, "envDir", envDir)
	chainDataDir := outRootDir + "chaindata"
	logFname := outRootDir + "chain.log"
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
	fmt.Println("geth pid:", cmd.Process.Pid)
	// in case geth exits with non-zero, exit test early
	// if geth is killed by ethProc.Signal, it exits w/ 0
	go func() {
		if err := cmd.Wait(); err != nil {
			fmt.Println("geth process failed:", err)
			os.Exit(1)
		}
	}()
	return cmd.Process, nil
}

func UpdateSGNConfig() {
	log.Println("Updating SGN's config.json")

	viper.SetConfigFile("../../config.json")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	viper.Set(flags.FlagEthWS, "ws://127.0.0.1:8546")
	viper.Set(flags.FlagEthGuardAddress, GuardAddr)
	viper.Set(flags.FlagEthLedgerAddress, E2eProfile.LedgerAddr)
	viper.WriteConfig()
}

func sleep(second time.Duration) {
	time.Sleep(second * time.Second)
}

// StartSidechainDefault starts sgn sidechain with the data in test/data
func StartSidechainDefault(rootDir string) (*os.Process, error) {
	cmd := exec.Command("make", "update-test-data")
	// set cmd.Dir under repo root path
	cmd.Dir, _ = filepath.Abs("../..")
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	cmd = exec.Command("sgn", "start")
	cmd.Dir, _ = filepath.Abs("../..")
	logFname := rootDir + "sgn.log"
	logF, _ := os.Create(logFname)
	cmd.Stderr = logF
	cmd.Stdout = logF
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	fmt.Println("sgn pid:", cmd.Process.Pid)
	// in case sgn exits with non-zero, exit test early
	// if sgn is killed by ethProc.Signal, it exits w/ 0
	go func() {
		if err := cmd.Wait(); err != nil {
			fmt.Println("sgn process failed:", err)
			os.Exit(1)
		}
	}()
	return cmd.Process, nil
}

func installBins() error {
	cmd := exec.Command("make", "install")
	cmd.Dir, _ = filepath.Abs("../..")
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
