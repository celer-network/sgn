package singlenode

import (
	"testing"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	tc "github.com/celer-network/sgn/test/common"
	"github.com/celer-network/sgn/x/global"
	govtypes "github.com/celer-network/sgn/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func setUpGov() []tc.Killable {
	res := setupNewSGNEnv(nil, "gov change parameter")
	tc.SleepWithLog(10, "sgn syncing")

	return res
}

// func TestE2EGov(t *testing.T) {
// 	toKill := setUpGov()
// 	defer tc.TearDown(toKill)

// 	t.Run("e2e-gov", func(t *testing.T) {
// 		t.Run("govTest", govTest)
// 	})
// }

func govTest(t *testing.T) {
	log.Info("=====================================================================")
	log.Info("======================== Test gov ===========================")

	transactor := tc.NewTransactor(
		t,
		CLIHome,
		viper.GetString(common.FlagSgnChainID),
		viper.GetString(common.FlagSgnNodeURI),
		viper.GetString(common.FlagSgnOperator),
		viper.GetString(common.FlagSgnPassphrase),
	)

	log.Info("======================== Test change epochlengh passed ===========================")
	paramChanges := []govtypes.ParamChange{govtypes.NewParamChange("global", "EpochLength", "\"3\"")}
	content := govtypes.NewParameterProposal("Global Param Change", "Update EpochLength", paramChanges)
	submitProposalmsg := govtypes.NewMsgSubmitProposal(content, sdk.NewInt(0), transactor.Key.GetAddress())
	transactor.AddTxMsg(submitProposalmsg)

	proposalID := uint64(1)
	proposal, err := tc.QueryProposal(transactor.CliCtx, proposalID, govtypes.StatusDepositPeriod)
	tc.ChkTestErr(t, err, "failed to query proposal 1 with deposit status")
	assert.Equal(t, content.GetTitle(), proposal.GetTitle(), "The proposal should have same title as submitted proposal")
	assert.Equal(t, content.GetDescription(), proposal.GetDescription(), "The proposal should have same description as submitted proposal")

	depositMsg := govtypes.NewMsgDeposit(transactor.Key.GetAddress(), proposalID, sdk.NewInt(10))
	transactor.AddTxMsg(depositMsg)
	proposal, err = tc.QueryProposal(transactor.CliCtx, proposalID, govtypes.StatusVotingPeriod)
	tc.ChkTestErr(t, err, "failed to query proposal 1 with voting status")

	byteVoteOption, _ := govtypes.VoteOptionFromString("Yes")
	voteMsg := govtypes.NewMsgVote(transactor.Key.GetAddress(), proposal.ProposalID, byteVoteOption)
	transactor.AddTxMsg(voteMsg)

	proposal, err = tc.QueryProposal(transactor.CliCtx, proposalID, govtypes.StatusPassed)
	tc.ChkTestErr(t, err, "failed to query proposal 1 with passed status")

	globalParams, err := global.CLIQueryParams(transactor.CliCtx, global.RouterKey)
	tc.ChkTestErr(t, err, "failed to query global params")
	assert.Equal(t, int64(3), globalParams.EpochLength, "EpochLength params should be updated to 3")

	log.Info("======================== Test change epochlengh rejected ===========================")
	paramChanges = []govtypes.ParamChange{govtypes.NewParamChange("global", "EpochLength", "\"5\"")}
	content = govtypes.NewParameterProposal("Global Param Change", "Update EpochLength", paramChanges)
	submitProposalmsg = govtypes.NewMsgSubmitProposal(content, sdk.NewInt(10), transactor.Key.GetAddress())
	transactor.AddTxMsg(submitProposalmsg)

	proposalID = uint64(2)
	proposal, err = tc.QueryProposal(transactor.CliCtx, proposalID, govtypes.StatusVotingPeriod)
	tc.ChkTestErr(t, err, "failed to query proposal 2 with voting status")

	byteVoteOption, _ = govtypes.VoteOptionFromString("NoWithVeto")
	voteMsg = govtypes.NewMsgVote(transactor.Key.GetAddress(), proposal.ProposalID, byteVoteOption)
	transactor.AddTxMsg(voteMsg)

	proposal, err = tc.QueryProposal(transactor.CliCtx, proposalID, govtypes.StatusRejected)
	tc.ChkTestErr(t, err, "failed to query proposal 2 with rejected status")

	globalParams, err = global.CLIQueryParams(transactor.CliCtx, global.RouterKey)
	tc.ChkTestErr(t, err, "failed to query global params")
	assert.Equal(t, int64(3), globalParams.EpochLength, "EpochLength params should stay 3")

	transactor.AddTxMsg(submitProposalmsg)
	proposalID = uint64(3)
	proposal, err = tc.QueryProposal(transactor.CliCtx, proposalID, govtypes.StatusVotingPeriod)
	assert.Error(t, err, "fail to submit proposal due to muted depositor")
}
