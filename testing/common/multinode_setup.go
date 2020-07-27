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

	log.Infoln("fund each validator's ETH address 100 ETH")
	addr0 := mainchain.Hex2Addr(ValEthAddrs[0])
	addr1 := mainchain.Hex2Addr(ValEthAddrs[1])
	addr2 := mainchain.Hex2Addr(ValEthAddrs[2])
	addr3 := mainchain.Hex2Addr(ClientEthAddrs[0])
	err := FundAddrsETH("1"+strings.Repeat("0", 20), []mainchain.Addr{addr0, addr1, addr2, addr3})
	ChkErr(err, "fund each validator ETH")

	log.Infoln("set up mainchain")
	SetupEthClients()
	SetupE2eProfile()

	// fund CELR to each eth account
	log.Infoln("fund each validator 10 million CELR")
	err = FundAddrsErc20(E2eProfile.CelrAddr, []mainchain.Addr{addr0, addr1, addr2, addr3}, "1"+strings.Repeat("0", 25))
	ChkErr(err, "fund each validator ERC20")
}

func SetupNewSGNEnv(sgnParams *SGNParams) {
	log.Infoln("Deploy DPoS and SGN contracts")
	if sgnParams == nil {
		sgnParams = &SGNParams{
			CelrAddr:               E2eProfile.CelrAddr,
			GovernProposalDeposit:  big.NewInt(1), // TODO: use a more practical value
			GovernVoteTimeout:      big.NewInt(3), // TODO: use a more practical value
			BlameTimeout:           big.NewInt(50),
			MinValidatorNum:        big.NewInt(1),
			MaxValidatorNum:        big.NewInt(7),
			MinStakingPool:         big.NewInt(100),
			IncreaseRateWaitTime:   big.NewInt(1), // TODO: use a more practical value
			SidechainGoLiveTimeout: big.NewInt(0),
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
	for i := 0; i < 3; i++ {
		configPath := fmt.Sprintf("../../../docker-volumes/node%d/config.json", i)
		viper.SetConfigFile(configPath)
		err = viper.ReadInConfig()
		ChkErr(err, "Failed to read config")
		viper.Set(common.FlagEthDPoSAddress, E2eProfile.DPoSAddr)
		viper.Set(common.FlagEthSGNAddress, E2eProfile.SGNAddr)
		viper.Set(common.FlagEthLedgerAddress, E2eProfile.LedgerAddr)
		err = viper.WriteConfig()
		ChkErr(err, "Failed to write config")
	}

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
