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

	newEpoch := NewEpoch(latestEpoch.Id.AddRaw(1), now)
	newEpoch.TotalFee = getTotalFee(ctx, keeper)
	newEpoch.ValidatorSnapshotKeys = getValidatorSnapshotKeys(ctx, keeper)
	keeper.SetLatestEpoch(ctx, newEpoch)
}

func getTotalFee(ctx sdk.Context, keeper Keeper) sdk.Int {
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

	return totalFee
}

func getValidatorSnapshotKeys(ctx sdk.Context, keeper Keeper) [][]byte {
	var snapshotKeys [][]byte
	validators := keeper.validatorKeeper.GetValidators(ctx)
	for _, validator := range validators {
		ethAddr := validator.Description.Moniker
		candidate := keeper.validatorKeeper.GetCandidate(ctx, ethAddr)
		snapshotKeys = append(snapshotKeys, candidate.GetSnapshotKey())
	}

	return snapshotKeys
}
