package multinode

import (
	"math/big"
	"testing"
	"time"

	"github.com/celer-network/goutils/log"
	tc "github.com/celer-network/sgn/testing/common"
	"github.com/celer-network/sgn/x/validator"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setUpQuota() {
	log.Infoln("Set up new sgn env")
	p := &tc.SGNParams{
		CelrAddr:               tc.E2eProfile.CelrAddr,
		GovernProposalDeposit:  big.NewInt(1),
		GovernVoteTimeout:      big.NewInt(5),
		SlashTimeout:           big.NewInt(50),
		MinValidatorNum:        big.NewInt(1),
		MaxValidatorNum:        big.NewInt(3),
		MinStakingPool:         big.NewInt(100),
		AdvanceNoticePeriod:    big.NewInt(1),
		SidechainGoLiveTimeout: big.NewInt(0),
		MinGasPrices:           "0.000001quota",
	}
	tc.SetupNewSGNEnv(p, false)
	tc.SleepWithLog(10, "sgn syncing")
}

func TestE2EQuota(t *testing.T) {
	setUpQuota()

	t.Run("e2e-quota", func(t *testing.T) {
		t.Run("quotaTest", quotaTest)
	})
}

func quotaTest(t *testing.T) {
	log.Info("=====================================================================")
	log.Info("======================== Test quota ===========================")

	transactor := tc.NewTestTransactor(
		t,
		tc.SgnCLIHomes[3],
		tc.SgnChainID,
		tc.SgnNodeURI,
		tc.ValAccounts[3],
		tc.SgnPassphrase,
	)

	amt1 := big.NewInt(3000000000000000000)
	amt2 := big.NewInt(4000000000000000000)
	amt3 := big.NewInt(3000000000000000000)
	amts := []*big.Int{amt1, amt2, amt3}
	tc.AddValidators(t, transactor, tc.ValEthKs[:3], tc.ValAccounts[:3], amts)

	ethAddr, auth, err := tc.GetAuth(tc.ValEthKs[3])
	require.NoError(t, err, "failed to get auth")
	tc.AddCandidateWithStake(
		t, transactor, ethAddr, auth, tc.ValAccounts[3],
		big.NewInt(1000000000000000000), big.NewInt(1), big.NewInt(1), big.NewInt(10000), false)

	msg := validator.NewMsgEditCandidateDescription(tc.ValEthAddrs[3], staking.Description{
		Website: "old.com",
	}, transactor.Key.GetAddress())
	for i := 0; i < 10; i++ {
		transactor.AddTxMsg(msg)
		time.Sleep(3 * time.Second)
	}

	candidate, err := validator.CLIQueryCandidate(transactor.CliCtx, validator.RouterKey, tc.ValEthAddrs[3])
	require.NoError(t, err, "failed to query candidate")
	assert.Equal(t, "old.com", candidate.Description.Website, "The expected result should be: old.com")

	msg = validator.NewMsgEditCandidateDescription(tc.ValEthAddrs[3], staking.Description{
		Website: "new.com",
	}, transactor.Key.GetAddress())
	transactor.AddTxMsg(msg)
	time.Sleep(10 * time.Second)

	candidate, err = validator.CLIQueryCandidate(transactor.CliCtx, validator.RouterKey, tc.ValEthAddrs[3])
	require.NoError(t, err, "failed to query candidate")
	assert.Equal(t, "old.com", candidate.Description.Website, "The expected result should be: old.com")

}
