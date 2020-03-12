package multinode

import (
	"math/big"
	"testing"

	"github.com/celer-network/goutils/log"
	tc "github.com/celer-network/sgn/test/common"
	"github.com/celer-network/sgn/x/global"
	govtypes "github.com/celer-network/sgn/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func setUpGov() {
	log.Infoln("Set up new sgn env")
	setupNewSGNEnv(nil)
	tc.SleepWithLog(10, "sgn syncing")
}

func TestE2EGov(t *testing.T) {
	setUpGov()

	t.Run("e2e-gov", func(t *testing.T) {
		t.Run("govTest", govTest)
	})
}

func govTest(t *testing.T) {
	log.Info("=====================================================================")
	log.Info("======================== Test gov ===========================")

	transactor0 := tc.NewTransactor(
		t,
		tc.SgnCLIHomes[0],
		tc.SgnChainID,
		tc.SgnNodeURI,
		tc.SgnOperators[0],
		tc.SgnPassphrase,
	)

	transactor1 := tc.NewTransactor(
		t,
		tc.SgnCLIHomes[1],
		tc.SgnChainID,
		tc.SgnNodeURI,
		tc.SgnOperators[1],
		tc.SgnPassphrase,
	)

	transactor2 := tc.NewTransactor(
		t,
		tc.SgnCLIHomes[2],
		tc.SgnChainID,
		tc.SgnNodeURI,
		tc.SgnOperators[2],
		tc.SgnPassphrase,
	)

	amt1 := big.NewInt(3000000000000000000)
	amt2 := big.NewInt(2000000000000000000)
	amt3 := big.NewInt(2000000000000000000)
	amts := []*big.Int{amt1, amt2, amt3}
	tc.AddValidators(t, transactor0, tc.ValEthKs[:], tc.SgnOperators[:], amts)

	log.Info("======================== Test change epochlengh rejected due to small quorum ===========================")
	paramChanges := []govtypes.ParamChange{govtypes.NewParamChange("global", "EpochLength", "\"3\"")}
	content := govtypes.NewParameterProposal("Global Param Change", "Update EpochLength", paramChanges)
	submitProposalmsg := govtypes.NewMsgSubmitProposal(content, sdk.NewInt(1), transactor1.Key.GetAddress())
	transactor1.AddTxMsg(submitProposalmsg)

	proposalID := uint64(1)
	proposal, err := tc.QueryProposal(transactor1.CliCtx, proposalID, govtypes.StatusVotingPeriod)
	tc.ChkTestErr(t, err, "failed to query proposal 1 with voting status")

	byteVoteOption, _ := govtypes.VoteOptionFromString("No")
	voteMsg := govtypes.NewMsgVote(transactor1.Key.GetAddress(), proposal.ProposalID, byteVoteOption)
	transactor1.AddTxMsg(voteMsg)

	proposal, err = tc.QueryProposal(transactor1.CliCtx, proposalID, govtypes.StatusRejected)
	tc.ChkTestErr(t, err, "failed to query proposal 1 with rejected status")

	globalParams, err := global.CLIQueryParams(transactor1.CliCtx, global.RouterKey)
	tc.ChkTestErr(t, err, "failed to query global params")
	assert.Equal(t, int64(1), globalParams.EpochLength, "EpochLength params should stay 1")

	transactor1.AddTxMsg(submitProposalmsg)
	proposalID = uint64(2)
	proposal, err = tc.QueryProposal(transactor1.CliCtx, proposalID, govtypes.StatusVotingPeriod)
	assert.Error(t, err, "fail to submit proposal due to muted depositor")

	log.Info("======================== Test change epochlengh passed for reaching quorun ===========================")
	paramChanges = []govtypes.ParamChange{govtypes.NewParamChange("global", "EpochLength", "\"3\"")}
	content = govtypes.NewParameterProposal("Global Param Change", "Update EpochLength", paramChanges)
	submitProposalmsg = govtypes.NewMsgSubmitProposal(content, sdk.NewInt(1), transactor0.Key.GetAddress())
	transactor0.AddTxMsg(submitProposalmsg)

	proposalID = uint64(2)
	proposal, err = tc.QueryProposal(transactor0.CliCtx, proposalID, govtypes.StatusVotingPeriod)
	tc.ChkTestErr(t, err, "failed to query proposal 2 with voting status")

	byteVoteOption, _ = govtypes.VoteOptionFromString("Yes")
	voteMsg = govtypes.NewMsgVote(transactor0.Key.GetAddress(), proposal.ProposalID, byteVoteOption)
	transactor0.AddTxMsg(voteMsg)

	proposal, err = tc.QueryProposal(transactor0.CliCtx, proposalID, govtypes.StatusPassed)
	tc.ChkTestErr(t, err, "failed to query proposal 2 with passed status")

	globalParams, err = global.CLIQueryParams(transactor0.CliCtx, global.RouterKey)
	tc.ChkTestErr(t, err, "failed to query global params")
	assert.Equal(t, int64(3), globalParams.EpochLength, "EpochLength params should change to 3")

	log.Info("======================== Test change epochlengh rejected due to 1/3 veto ===========================")
	paramChanges = []govtypes.ParamChange{govtypes.NewParamChange("global", "EpochLength", "\"5\"")}
	content = govtypes.NewParameterProposal("Global Param Change", "Update EpochLength", paramChanges)
	submitProposalmsg = govtypes.NewMsgSubmitProposal(content, sdk.NewInt(1), transactor0.Key.GetAddress())
	transactor1.AddTxMsg(submitProposalmsg)

	proposalID = uint64(3)
	proposal, err = tc.QueryProposal(transactor0.CliCtx, proposalID, govtypes.StatusVotingPeriod)
	tc.ChkTestErr(t, err, "failed to query proposal 3 with voting status")

	byteVoteOption, _ = govtypes.VoteOptionFromString("NoWithVeto")
	voteMsg = govtypes.NewMsgVote(transactor0.Key.GetAddress(), proposal.ProposalID, byteVoteOption)
	transactor0.AddTxMsg(voteMsg)
	byteVoteOption, _ = govtypes.VoteOptionFromString("Yes")
	voteMsg = govtypes.NewMsgVote(transactor1.Key.GetAddress(), proposal.ProposalID, byteVoteOption)
	transactor1.AddTxMsg(voteMsg)
	voteMsg = govtypes.NewMsgVote(transactor2.Key.GetAddress(), proposal.ProposalID, byteVoteOption)
	transactor2.AddTxMsg(voteMsg)

	proposal, err = tc.QueryProposal(transactor0.CliCtx, proposalID, govtypes.StatusRejected)
	tc.ChkTestErr(t, err, "failed to query proposal 3 with rejected status")

	globalParams, err = global.CLIQueryParams(transactor0.CliCtx, global.RouterKey)
	tc.ChkTestErr(t, err, "failed to query global params")
	assert.Equal(t, int64(3), globalParams.EpochLength, "EpochLength params should stay 3")

	transactor1.AddTxMsg(submitProposalmsg)
	proposalID = uint64(4)
	proposal, err = tc.QueryProposal(transactor1.CliCtx, proposalID, govtypes.StatusVotingPeriod)
	assert.Error(t, err, "fail to submit proposal due to muted depositor")

	log.Info("======================== Test change epochlengh rejected due to 1/2 No ===========================")
	paramChanges = []govtypes.ParamChange{govtypes.NewParamChange("global", "EpochLength", "\"5\"")}
	content = govtypes.NewParameterProposal("Global Param Change", "Update EpochLength", paramChanges)
	submitProposalmsg = govtypes.NewMsgSubmitProposal(content, sdk.NewInt(1), transactor0.Key.GetAddress())
	transactor2.AddTxMsg(submitProposalmsg)

	proposalID = uint64(4)
	proposal, err = tc.QueryProposal(transactor0.CliCtx, proposalID, govtypes.StatusVotingPeriod)
	tc.ChkTestErr(t, err, "failed to query proposal 4 with voting status")

	byteVoteOption, _ = govtypes.VoteOptionFromString("No")
	voteMsg = govtypes.NewMsgVote(transactor0.Key.GetAddress(), proposal.ProposalID, byteVoteOption)
	transactor0.AddTxMsg(voteMsg)
	byteVoteOption, _ = govtypes.VoteOptionFromString("Yes")
	voteMsg = govtypes.NewMsgVote(transactor1.Key.GetAddress(), proposal.ProposalID, byteVoteOption)
	transactor1.AddTxMsg(voteMsg)

	proposal, err = tc.QueryProposal(transactor0.CliCtx, proposalID, govtypes.StatusRejected)
	tc.ChkTestErr(t, err, "failed to query proposal 4 with rejected status")

	globalParams, err = global.CLIQueryParams(transactor0.CliCtx, global.RouterKey)
	tc.ChkTestErr(t, err, "failed to query global params")
	assert.Equal(t, int64(3), globalParams.EpochLength, "EpochLength params should stay 3")

	log.Info("======================== Test change epochlengh passed for over 1/2 yes ===========================")
	paramChanges = []govtypes.ParamChange{govtypes.NewParamChange("global", "EpochLength", "\"5\"")}
	content = govtypes.NewParameterProposal("Global Param Change", "Update EpochLength", paramChanges)
	submitProposalmsg = govtypes.NewMsgSubmitProposal(content, sdk.NewInt(1), transactor0.Key.GetAddress())
	transactor2.AddTxMsg(submitProposalmsg)

	proposalID = uint64(5)
	proposal, err = tc.QueryProposal(transactor0.CliCtx, proposalID, govtypes.StatusVotingPeriod)
	tc.ChkTestErr(t, err, "failed to query proposal 5 with voting status")

	byteVoteOption, _ = govtypes.VoteOptionFromString("No")
	voteMsg = govtypes.NewMsgVote(transactor0.Key.GetAddress(), proposal.ProposalID, byteVoteOption)
	transactor0.AddTxMsg(voteMsg)
	byteVoteOption, _ = govtypes.VoteOptionFromString("Yes")
	voteMsg = govtypes.NewMsgVote(transactor1.Key.GetAddress(), proposal.ProposalID, byteVoteOption)
	transactor1.AddTxMsg(voteMsg)
	voteMsg = govtypes.NewMsgVote(transactor2.Key.GetAddress(), proposal.ProposalID, byteVoteOption)
	transactor2.AddTxMsg(voteMsg)

	proposal, err = tc.QueryProposal(transactor0.CliCtx, proposalID, govtypes.StatusRejected)
	tc.ChkTestErr(t, err, "failed to query proposal 5 with rejected status")

	globalParams, err = global.CLIQueryParams(transactor0.CliCtx, global.RouterKey)
	tc.ChkTestErr(t, err, "failed to query global params")
	assert.Equal(t, int64(5), globalParams.EpochLength, "EpochLength params should stay 5")

}
