package common

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/viper"
)

func SetupMainchain() {
	repoRoot, _ := filepath.Abs("../../..")
	log.Infoln("make localnet-down")
	cmd := exec.Command("make", "localnet-down")
	cmd.Dir = repoRoot
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	log.Infoln("build dockers, get geth, build sgn binary")
	cmd = exec.Command("make", "prepare-docker-env")
	cmd.Dir = repoRoot
	if err := cmd.Run(); err != nil {
		log.Error(err)
	}

	log.Infoln("start geth container")
	cmd = exec.Command("make", "localnet-start-geth")
	cmd.Dir = repoRoot
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	SleepWithLog(5, "geth start")

	// set up mainchain: deploy contracts, fund addrs, etc
	addrs := []mainchain.Addr{
		mainchain.Hex2Addr(ValEthAddrs[0]),
		mainchain.Hex2Addr(ValEthAddrs[1]),
		mainchain.Hex2Addr(ValEthAddrs[2]),
		mainchain.Hex2Addr(ValEthAddrs[3]),
		mainchain.Hex2Addr(DelEthAddrs[0]),
		mainchain.Hex2Addr(DelEthAddrs[1]),
		mainchain.Hex2Addr(DelEthAddrs[2]),
		mainchain.Hex2Addr(DelEthAddrs[3]),
		mainchain.Hex2Addr(ClientEthAddrs[0]),
		mainchain.Hex2Addr(ClientEthAddrs[1]),
	}
	log.Infoln("fund each test addr 100 ETH")
	err := FundAddrsETH("1"+strings.Repeat("0", 20), addrs)
	ChkErr(err, "fund each test addr 100 ETH")

	log.Infoln("set up mainchain")
	SetupEthClients()
	SetupE2eProfile()

	// fund CELR to each eth account
	log.Infoln("fund each test addr 10 million CELR")
	err = FundAddrsErc20(E2eProfile.CelrAddr, addrs, "1"+strings.Repeat("0", 25))
	ChkErr(err, "fund each test addr 10 million CELR")
}

func SetupNewSGNEnv(sgnParams *SGNParams, manual bool) {
	log.Infoln("Deploy DPoS and SGN contracts")
	if sgnParams == nil {
		sgnParams = &SGNParams{
			CelrAddr:               E2eProfile.CelrAddr,
			GovernProposalDeposit:  big.NewInt(1),
			GovernVoteTimeout:      big.NewInt(5),
			SlashTimeout:           big.NewInt(50),
			MinValidatorNum:        big.NewInt(1),
			MaxValidatorNum:        big.NewInt(7),
			MinStakingPool:         big.NewInt(100),
			AdvanceNoticePeriod:    big.NewInt(1),
			SidechainGoLiveTimeout: big.NewInt(0),
			MinGasPrices:           "0.000001quota",
		}
	}
	var tx *types.Transaction
	tx, E2eProfile.DPoSAddr, E2eProfile.SGNAddr = DeployDPoSSGNContracts(sgnParams)
	WaitMinedWithChk(context.Background(), EthClient, tx, BlockDelay, PollingInterval, "DeployDPoSSGNContracts")

	log.Infoln("make localnet-down-nodes")
	cmd := exec.Command("make", "localnet-down-nodes")
	repoRoot, _ := filepath.Abs("../../..")
	cmd.Dir = repoRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	ChkErr(err, "Failed to make localnet-down-nodes")

	log.Infoln("make prepare-sgn-data")
	cmd = exec.Command("make", "prepare-sgn-data")
	cmd.Dir = repoRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	ChkErr(err, "Failed to make prepare-sgn-data")

	log.Infoln("Updating config files of SGN nodes")
	for i := 0; i < len(ValEthKs); i++ {
		configPath := fmt.Sprintf("../../../docker-volumes/node%d/sgncli/config/sgn.toml", i)
		configFileViper := viper.New()
		configFileViper.SetConfigFile(configPath)
		err = configFileViper.ReadInConfig()
		ChkErr(err, "Failed to read config")
		configFileViper.Set(common.FlagEthCelrAddress, E2eProfile.CelrAddr.Hex())
		configFileViper.Set(common.FlagEthDPoSAddress, E2eProfile.DPoSAddr.Hex())
		configFileViper.Set(common.FlagEthLedgerAddress, E2eProfile.LedgerAddr.Hex())
		configFileViper.Set(common.FlagEthSGNAddress, E2eProfile.SGNAddr.Hex())

		if len(sgnParams.MinGasPrices) > 0 {
			configFileViper.Set(server.FlagMinGasPrices, sgnParams.MinGasPrices)
			configFileViper.Set(flags.FlagGasPrices, sgnParams.MinGasPrices)
		}
		err = configFileViper.WriteConfig()
		ChkErr(err, "Failed to write config")

		if manual {
			genesisPath := fmt.Sprintf("../../../docker-volumes/node%d/sgnd/config/genesis.json", i)
			genesisViper := viper.New()
			genesisViper.SetConfigFile(genesisPath)
			err = genesisViper.ReadInConfig()
			ChkErr(err, "Failed to read genesis")
			genesisViper.Set("app_state.govern.voting_params.voting_period", "120000000000")
			err = genesisViper.WriteConfig()
			ChkErr(err, "Failed to write genesis")
		}
	}

	// Update global viper
	node0ConfigPath := "../../../docker-volumes/node0/sgncli/config/sgn.toml"
	viper.SetConfigFile(node0ConfigPath)
	err = viper.ReadInConfig()
	ChkErr(err, "Failed to read config")
	viper.Set(common.FlagEthCelrAddress, E2eProfile.CelrAddr.Hex())
	viper.Set(common.FlagEthDPoSAddress, E2eProfile.DPoSAddr.Hex())
	viper.Set(common.FlagEthLedgerAddress, E2eProfile.LedgerAddr.Hex())
	viper.Set(common.FlagEthSGNAddress, E2eProfile.SGNAddr.Hex())

	err = SetContracts(E2eProfile.DPoSAddr, E2eProfile.SGNAddr, E2eProfile.LedgerAddr)
	ChkErr(err, "Failed to SetContracts")

	log.Infoln("make localnet-up-nodes")
	cmd = exec.Command("make", "localnet-up-nodes")
	cmd.Dir = repoRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	ChkErr(err, "Failed to make localnet-up-nodes")
}
