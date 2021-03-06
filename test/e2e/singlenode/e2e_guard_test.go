package singlenode

import (
	"math/big"
	"strings"
	"testing"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	tec "github.com/celer-network/sgn/test/e2e/common"
	tc "github.com/celer-network/sgn/testing/common"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func setupGuard() []tc.Killable {
	p := &tc.SGNParams{
		CelrAddr:               tc.E2eProfile.CelrAddr,
		GovernProposalDeposit:  big.NewInt(1),
		GovernVoteTimeout:      big.NewInt(5),
		SlashTimeout:           big.NewInt(10),
		MinValidatorNum:        big.NewInt(0),
		MaxValidatorNum:        big.NewInt(11),
		MinStakingPool:         big.NewInt(0),
		AdvanceNoticePeriod:    big.NewInt(1),
		SidechainGoLiveTimeout: big.NewInt(0),
	}
	res := setupNewSGNEnv(p, "guard")
	tc.SleepWithLog(10, "sgn being ready")

	return res
}

func TestE2EGuard(t *testing.T) {
	toKill := setupGuard()
	defer tc.TearDown(toKill)

	t.Run("e2e-guard", func(t *testing.T) {
		t.Run("guardTest", guardTest)
	})
}

func guardTest(t *testing.T) {
	log.Infoln("===================================================================")
	log.Infoln("======================== Test guard ===========================")

	transactor := tc.NewTestTransactor(
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
	require.NoError(t, err, "failed to get auth")
	tc.AddCandidateWithStake(
		t, transactor, ethAddr, auth, tc.ValAccounts[0], amt, big.NewInt(1), big.NewInt(1), big.NewInt(10000), true)

	tec.Subscribe(t, transactor, amt)
	tec.TestGuard(t, transactor, []string{tc.ValEthAddrs[0], tc.ValEthAddrs[0]})
	// 2 requests, pre request cost is 1000000000000000000, all goes to validator
	tec.CheckReward(t, transactor, tc.ValEthAddrs[0], "2000000000000000000", 1)
}
