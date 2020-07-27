package main

import (
	"flag"
	"math/big"

	tc "github.com/celer-network/sgn/testing/common"
)

func main() {
	flag.Parse()
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
}
