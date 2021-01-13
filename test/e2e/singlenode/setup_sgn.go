package singlenode

import (
	"context"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	tc "github.com/celer-network/sgn/testing/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/viper"
)

var (
	CLIHome = os.ExpandEnv("$HOME/.sgncli")

	// root dir with ending / for all files, OutRootDirPrefix + epoch seconds
	// due to testframework etc in a different testing package, we have to define
	// same var in testframework.go and expose a set api
	outRootDir string
)

func setupNewSGNEnv(sgnParams *tc.SGNParams, testName string) []tc.Killable {
	if sgnParams == nil {
		sgnParams = &tc.SGNParams{
			CelrAddr:               tc.E2eProfile.CelrAddr,
			GovernProposalDeposit:  big.NewInt(1),
			GovernVoteTimeout:      big.NewInt(1),
			SlashTimeout:           big.NewInt(50),
			MinValidatorNum:        big.NewInt(1),
			MaxValidatorNum:        big.NewInt(11),
			MinStakingPool:         big.NewInt(100),
			AdvanceNoticePeriod:    big.NewInt(1),
			SidechainGoLiveTimeout: big.NewInt(0),
		}
	}
	var tx *types.Transaction
	tx, tc.E2eProfile.DPoSAddr, tc.E2eProfile.SGNAddr = tc.DeployDPoSSGNContracts(sgnParams)
	tc.WaitMinedWithChk(context.Background(), tc.EthClient, tx, tc.BlockDelay, tc.PollingInterval, "DeployDPoSSGNContracts")

	updateSGNConfig()

	sgnProc, err := startSidechain(outRootDir, testName)
	tc.ChkErr(err, "start sidechain")
	tc.SetContracts(tc.E2eProfile.DPoSAddr, tc.E2eProfile.SGNAddr, tc.E2eProfile.LedgerAddr)

	killable := []tc.Killable{sgnProc}
	if sgnParams.StartGateway {
		gatewayProc, err := StartGateway(outRootDir, testName)
		tc.ChkErr(err, "start gateway")
		killable = append(killable, gatewayProc)
	}

	return killable
}

func updateSGNConfig() {
	cmd := exec.Command("make", "reset-test-data")
	// set cmd.Dir under repo root path
	cmd.Dir, _ = filepath.Abs("../../..")
	err := cmd.Run()
	tc.ChkErr(err, "reset test data")

	log.Infoln("Updating genesis.json")
	genesisPath := os.ExpandEnv("$HOME/.sgnd/config/genesis.json")
	genesisViper := viper.New()
	genesisViper.SetConfigFile(genesisPath)
	err = genesisViper.ReadInConfig()
	tc.ChkErr(err, "Failed to read genesis")
	genesisViper.Set("app_state.guard.params.ledger_address", tc.E2eProfile.LedgerAddr.Hex())
	err = genesisViper.WriteConfig()
	tc.ChkErr(err, "Failed to write genesis")

	log.Infoln("Updating sgn.toml")

	configFilePath := os.ExpandEnv("$HOME/.sgncli/config/sgn.toml")
	configFileViper := viper.New()
	configFileViper.SetConfigFile(configFilePath)
	err = configFileViper.ReadInConfig()
	tc.ChkErr(err, "failed to read config")

	keystore, err := filepath.Abs("../../keys/vethks0.json")
	tc.ChkErr(err, "get keystore path")

	configFileViper.Set(common.FlagEthGateway, tc.LocalGeth)
	configFileViper.Set(common.FlagEthCelrAddress, tc.E2eProfile.CelrAddr.Hex())
	configFileViper.Set(common.FlagEthDPoSAddress, tc.E2eProfile.DPoSAddr.Hex())
	configFileViper.Set(common.FlagEthSGNAddress, tc.E2eProfile.SGNAddr.Hex())
	configFileViper.Set(common.FlagEthKeystore, keystore)
	err = configFileViper.WriteConfig()
	tc.ChkErr(err, "failed to write config")
	// Update global viper
	viper.SetConfigFile(configFilePath)
	err = viper.ReadInConfig()
	tc.ChkErr(err, "failed to read config")
}

func installSgn() error {
	cmd := exec.Command("make", "install")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "WITH_CLEVELDB=yes")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// set cmd.Dir under repo root path
	cmd.Dir, _ = filepath.Abs("../../..")
	return cmd.Run()
}

// startSidechain starts sgn sidechain with the data in test/data
func startSidechain(rootDir, testName string) (*os.Process, error) {
	cmd := exec.Command("sgnd", "start")
	cmd.Dir, _ = filepath.Abs("../../..")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.Errorln("Failed to run \"sgnd start\": ", err)
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
