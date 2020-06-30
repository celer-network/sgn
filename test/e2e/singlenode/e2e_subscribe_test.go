package singlenode

import (
	"math/big"
	"strings"
	"testing"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	tc "github.com/celer-network/sgn/test/common"
	e2ecommon "github.com/celer-network/sgn/test/e2e/common"
	"github.com/spf13/viper"
)

func setUpSubscribe() []tc.Killable {
	p := &tc.SGNParams{
		CelrAddr:               tc.E2eProfile.CelrAddr,
		GovernProposalDeposit:  big.NewInt(1), // TODO: use a more practical value
		GovernVoteTimeout:      big.NewInt(1), // TODO: use a more practical value
		BlameTimeout:           big.NewInt(10),
		MinValidatorNum:        big.NewInt(0),
		MaxValidatorNum:        big.NewInt(11),
		MinStakingPool:         big.NewInt(0),
		IncreaseRateWaitTime:   big.NewInt(1), // TODO: use a more practical value
		SidechainGoLiveTimeout: big.NewInt(0),
	}
	res := setupNewSGNEnv(p, "subscribe")
	tc.SleepWithLog(10, "sgn being ready")

	return res
}

func TestE2ESubscribe(t *testing.T) {
	toKill := setUpSubscribe()
	defer tc.TearDown(toKill)

	t.Run("e2e-subscribe", func(t *testing.T) {
		t.Run("subscribeTest", subscribeTest)
	})
}

func subscribeTest(t *testing.T) {
	log.Infoln("===================================================================")
	log.Infoln("======================== Test subscribe ===========================")

	transactor := tc.NewTransactor(
		t,
		CLIHome,
		viper.GetString(common.FlagSgnChainID),
		viper.GetString(common.FlagSgnNodeURI),
		viper.GetStringSlice(common.FlagSgnTransactors)[0],
		viper.GetString(common.FlagSgnPassphrase),
	)

	amt := new(big.Int)
	amt.SetString("1"+strings.Repeat("0", 20), 10)
	ethAddr, auth, err := tc.GetAuth(tc.ValEthKs[0])
	tc.ChkTestErr(t, err, "failed to get auth")
	tc.AddCandidateWithStake(t, transactor, ethAddr, auth, tc.SgnOperators[0], amt, big.NewInt(1), big.NewInt(1), big.NewInt(10000), true)

	e2ecommon.SubscribteTestCommon(t, transactor, amt, "", 1)
}
