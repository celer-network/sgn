package multinode

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	tc "github.com/celer-network/sgn/testing/common"
	"github.com/celer-network/sgn/transactor"
	govtypes "github.com/celer-network/sgn/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupSlash() {
	log.Infoln("set up new sgn env")
	p := &tc.SGNParams{
		CelrAddr:               tc.E2eProfile.CelrAddr,
		GovernProposalDeposit:  big.NewInt(1),
		GovernVoteTimeout:      big.NewInt(1),
		SlashTimeout:           big.NewInt(10),
		MinValidatorNum:        big.NewInt(0),
		MaxValidatorNum:        big.NewInt(11),
		MinStakingPool:         big.NewInt(0),
		AdvanceNoticePeriod:    big.NewInt(1),
		SidechainGoLiveTimeout: big.NewInt(0),
	}
	tc.SetupNewSGNEnv(p, false)
	tc.SleepWithLog(10, "sgn syncing")
}

func setupValidators(t *testing.T, transactor *transactor.Transactor) {
	amts := []*big.Int{big.NewInt(8000000000000000000), big.NewInt(1000000000000000000)}
	tc.AddValidators(t, transactor, tc.ValEthKs[:2], tc.ValAccounts[:2], amts)

	_, auth, err := tc.GetAuth(tc.DelEthKs[0])
	require.NoError(t, err, "failed to get auth")
	err = tc.DelegateStake(auth, mainchain.Hex2Addr(tc.ValEthAddrs[1]), big.NewInt(1000000000000000000))
	require.NoError(t, err, "failed to delegate stake")
	_, auth, err = tc.GetAuth(tc.DelEthKs[1])
	require.NoError(t, err, "failed to get auth")
	err = tc.DelegateStake(auth, mainchain.Hex2Addr(tc.ValEthAddrs[1]), big.NewInt(1000000000000000000))
	require.NoError(t, err, "failed to delegate stake")
}

func TestE2ESlash(t *testing.T) {
	t.Run("e2e-slash", func(t *testing.T) {
		t.Run("slashTest", slashTest)
		t.Run("disableSlashTest", disableSlashTest)
		t.Run("expirePenaltyTest", expirePenaltyTest)
	})
}

// Test penalty slash when a validator is offline
func slashTest(t *testing.T) {
	log.Infoln("===================================================================")
	log.Infoln("======================== Test slash ===========================")

	setupSlash()

	transactor := tc.NewTestTransactor(
		t,
		tc.SgnCLIHomes[0],
		tc.SgnChainID,
		tc.SgnNodeURI,
		tc.SgnCLIAddr,
		tc.SgnPassphrase,
	)

	setupValidators(t, transactor)
	shutdownNode(1)
	prevBalance, _ := tc.E2eProfile.CelrContract.BalanceOf(&bind.CallOpts{}, mainchain.Hex2Addr(tc.ValEthAddrs[0]))

	log.Infoln("Query sgn about penalty info...")
	nonce := uint64(0)
	penalty, err := tc.QueryPenalty(transactor.CliCtx, nonce, 1)
	require.NoError(t, err, "failed to query penalty")
	log.Infoln("Query sgn about penalty info:", penalty.String())
	expRes1 := fmt.Sprintf(`Nonce: %d, Reason: missing_signature, ValidatorAddr: %s, TotalPenalty: 20000000000000000`, nonce, tc.ValEthAddrs[1])
	expRes2 := fmt.Sprintf(`Account: %s, Amount: 10000000000000000`, tc.ValEthAddrs[1])
	expRes3 := fmt.Sprintf(`Account: %s, Amount: 10000000000000000`, tc.DelEthAddrs[0])
	expRes4 := fmt.Sprintf(`Account: 0000000000000000000000000000000000000001, Amount: 10000000000000000`)
	assert.Equal(t, expRes1, penalty.String(), fmt.Sprintf("The expected result should be \"%s\"", expRes1))
	assert.Equal(t, expRes2, penalty.PenalizedDelegators[0].String(), fmt.Sprintf("The expected result should be \"%s\"", expRes2))
	assert.Equal(t, expRes3, penalty.PenalizedDelegators[1].String(), fmt.Sprintf("The expected result should be \"%s\"", expRes3))
	assert.Equal(t, expRes4, penalty.Beneficiaries[0].String(), fmt.Sprintf("The expected result should be \"%s\"", expRes4))

	nonce = uint64(1)
	penalty, err = tc.QueryPenalty(transactor.CliCtx, nonce, 1)
	require.NoError(t, err, "failed to query penalty")
	log.Infoln("Query sgn about penalty info:", penalty.String())
	expRes1 = fmt.Sprintf(`Nonce: %d, Reason: missing_signature, ValidatorAddr: %s, TotalPenalty: 10000000000000000`, nonce, tc.ValEthAddrs[1])
	expRes2 = fmt.Sprintf(`Account: %s, Amount: 10000000000000000`, tc.DelEthAddrs[1])
	assert.Equal(t, expRes1, penalty.String(), fmt.Sprintf("The expected result should be \"%s\"", expRes1))
	assert.Equal(t, expRes2, penalty.PenalizedDelegators[0].String(), fmt.Sprintf("The expected result should be \"%s\"", expRes2))

	log.Infoln("Query onchain staking pool")
	var poolAmt string
	for retry := 0; retry < tc.RetryLimit; retry++ {
		ci, _ := tc.DposContract.GetCandidateInfo(&bind.CallOpts{}, mainchain.Hex2Addr(tc.ValEthAddrs[1]))
		poolAmt = ci.StakingPool.String()
		if poolAmt == "2970000000000000000" {
			break
		}
		time.Sleep(tc.RetryPeriod)
	}
	assert.Equal(t, "2970000000000000000", poolAmt, fmt.Sprintf("The expected StakingPool should be 2970000000000000000"))

	log.Infoln("Query onchain validator 0 balance")
	// The validator 0 needs to submit two transactions, each transaction will reward 10000000000000000
	// so it will receive 20000000000000000 in total
	expectedBalance := new(big.Int).Add(prevBalance, big.NewInt(20000000000000000)).String()
	var balance string
	for retry := 0; retry < tc.RetryLimit; retry++ {
		b, _ := tc.E2eProfile.CelrContract.BalanceOf(&bind.CallOpts{}, mainchain.Hex2Addr(tc.ValEthAddrs[0]))
		balance = b.String()
		if balance == expectedBalance {
			break
		}
		time.Sleep(tc.RetryPeriod)
	}

	assert.Equal(t, expectedBalance, balance, fmt.Sprintf("The expected balance should be %s", expectedBalance))
}

