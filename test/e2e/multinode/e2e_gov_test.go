package multinode

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	tc "github.com/celer-network/sgn/testing/common"
	govtypes "github.com/celer-network/sgn/x/gov/types"
	"github.com/celer-network/sgn/x/validator"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupGov() {
	log.Infoln("Set up new sgn env")
	tc.SetupNewSGNEnv(nil, false)
	tc.SleepWithLog(10, "sgn syncing")
}

func TestE2EGov(t *testing.T) {
	t.Run("e2e-gov", func(t *testing.T) {
		t.Run("sidechainGovTest", sidechainGovTest)
		t.Run("mainchainGovTest", mainchainGovTest)
	})
}

func sidechainGovTest(t *testing.T) {
	log.Info("=====================================================================")
	log.Info("======================== Test sidechain gov ===========================")

	setupGov()

	transactor0 := tc.NewTestTransactor(
		t,
		tc.SgnCLIHomes[0],
		tc.SgnChainID,
		tc.SgnNodeURI,
		tc.ValAccounts[0],
		tc.SgnPassphrase,
	)

	transactor1 := tc.NewTestTransactor(
		t,
		tc.SgnCLIHomes[1],
		tc.SgnChainID,
		tc.SgnNodeURI,
		tc.ValAccounts[1],
		tc.SgnPassphrase,
	)

	transactor2 := tc.NewTestTransactor(
		t,
		tc.SgnCLIHomes[2],
		tc.SgnChainID,
		tc.SgnNodeURI,
		tc.ValAccounts[2],
		tc.SgnPassphrase,
	)

	amt1 := big.NewInt(3000000000000000000)
	amt2 := big.NewInt(2000000000000000000)
	amt3 := big.NewInt(2000000000000000000)
	amts := []*big.Int{amt1, amt2, amt3}
	tc.AddValidators(t, transactor0, tc.ValEthKs[:3], tc.ValAccounts[:3], amts)

	log.Info("======================== Test change epochlengh rejected due to small quorum ===========================")
	paramChanges := []govtypes.ParamChange{govtypes.NewParamChange("validator", "EpochLength", "\"2\"")}
	content := govtypes.NewParameterProposal("Guard Param Change", "Update EpochLength", paramChanges)
	submitProposalmsg := govtypes.NewMsgSubmitProposal(content, sdk.NewInt(1), transactor1.Key.GetAddress())
	transactor1.AddTxMsg(submitProposalmsg)

	proposalID := uint64(1)
	proposal, err := tc.QueryProposal(transactor1.CliCtx, proposalID, govtypes.StatusVotingPeriod)
	require.NoError(t, err, "failed to query proposal 1 with voting status")

	byteVoteOption, _ := govtypes.VoteOptionFromString("No")
	voteMsg := govtypes.NewMsgVote(transactor1.Key.GetAddress(), proposal.ProposalID, byteVoteOption)
	transactor1.AddTxMsg(voteMsg)

	proposal, err = tc.QueryProposal(transactor1.CliCtx, proposalID, govtypes.StatusRejected)
	require.NoError(t, err, "failed to query proposal 1 with rejected status")

	validatorParams, err := validator.CLIQueryParams(transactor1.CliCtx, validator.RouterKey)
	require.NoError(t, err, "failed to query validator params")
	assert.Equal(t, uint(3), validatorParams.EpochLength, "EpochLength params should stay 3")

	nonce := uint64(0)
	penalty, err := tc.QueryPenalty(transactor1.CliCtx, nonce, 3)
	require.NoError(t, err, "failed to query penalty 0")
	expRes1 := fmt.Sprintf(`Nonce: %d, Reason: deposit_burn, ValidatorAddr: %s, TotalPenalty: 1000000000000000000`, nonce, tc.ValEthAddrs[1])
	expRes2 := fmt.Sprintf(`Account: %s, Amount: 1000000000000000000`, tc.ValEthAddrs[1])
	assert.Equal(t, expRes1, penalty.String(), fmt.Sprintf("The expected result should be \"%s\"", expRes1))
	assert.Equal(t, expRes2, penalty.PenalizedDelegators[0].String(), fmt.Sprintf("The expected result should be \"%s\"", expRes2))

	log.Info("======================== Test change epochlengh passed for reaching quorun ===========================")
	paramChanges = []govtypes.ParamChange{govtypes.NewParamChange("validator", "EpochLength", "\"2\"")}
	content = govtypes.NewParameterProposal("Guard Param Change", "Update EpochLength", paramChanges)
	submitProposalmsg = govtypes.NewMsgSubmitProposal(content, sdk.NewInt(1), transactor0.Key.GetAddress())
	transactor0.AddTxMsg(submitProposalmsg)

	proposalID = uint64(2)
	proposal, err = tc.QueryProposal(transactor0.CliCtx, proposalID, govtypes.StatusVotingPeriod)
	require.NoError(t, err, "failed to query proposal 2 with voting status")

	byteVoteOption, _ = govtypes.VoteOptionFromString("Yes")
	voteMsg = govtypes.NewMsgVote(transactor0.Key.GetAddress(), proposal.ProposalID, byteVoteOption)
	transactor0.AddTxMsg(voteMsg)

	proposal, err = tc.QueryProposal(transactor0.CliCtx, proposalID, govtypes.StatusPassed)
	require.NoError(t, err, "failed to query proposal 2 with passed status")

	validatorParams, err = validator.CLIQueryParams(transactor0.CliCtx, validator.RouterKey)
	require.NoError(t, err, "failed to query validator params")
	assert.Equal(t, uint(2), validatorParams.EpochLength, "EpochLength params should change to 2")

	log.Info("======================== Test change epochlengh rejected due to 1/3 veto ===========================")
	paramChanges = []govtypes.ParamChange{govtypes.NewParamChange("validator", "EpochLength", "\"5\"")}
	content = govtypes.NewParameterProposal("Guard Param Change", "Update EpochLength", paramChanges)
	submitProposalmsg = govtypes.NewMsgSubmitProposal(content, sdk.NewInt(1), transactor1.Key.GetAddress())
	transactor1.AddTxMsg(submitProposalmsg)

	proposalID = uint64(3)
	proposal, err = tc.QueryProposal(transactor0.CliCtx, proposalID, govtypes.StatusVotingPeriod)
	require.NoError(t, err, "failed to query proposal 3 with voting status")

	byteVoteOption, _ = govtypes.VoteOptionFromString("NoWithVeto")
	voteMsg = govtypes.NewMsgVote(transactor0.Key.GetAddress(), proposal.ProposalID, byteVoteOption)
	transactor0.AddTxMsg(voteMsg)
	byteVoteOption, _ = govtypes.VoteOptionFromString("Yes")
	voteMsg = govtypes.NewMsgVote(transactor1.Key.GetAddress(), proposal.ProposalID, byteVoteOption)
	transactor1.AddTxMsg(voteMsg)
	voteMsg = govtypes.NewMsgVote(transactor2.Key.GetAddress(), proposal.ProposalID, byteVoteOption)
	transactor2.AddTxMsg(voteMsg)

	proposal, err = tc.QueryProposal(transactor0.CliCtx, proposalID, govtypes.StatusRejected)
	require.NoError(t, err, "failed to query proposal 3 with rejected status")

	validatorParams, err = validator.CLIQueryParams(transactor0.CliCtx, validator.RouterKey)
	require.NoError(t, err, "failed to query validator params")
	assert.Equal(t, uint(2), validatorParams.EpochLength, "EpochLength params should stay 2")

	nonce = uint64(1)
	penalty, err = tc.QueryPenalty(transactor1.CliCtx, nonce, 3)
	require.NoError(t, err, "failed to query penalty 1")
	expRes1 = fmt.Sprintf(`Nonce: %d, Reason: deposit_burn, ValidatorAddr: %s, TotalPenalty: 1000000000000000000`, nonce, tc.ValEthAddrs[1])
	assert.Equal(t, expRes1, penalty.String(), fmt.Sprintf("The expected result should be \"%s\"", expRes1))
	assert.Equal(t, expRes2, penalty.PenalizedDelegators[0].String(), fmt.Sprintf("The expected result should be \"%s\"", expRes2))

	log.Info("======================== Test change epochlengh rejected due to 1/2 No ===========================")
	paramChanges = []govtypes.ParamChange{govtypes.NewParamChange("validator", "EpochLength", "\"5\"")}
	content = govtypes.NewParameterProposal("Guard Param Change", "Update EpochLength", paramChanges)
	submitProposalmsg = govtypes.NewMsgSubmitProposal(content, sdk.NewInt(1), transactor2.Key.GetAddress())
	transactor2.AddTxMsg(submitProposalmsg)

	proposalID = uint64(4)
	proposal, err = tc.QueryProposal(transactor0.CliCtx, proposalID, govtypes.StatusVotingPeriod)
	require.NoError(t, err, "failed to query proposal 4 with voting status")

	byteVoteOption, _ = govtypes.VoteOptionFromString("No")
	voteMsg = govtypes.NewMsgVote(transactor0.Key.GetAddress(), proposal.ProposalID, byteVoteOption)
	transactor0.AddTxMsg(voteMsg)
	byteVoteOption, _ = govtypes.VoteOptionFromString("Yes")
	voteMsg = govtypes.NewMsgVote(transactor1.Key.GetAddress(), proposal.ProposalID, byteVoteOption)
	transactor1.AddTxMsg(voteMsg)

	proposal, err = tc.QueryProposal(transactor0.CliCtx, proposalID, govtypes.StatusRejected)
	require.NoError(t, err, "failed to query proposal 4 with rejected status")

	validatorParams, err = validator.CLIQueryParams(transactor0.CliCtx, validator.RouterKey)
	require.NoError(t, err, "failed to query validator params")
	assert.Equal(t, uint(2), validatorParams.EpochLength, "EpochLength params should stay 2")

	log.Info("======================== Test change epochlengh passed for over 1/2 yes ===========================")
	paramChanges = []govtypes.ParamChange{govtypes.NewParamChange("validator", "EpochLength", "\"5\"")}
	content = govtypes.NewParameterProposal("Gubscribe Param Change", "Update EpochLength", paramChanges)
	submitProposalmsg = govtypes.NewMsgSubmitProposal(content, sdk.NewInt(1), transactor2.Key.GetAddress())
	transactor2.AddTxMsg(submitProposalmsg)

	proposalID = uint64(5)
	proposal, err = tc.QueryProposal(transactor0.CliCtx, proposalID, govtypes.StatusVotingPeriod)
	require.NoError(t, err, "failed to query proposal 5 with voting status")

	byteVoteOption, _ = govtypes.VoteOptionFromString("No")
	voteMsg = govtypes.NewMsgVote(transactor2.Key.GetAddress(), proposal.ProposalID, byteVoteOption)
	transactor2.AddTxMsg(voteMsg)
	byteVoteOption, _ = govtypes.VoteOptionFromString("Yes")
	voteMsg = govtypes.NewMsgVote(transactor0.Key.GetAddress(), proposal.ProposalID, byteVoteOption)
	transactor0.AddTxMsg(voteMsg)

	proposal, err = tc.QueryProposal(transactor0.CliCtx, proposalID, govtypes.StatusPassed)
	require.NoError(t, err, "failed to query proposal 5 with passed status")

	validatorParams, err = validator.CLIQueryParams(transactor0.CliCtx, validator.RouterKey)
	require.NoError(t, err, "failed to query validator params")
	assert.Equal(t, uint(5), validatorParams.EpochLength, "EpochLength params should change to 5")
}

