package cron

import (
	"github.com/celer-network/sgn/common"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// EndBlocker called every block, process inflation, update cron set.
func EndBlocker(ctx sdk.Context, req abci.RequestEndBlock, keeper Keeper) {
	dailyTimestamp := keeper.GetDailyTimestamp(ctx)

	if ctx.BlockTime().Sub(dailyTimestamp).Hours() > 24 {
		keeper.SetDailyTimestamp(ctx, ctx.BlockTime())
		resetRateLimit(ctx, keeper)
	}
}

func resetRateLimit(ctx sdk.Context, keeper Keeper) {
	candidates := keeper.validatorKeeper.GetAllCandidates(ctx)

	for _, candidate := range candidates {
		// NOTE: make sure sendEnable is false
		keeper.bankKeeper.SetCoins(ctx, candidate.Operator,
			sdk.NewCoins(sdk.NewCoin(common.StakeName, candidate.StakingPool.QuoRaw(common.StakeDec))))
	}
}