// Test disable slash
func disableSlashTest(t *testing.T) {
	log.Infoln("===================================================================")
	log.Infoln("======================== Test disableSlash ===========================")

	setupSlash()

	transactor := tc.NewTestTransactor(
		t,
		tc.SgnCLIHomes[0],
		tc.SgnChainID,
		tc.SgnNodeURI,
		tc.ValAccounts[0],
		tc.SgnPassphrase,
	)

	setupValidators(t, transactor)

	prevBalance, _ := tc.E2eProfile.CelrContract.BalanceOf(&bind.CallOpts{}, mainchain.Hex2Addr(tc.ValEthAddrs[0]))

	paramChanges := []govtypes.ParamChange{govtypes.NewParamChange("slash", "EnableSlash", "false")}
	content := govtypes.NewParameterProposal("Slash Param Change", "Update EnableSlash", paramChanges)
	submitProposalmsg := govtypes.NewMsgSubmitProposal(content, sdk.NewInt(1), transactor.Key.GetAddress())
	transactor.AddTxMsg(submitProposalmsg)

	proposalID := uint64(1)
	byteVoteOption, _ := govtypes.VoteOptionFromString("Yes")
	voteMsg := govtypes.NewMsgVote(transactor.Key.GetAddress(), proposalID, byteVoteOption)
	transactor.AddTxMsg(voteMsg)

	_, err := tc.QueryProposal(transactor.CliCtx, proposalID, govtypes.StatusPassed)
	require.NoError(t, err, "failed to query proposal 1 with passed status")

	shutdownNode(1)

	nonce := uint64(0)
	_, err = tc.QueryPenalty(transactor.CliCtx, nonce, 1)
	assert.Error(t, err, "get penalty 0 with 1 sig should fail")

	penalty, err := tc.QueryPenalty(transactor.CliCtx, nonce, 0)
	require.NoError(t, err, "failed to query penalty 0")
	log.Infoln("Query sgn about penalty info:", penalty.String())
	expRes1 := fmt.Sprintf(`Nonce: %d, Reason: missing_signature, ValidatorAddr: %s, TotalPenalty: 20000000000000000`, nonce, tc.ValEthAddrs[1])
	expRes2 := fmt.Sprintf(`Account: %s, Amount: 10000000000000000`, tc.ValEthAddrs[1])
	expRes3 := fmt.Sprintf(`Account: %s, Amount: 10000000000000000`, tc.DelEthAddrs[0])
	expRes4 := fmt.Sprintf(`Account: 0000000000000000000000000000000000000001, Amount: 10000000000000000`)
	assert.Equal(t, expRes1, penalty.String(), fmt.Sprintf("The expected result should be \"%s\"", expRes1))
	assert.Equal(t, expRes2, penalty.PenalizedDelegators[0].String(), fmt.Sprintf("The expected result should be \"%s\"", expRes2))
	assert.Equal(t, expRes3, penalty.PenalizedDelegators[1].String(), fmt.Sprintf("The expected result should be \"%s\"", expRes3))
	assert.Equal(t, expRes4, penalty.Beneficiaries[0].String(), fmt.Sprintf("The expected result should be \"%s\"", expRes4))

	tc.SleepWithLog(30, "wait for submitting penalty")
	currentBalance, _ := tc.E2eProfile.CelrContract.BalanceOf(&bind.CallOpts{}, mainchain.Hex2Addr(tc.ValEthAddrs[0]))
	assert.Equal(t, prevBalance, currentBalance, fmt.Sprintf("The expected balance should be %s", prevBalance))
}

