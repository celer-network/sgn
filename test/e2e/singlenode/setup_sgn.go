// Copyright 2018 Celer Network

package singlenode

import (
	"math/big"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	tf "github.com/celer-network/sgn/testing"
	"github.com/celer-network/goutils/log"
	"github.com/spf13/viper"
	homedir "github.com/mitchellh/go-homedir"
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
	ks_path, _ := filepath.Abs("../../keys/client0.json")
	tf.SetupEthClient(ks_path)
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

	viper.SetConfigFile("../../../config.json")
	err := viper.ReadInConfig()
	tf.ChkErr(err, "failed to read config")
	viper.Set(common.FlagEthWS, "ws://127.0.0.1:8546")
	viper.Set(common.FlagEthGuardAddress, guardAddr.String())
	viper.Set(common.FlagEthLedgerAddress, e2eProfile.LedgerAddr)
	path, err := homedir.Expand("~/.sgncli")
	tf.ChkErr(err, "failed to get sgncli abs path")
	viper.Set(common.FlagSgnCLIHome, path);
	path, err = homedir.Expand("~/.sgn")
	tf.ChkErr(err, "failed to get sgn abs path")
	viper.Set(common.FlagSgnNodeHome, path);
	viper.WriteConfig()
}

func installSgn() error {
	cmd := exec.Command("make", "install")
	// set cmd.Dir under repo root path
	cmd.Dir, _ = filepath.Abs("../../..")
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// startSidechain starts sgn sidechain with the data in test/data
func startSidechain(rootDir, testName string) (*os.Process, error) {
	cmd := exec.Command("make", "update-test-data")
	// set cmd.Dir under repo root path
	cmd.Dir, _ = filepath.Abs("../../..")
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	cmd = exec.Command("sgn", "start")
	cmd.Dir, _ = filepath.Abs("../../..")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	log.Infoln("sgn pid:", cmd.Process.Pid)
	return cmd.Process, nil
}

func startGateway(rootDir, testName string) (*os.Process, error) {
	cmd := exec.Command("sgncli", "gateway")
	cmd.Dir, _ = filepath.Abs("../../..")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	log.Infoln("gateway pid:", cmd.Process.Pid)
	return cmd.Process, nil
}
