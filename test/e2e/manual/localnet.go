package main

import (
	"flag"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/celer-network/goutils/log"
	tc "github.com/celer-network/sgn/testing/common"
)

var (
	up   = flag.Bool("up", false, "start local testnet")
	down = flag.Bool("down", false, "shutdown local testnet")
)

func main() {
	flag.Parse()
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
	} else if *down {
		repoRoot, _ := filepath.Abs("../../..")
		log.Infoln("Tearing down all containers...")
		cmd := exec.Command("make", "localnet-down")
		cmd.Dir = repoRoot
		if err := cmd.Run(); err != nil {
			log.Error(err)
		}
		os.Exit(0)
	}
}