// Test expire penalty
func expirePenaltyTest(t *testing.T) {
	log.Infoln("===================================================================")
	log.Infoln("======================== Test expirePenalty ===========================")

	setupSlash()

	transactor := tc.NewTestTransactor(
		t,
		tc.SgnCLIHomes[0],
		tc.SgnChainID,
		tc.SgnNodeURI,
		tc.ValAccounts[0],
		tc.SgnPassphrase,
	)

	setupValidators(t, transactor)
	prevBalance, _ := tc.E2eProfile.CelrContract.BalanceOf(&bind.CallOpts{}, mainchain.Hex2Addr(tc.ValEthAddrs[0]))

	paramChanges := []govtypes.ParamChange{govtypes.NewParamChange("slash", "PenaltyTimeout", string(transactor.CliCtx.Codec.MustMarshalJSON(1)))}
	content := govtypes.NewParameterProposal("Slash Param Change", "Update PenaltyTimeout", paramChanges)
	submitProposalmsg := govtypes.NewMsgSubmitProposal(content, sdk.NewInt(1), transactor.Key.GetAddress())
	transactor.AddTxMsg(submitProposalmsg)

	proposalID := uint64(1)
	byteVoteOption, _ := govtypes.VoteOptionFromString("Yes")
	voteMsg := govtypes.NewMsgVote(transactor.Key.GetAddress(), proposalID, byteVoteOption)
	transactor.AddTxMsg(voteMsg)

	_, err := tc.QueryProposal(transactor.CliCtx, proposalID, govtypes.StatusPassed)
	require.NoError(t, err, "failed to query proposal 1 with passed status")

	shutdownNode(1)

	nonce := uint64(0)
	penalty, err := tc.QueryPenalty(transactor.CliCtx, nonce, 1)
	require.NoError(t, err, "failed to query penalty")

	log.Infoln("Query sgn about penalty info:", penalty.String())
	expRes1 := fmt.Sprintf(`Nonce: %d, Reason: missing_signature, ValidatorAddr: %s, TotalPenalty: 20000000000000000`, nonce, tc.ValEthAddrs[1])
	expRes2 := fmt.Sprintf(`Account: %s, Amount: 10000000000000000`, tc.ValEthAddrs[1])
	expRes3 := fmt.Sprintf(`Account: %s, Amount: 10000000000000000`, tc.DelEthAddrs[0])
	expRes4 := fmt.Sprintf(`Account: 0000000000000000000000000000000000000001, Amount: 10000000000000000`)
	assert.Equal(t, expRes1, penalty.String(), fmt.Sprintf("The expected result should be \"%s\"", expRes1))
	assert.Equal(t, expRes2, penalty.PenalizedDelegators[0].String(), fmt.Sprintf("The expected result should be \"%s\"", expRes2))
	assert.Equal(t, expRes3, penalty.PenalizedDelegators[1].String(), fmt.Sprintf("The expected result should be \"%s\"", expRes3))
	assert.Equal(t, expRes4, penalty.Beneficiaries[0].String(), fmt.Sprintf("The expected result should be \"%s\"", expRes4))

	tc.SleepWithLog(30, "wait for submitting penalty")
	currentBalance, _ := tc.E2eProfile.CelrContract.BalanceOf(&bind.CallOpts{}, mainchain.Hex2Addr(tc.ValEthAddrs[0]))
	assert.Equal(t, prevBalance, currentBalance, fmt.Sprintf("The expected balance should be %s", prevBalance))
}
