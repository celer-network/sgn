package sync

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis - store genesis parameters
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	k.SetChangeID(ctx, data.StartingChangeID)
	k.SetVotingParams(ctx, data.VotingParams)
	k.SetTallyParams(ctx, data.TallyParams)

	for _, change := range data.Changes {
		k.SetChange(ctx, change)
	}
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	startingChangeID, _ := k.GetChangeID(ctx)
	votingParams := k.GetVotingParams(ctx)
	tallyParams := k.GetTallyParams(ctx)
	changes := k.GetChanges(ctx)

	return GenesisState{
		StartingChangeID: startingChangeID,
		Changes:          changes,
		VotingParams:     votingParams,
		TallyParams:      tallyParams,
	}
}
