package subscribe

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, keeper Keeper) {
	latestEpoch := keeper.GetLatestEpoch(ctx)
	epochLength := keeper.EpochLength(ctx)
	now := ctx.BlockTime().Unix()

	if now-latestEpoch.Timestamp < epochLength {
		return
	}

	costPerEpoch := keeper.CostPerEpoch(ctx)
	totalFee := sdk.ZeroInt()
	keeper.IterateSubscriptions(ctx, func(subscription Subscription) (stop bool) {
		subscription.Subscribing = false
		if subscription.Deposit.Sub(subscription.Spend).GTE(costPerEpoch) {
			subscription.Spend = subscription.Spend.Add(costPerEpoch)
			totalFee = totalFee.Add(costPerEpoch)
			subscription.Subscribing = true
		}

		keeper.SetSubscription(ctx, subscription)
		return false
	})

	newEpoch := NewEpoch(latestEpoch.Id.AddRaw(1), now)
	newEpoch.TotalFee = totalFee
	keeper.SetLatestEpoch(ctx, newEpoch)
}
