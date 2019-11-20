// Copyright 2018 Celer Network

package e2e

import (
	"math/big"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/celer-network/sgn/flags"
	"github.com/celer-network/sgn/mainchain"
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

func setupNewSGNEnv(sgnParams *SGNParams, testName string) []tf.Killable {
	if sgnParams == nil {
		sgnParams = &SGNParams{
			blameTimeout:           big.NewInt(50),
			minValidatorNum:        big.NewInt(1),
			minStakingPool:         big.NewInt(100),
			sidechainGoLiveTimeout: big.NewInt(0),
		}
	}

	guardAddr = deployGuardContract(sgnParams)

	updateSGNConfig()
	sgnProc, err := startSidechain(outRootDir, testName)
	tf.ChkErr(err, "start sidechain")
	tf.SetupEthClient()
	tf.SetupTransactor()

	celrContract, err = mainchain.NewERC20(mockCelerAddr, tf.EthClient.Client)
	tf.ChkErr(err, "NewERC20 error")

	killable := []tf.Killable{sgnProc}
	if sgnParams.startGateway {
		gatewayProc, err := startGateway(outRootDir, testName)
		tf.ChkErr(err, "start gateway")
		killable = append(killable, gatewayProc)
	}

	return killable
}

func updateSGNConfig() {
	log.Infoln("Updating SGN's config.json")

	viper.SetConfigFile("../../config.json")
	err := viper.ReadInConfig()
	tf.ChkErr(err, "failed to read config")
	viper.Set(flags.FlagEthWS, "ws://127.0.0.1:8546")
	viper.Set(flags.FlagEthGuardAddress, guardAddr.String())
	viper.Set(flags.FlagEthLedgerAddress, e2eProfile.LedgerAddr)
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
