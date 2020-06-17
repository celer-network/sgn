package cron

import (
	"github.com/celer-network/goutils/log"
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
		resetRequestCount(ctx, keeper)
	}
}

func resetRateLimit(ctx sdk.Context, keeper Keeper) {
	candidates := keeper.validatorKeeper.GetAllCandidates(ctx)

	for _, candidate := range candidates {
		totalAccounts := int64(len(candidate.Transactors) + 1)
		quota := sdk.NewCoins(sdk.NewCoin(common.QuotaCoinName, candidate.StakingPool.QuoRaw(common.TokenDec).QuoRaw(totalAccounts)))
		// NOTE: make sure sendEnable is false
		err := keeper.bankKeeper.SetCoins(ctx, candidate.Operator, quota)
		if err != nil {
			log.Errorln("SetCoins err", err)
			continue
		}

		for _, transactor := range candidate.Transactors {
			err = keeper.bankKeeper.SetCoins(ctx, transactor, quota)
			if err != nil {
				log.Errorln("SetCoins err", err)
				continue
			}
		}
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
