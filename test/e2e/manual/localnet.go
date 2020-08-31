package main

import (
	"flag"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	tc "github.com/celer-network/sgn/testing/common"
	"github.com/spf13/viper"
)

var (
	up   = flag.Bool("up", false, "start local testnet")
	down = flag.Bool("down", false, "shutdown local testnet")
)

func main() {
	flag.Parse()
	repoRoot, _ := filepath.Abs("../../..")
	if *up {
		tc.SetupMainchain()
		p := &tc.SGNParams{
			CelrAddr:               tc.E2eProfile.CelrAddr,
			GovernProposalDeposit:  big.NewInt(1000000000000000000),
			GovernVoteTimeout:      big.NewInt(30),
			SlashTimeout:           big.NewInt(15),
			MinValidatorNum:        big.NewInt(1),
			MaxValidatorNum:        big.NewInt(5),
			MinStakingPool:         big.NewInt(5000000000000000000), // 5 CELR
			IncreaseRateWaitTime:   big.NewInt(30),
			SidechainGoLiveTimeout: big.NewInt(0),
		}
		tc.SetupNewSGNEnv(p, true)

		log.Infoln("install sgnd, sgncli, sgnops")
		cmd := exec.Command("make", "install-all")
		cmd.Dir = repoRoot
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "WITH_CLEVELDB=yes")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Error(err)
		}

		log.Infoln("copy config files")
		cmd2 := exec.Command("make", "copy-manual-test-data")
		cmd2.Dir = repoRoot
		if err := cmd2.Run(); err != nil {
			log.Error(err)
		}
		log.Infoln("update config files")
		for i := 0; i < 3; i++ {
			configPath := fmt.Sprintf("./data/node%d/config.json", i)
			configFileViper := viper.New()
			configFileViper.SetConfigFile(configPath)
			if err := configFileViper.ReadInConfig(); err != nil {
				log.Error(err)
			}
			ksPath, _ := filepath.Abs(fmt.Sprintf("./data/node%d/keys/ethks%d.json", i, i))
			configFileViper.Set(common.FlagEthKeystore, ksPath)
			configFileViper.Set(common.FlagEthGateway, tc.LocalGeth)
			configFileViper.Set(common.FlagSgnNodeURI, tc.SgnNodesURI[i])

			if err := configFileViper.WriteConfig(); err != nil {
				log.Error(err)
			}
		}

	} else if *down {
		log.Infoln("Tearing down all containers...")
		cmd := exec.Command("make", "localnet-down")
		cmd.Dir = repoRoot
		if err := cmd.Run(); err != nil {
			log.Error(err)
		}
		os.Exit(0)
	}
}
