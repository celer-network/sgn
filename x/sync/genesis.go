package sync

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis - store genesis parameters
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {

	k.SetChangeID(ctx, data.StartingChangeID)
	k.SetDepositParams(ctx, data.DepositParams)
	k.SetVotingParams(ctx, data.VotingParams)
	k.SetTallyParams(ctx, data.TallyParams)

	for _, vote := range data.Votes {
		k.SetVote(ctx, vote)
	}
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	startingChangeID, _ := k.GetChangeID(ctx)
	depositParams := k.GetDepositParams(ctx)
	votingParams := k.GetVotingParams(ctx)
	tallyParams := k.GetTallyParams(ctx)
	changes := k.GetChanges(ctx)

	var changesDeposits Deposits
	var changesVotes Votes
	for _, change := range changes {
		deposits := k.GetDeposits(ctx, change.ChangeID)
		changesDeposits = append(changesDeposits, deposits...)

		votes := k.GetVotes(ctx, change.ChangeID)
		changesVotes = append(changesVotes, votes...)
	}

	return GenesisState{
		StartingChangeID: startingChangeID,
		Deposits:           changesDeposits,
		Votes:              changesVotes,
		Changes:          changes,
		DepositParams:      depositParams,
		VotingParams:       votingParams,
		TallyParams:        tallyParams,
	}
}