func mainchainGovTest(t *testing.T) {
	log.Info("=====================================================================")
	log.Info("======================== Test mainchain gov ===========================")

	setupGov()

	transactor := tc.NewTestTransactor(
		t,
		tc.SgnCLIHomes[0],
		tc.SgnChainID,
		tc.SgnNodeURI,
		tc.ValAccounts[0],
		tc.SgnPassphrase,
	)

	amt1 := big.NewInt(3000000000000000000)
	amt2 := big.NewInt(2000000000000000000)
	amt3 := big.NewInt(2000000000000000000)
	amts := []*big.Int{amt1, amt2, amt3}
	tc.AddValidators(t, transactor, tc.ValEthKs[:3], tc.ValAccounts[:3], amts)

	_, auth0, err := tc.GetAuth(tc.ValEthKs[0])
	require.NoError(t, err, "failed to get auth0")
	_, auth1, err := tc.GetAuth(tc.ValEthKs[1])
	require.NoError(t, err, "failed to get auth1")

	ctx := context.Background()
	tx, err := tc.E2eProfile.CelrContract.Approve(auth0, tc.E2eProfile.DPoSAddr, amt1)
	require.NoError(t, err, "failed to approve CELR to DPoS contract")
	tc.WaitMinedWithChk(ctx, tc.EthClient, tx, tc.BlockDelay, tc.PollingInterval, "Approve CELR to DPoS contract")

	tx, err = tc.DposContract.CreateParamProposal(auth0, big.NewInt(mainchain.MaxValidatorNum), big.NewInt(20))
	require.NoError(t, err, "failed to create param proposal")
	tc.WaitMinedWithChk(ctx, tc.EthClient, tx, tc.BlockDelay, tc.PollingInterval, "Create param proposal")

	tx, err = tc.DposContract.VoteParam(auth0, big.NewInt(0), 1)
	require.NoError(t, err, "failed to vote param proposal 0 for validator 0")
	tx, err = tc.DposContract.VoteParam(auth1, big.NewInt(0), 1)
	require.NoError(t, err, "failed to vote param proposal 0 for validator 1")
	tc.WaitMinedWithChk(ctx, tc.EthClient, tx, tc.BlockDelay, tc.PollingInterval, "Vote param proposal 0")

	time.Sleep(4 * time.Second)
	tx, err = tc.DposContract.ConfirmParamProposal(auth0, big.NewInt(0))
	require.NoError(t, err, "failed to confirm param proposal 0")
	tc.WaitMinedWithChk(ctx, tc.EthClient, tx, tc.BlockDelay, tc.PollingInterval, "Confirm param proposal 0")

	var params staking.Params
	expectedMaxValidators := uint16(25) // 20 (mainchain value) + 5 (max_validator_diff)
	for retry := 0; retry < tc.RetryLimit; retry++ {
		bz, _, err := transactor.CliCtx.Query(fmt.Sprintf("custom/%s/%s", staking.StoreKey, staking.QueryParameters))
		require.NoError(t, err, "failed to query staking params")
		transactor.CliCtx.Codec.MustUnmarshalJSON(bz, &params)
		if params.MaxValidators == expectedMaxValidators {
			break
		}
		time.Sleep(tc.RetryPeriod)
	}

	assert.Equal(t, expectedMaxValidators, params.MaxValidators, "MaxValidators param not match")
}
