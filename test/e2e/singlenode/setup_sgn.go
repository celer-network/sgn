// Copyright 2018 Celer Network

package singlenode

import (
	"math/big"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	tc "github.com/celer-network/sgn/test/common"
	"github.com/spf13/viper"
)

func setupNewSGNEnv(sgnParams *tc.SGNParams, testName string) []tc.Killable {
	if sgnParams == nil {
		sgnParams = &tc.SGNParams{
			BlameTimeout:           big.NewInt(50),
			MinValidatorNum:        big.NewInt(1),
			MinStakingPool:         big.NewInt(100),
			SidechainGoLiveTimeout: big.NewInt(0),
			CelrAddr:               tc.E2eProfile.CelrAddr,
			MaxValidatorNum:        big.NewInt(11),
		}
	}
	tc.E2eProfile.GuardAddr = tc.DeployGuardContract(sgnParams)

	updateSGNConfig()

	sgnProc, err := startSidechain(outRootDir, testName)
	tc.ChkErr(err, "start sidechain")
	tc.DefaultTestEthClient.SetContracts(tc.E2eProfile.GuardAddr.String(), tc.E2eProfile.LedgerAddr.String())

	killable := []tc.Killable{sgnProc}
	if sgnParams.StartGateway {
		gatewayProc, err := StartGateway(outRootDir, testName)
		tc.ChkErr(err, "start gateway")
		killable = append(killable, gatewayProc)
	}

	return killable
}

func updateSGNConfig() {
	log.Infoln("Updating SGN's config.json")

	viper.SetConfigFile("../../../config.json")
	err := viper.ReadInConfig()
	tc.ChkErr(err, "failed to read config")

	clientKeystore, err := filepath.Abs("../../keys/ethnode0.json")
	tc.ChkErr(err, "get client keystore path")

	viper.Set(common.FlagEthWS, tc.LocalGeth)
	viper.Set(common.FlagEthGuardAddress, tc.E2eProfile.GuardAddr)
	viper.Set(common.FlagEthLedgerAddress, tc.E2eProfile.LedgerAddr)
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

	cmd = exec.Command("cp", "./test/config/local_config.json", "./config.json")
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
