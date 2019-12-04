// Copyright 2018 Celer Network

package e2e

import (
	"math/big"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/celer-network/sgn/common"
	tf "github.com/celer-network/sgn/testing"
	"github.com/celer-network/goutils/log"
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

	deployGuardContract(sgnParams)

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

func updateSGNConfig() {
	log.Infoln("Updating SGN's config.json")

	viper.SetConfigFile("../../config.json")
	err := viper.ReadInConfig()
	tf.ChkErr(err, "failed to read config")

	clientKeystore, err := filepath.Abs("../keys/client0.json")
	tf.ChkErr(err, "get client keystore path")

	viper.Set(common.FlagEthWS, "ws://127.0.0.1:8546")
	viper.Set(common.FlagEthGuardAddress, e2eProfile.GuardAddr.String())
	viper.Set(common.FlagEthLedgerAddress, e2eProfile.LedgerAddr)
	viper.Set(common.FlagEthKeystore, clientKeystore)
	viper.WriteConfig()
}

func installSgn() error {
	cmd := exec.Command("make", "install")
	// set cmd.Dir under repo root path
	cmd.Dir, _ = filepath.Abs("../..")
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
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
	cmd.Dir, _ = filepath.Abs("../..")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	log.Infoln("gateway pid:", cmd.Process.Pid)
	return cmd.Process, nil
}
