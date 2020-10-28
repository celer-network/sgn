package multinode

import (
	"math/big"
	"testing"

	"github.com/celer-network/goutils/log"
	tc "github.com/celer-network/sgn/testing/common"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func setupValidator(maxValidatorNum *big.Int) {
	log.Infoln("set up new sgn env")
	p := &tc.SGNParams{
		CelrAddr:               tc.E2eProfile.CelrAddr,
		GovernProposalDeposit:  big.NewInt(1),
		GovernVoteTimeout:      big.NewInt(1),
		SlashTimeout:           big.NewInt(0),
		MinValidatorNum:        big.NewInt(1),
		MaxValidatorNum:        maxValidatorNum,
		MinStakingPool:         big.NewInt(1),
		AdvanceNoticePeriod:    big.NewInt(1),
		SidechainGoLiveTimeout: big.NewInt(0),
	}
	tc.SetupNewSGNEnv(p, false)
	tc.SleepWithLog(10, "sgn being ready")
}

func TestE2EValidator(t *testing.T) {
	t.Run("e2e-validator", func(t *testing.T) {
		t.Run("validatorTest", validatorTest)
	})
}

func validatorTest(t *testing.T) {
	log.Info("===================================================================")
	log.Info("======================== Test validator ===========================")
	setupValidator(big.NewInt(3))

	transactor := tc.NewTestTransactor(
		t,
		tc.SgnCLIHomes[0],
		tc.SgnChainID,
		tc.SgnNodeURI,
		tc.SgnCLIAddr,
		tc.SgnPassphrase,
	)

	amts := []*big.Int{
		big.NewInt(8000000000000000000),
		big.NewInt(5000000000000000000),
		big.NewInt(4000000000000000000),
		big.NewInt(6000000000000000000),
	}
	minAmts := []*big.Int{
		big.NewInt(4000000000000000000),
		big.NewInt(2000000000000000000),
		big.NewInt(2000000000000000000),
		big.NewInt(2000000000000000000),
	}
	commissionRate := big.NewInt(200)

	log.Infoln("---------- It should add two validators successfully ----------")
	for i := 0; i < 2; i++ {
		log.Infoln("Adding validator", i)
		ethAddr, auth, err := tc.GetAuth(tc.ValEthKs[i])
		require.NoError(t, err, "failed to get auth")
		tc.AddCandidateWithStake(
			t, transactor, ethAddr, auth, tc.ValAccounts[i],
			amts[i], minAmts[i], commissionRate, big.NewInt(10000), true)
		tc.CheckValidatorNum(t, transactor, i+1)
	}

	log.Infoln("---------- It should fail to add validator 2 without enough delegation ----------")
	ethAddr, auth, err := tc.GetAuth(tc.ValEthKs[2])
	require.NoError(t, err, "failed to get auth")
	initialDelegation := big.NewInt(1000000000000000000)
	tc.AddCandidateWithStake(
		t, transactor, ethAddr, auth, tc.ValAccounts[2],
		initialDelegation, minAmts[2], commissionRate, big.NewInt(10000), false)
	log.Info("Query sgn about validators to check if validator 2 is not added...")
	tc.CheckValidatorNum(t, transactor, 2)

	log.Infoln("---------- It should correctly add validator 2 with enough delegation ----------")
	err = tc.DelegateStake(auth, ethAddr, big.NewInt(0).Sub(amts[2], initialDelegation))
	require.NoError(t, err, "failed to delegate stake")
	tc.CheckValidatorNum(t, transactor, 3)
	tc.CheckValidator(t, transactor, tc.ValAccounts[2], amts[2], sdk.Bonded)

	log.Infoln("---------- It should successfully remove validator 2 caused by intendWithdraw ----------")
	err = tc.IntendWithdraw(auth, ethAddr, amts[2])
	require.NoError(t, err, "failed to intendWithdraw stake")
	log.Info("Query sgn about the validators to check if it has correct number of validators...")
	tc.CheckValidatorNum(t, transactor, 2)
	tc.CheckValidatorStatus(t, transactor, tc.ValAccounts[2], sdk.Unbonding)

	err = tc.ConfirmUnbondedCandidate(auth, ethAddr)
	require.NoError(t, err, "failed to confirmUnbondedCandidate")
	tc.CheckCandidate(t, transactor, ethAddr, tc.ValAccounts[2], big.NewInt(0))

	log.Infoln("---------- It should successfully add back validator 2 with enough delegation ----------")
	err = tc.DelegateStake(auth, ethAddr, amts[2])
	require.NoError(t, err, "failed to delegate stake")
	tc.CheckValidatorNum(t, transactor, 3)
	tc.CheckValidator(t, transactor, tc.ValAccounts[2], amts[2], sdk.Bonded)

	log.Infoln("---------- It should correctly replace validator 2 with validator 3 ----------")
	ethAddr, auth, err = tc.GetAuth(tc.ValEthKs[3])
	require.NoError(t, err, "failed to get auth")
	tc.AddCandidateWithStake(
		t, transactor, ethAddr, auth, tc.ValAccounts[3],
		amts[3], minAmts[3], commissionRate, big.NewInt(10000), true)
	tc.CheckValidatorNum(t, transactor, 3)
	tc.CheckValidator(t, transactor, tc.ValAccounts[2], amts[2], sdk.Unbonding)
}
