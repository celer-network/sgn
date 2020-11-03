package cron

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// EndBlocker called every block, process inflation, update cron set.
func EndBlocker(ctx sdk.Context, req abci.RequestEndBlock, keeper Keeper) {
	dailyTimestamp := keeper.GetDailyTimestamp(ctx)

	if ctx.BlockTime().Sub(dailyTimestamp).Hours() > 24 {
		keeper.SetDailyTimestamp(ctx, ctx.BlockTime())
		resetRequestCount(ctx, keeper)
	}
}

func resetRequestCount(ctx sdk.Context, keeper Keeper) {
	candidates := keeper.validatorKeeper.GetAllCandidates(ctx)

	for _, candidate := range candidates {
		if candidate.RequestCount.IsPositive() {
			candidate.RequestCount = sdk.ZeroInt()
			keeper.validatorKeeper.SetCandidate(ctx, candidate)
		}
	}
}
