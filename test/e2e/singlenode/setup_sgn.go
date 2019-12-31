// Copyright 2018 Celer Network

package singlenode

import (
	"math/big"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	tf "github.com/celer-network/sgn/testing"
	"github.com/spf13/viper"
)

func setupNewSGNEnv(sgnParams *tf.SGNParams, testName string) []tf.Killable {
	if sgnParams == nil {
		sgnParams = &tf.SGNParams{
			BlameTimeout:           big.NewInt(50),
			MinValidatorNum:        big.NewInt(1),
			MinStakingPool:         big.NewInt(100),
			SidechainGoLiveTimeout: big.NewInt(0),
		}
	}

	sgnParams.CelrAddr = tf.E2eProfile.CelrAddr
	tf.E2eProfile.GuardAddr = tf.DeployGuardContract(sgnParams)

	updateSGNConfig()

	sgnProc, err := startSidechain(outRootDir, testName)
	tf.ChkErr(err, "start sidechain")
	tf.DefaultTestEthClient.SetContracts(tf.E2eProfile.GuardAddr.String(), tf.E2eProfile.LedgerAddr.String())

	killable := []tf.Killable{sgnProc}
	if sgnParams.StartGateway {
		gatewayProc, err := StartGateway(outRootDir, testName)
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

	clientKeystore, err := filepath.Abs("../../keys/client0.json")
	tf.ChkErr(err, "get client keystore path")

	viper.Set(common.FlagEthWS, tf.EthInstance)
	viper.Set(common.FlagEthGuardAddress, tf.E2eProfile.GuardAddr)
	viper.Set(common.FlagEthLedgerAddress, tf.E2eProfile.LedgerAddr)
	viper.Set(common.FlagEthKeystore, clientKeystore)
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
		log.Errorln("Failed to run \"make update-test-data\": ", err)
		return nil, err
	}

	cmd = exec.Command("sgn", "start")
	cmd.Dir, _ = filepath.Abs("../../..")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.Errorln("Failed to run \"sgn start\": ", err)
		return nil, err
	}

	log.Infoln("sgn pid:", cmd.Process.Pid)
	return cmd.Process, nil
}

func StartGateway(rootDir, testName string) (*os.Process, error) {
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
