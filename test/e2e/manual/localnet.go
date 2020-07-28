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
			GovernProposalDeposit:  big.NewInt(1),
			GovernVoteTimeout:      big.NewInt(1),
			BlameTimeout:           big.NewInt(0),
			MinValidatorNum:        big.NewInt(1),
			MaxValidatorNum:        big.NewInt(5),
			MinStakingPool:         big.NewInt(1),
			IncreaseRateWaitTime:   big.NewInt(1),
			SidechainGoLiveTimeout: big.NewInt(0),
		}
		tc.SetupNewSGNEnv(p)
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
