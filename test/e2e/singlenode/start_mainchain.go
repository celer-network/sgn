package singlenode

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/celer-network/goutils/log"
	tc "github.com/celer-network/sgn/testing/common"
)

// start process to handle eth rpc, and fund etherbase and server account
func startMainchain() (*os.Process, error) {
	log.Infoln("outRootDir", outRootDir, "envDir", tc.EnvDir)
	chainDataDir := outRootDir + "mainchaindata"
	logFname := outRootDir + "mainchain.log"
	if err := os.MkdirAll(chainDataDir, os.ModePerm); err != nil {
		return nil, err
	}

	// geth init
	cmdInit := exec.Command("geth", "--datadir", chainDataDir, "init", "mainchain_genesis.json")
	// set cmd.Dir because relative files are under testing/env
	cmdInit.Dir, _ = filepath.Abs(tc.EnvDir)
	if err := cmdInit.Run(); err != nil {
		return nil, err
	}

	// actually run geth, blocking. set syncmode full to avoid bloom mem cache by fast sync
	cmd := exec.Command("geth", "--networkid", "883", "--cache", "256", "--nousb", "--syncmode", "full", "--nodiscover", "--maxpeers", "0",
		"--netrestrict", "127.0.0.1/8", "--datadir", chainDataDir, "--keystore", "keystore", "--miner.gastarget", "8000000",
		"--ws", "--ws.addr", "localhost", "--ws.port", "8546", "--ws.api", "admin,debug,eth,miner,net,personal,shh,txpool,web3",
		"--mine", "--allow-insecure-unlock", "--unlock", "0xb5BB8b7f6f1883e0c01ffb8697024532e6F3238C", "--password", "empty_password.txt",
		"--http", "--http.corsdomain", "*", "--http.addr", "localhost", "--http.port", "8545", "--http.api",
		"admin,debug,eth,miner,net,personal,shh,txpool,web3")
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
