package subscribe

import (
	"github.com/celer-network/sgn/x/global"
	"github.com/celer-network/sgn/x/validator"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, keeper Keeper) {
	latestEpoch := keeper.globalKeeper.GetLatestEpoch(ctx)
	epochLength := keeper.globalKeeper.EpochLength(ctx)
	now := ctx.BlockTime().Unix()

	if now-latestEpoch.Timestamp < epochLength {
		return
	}

	newEpoch := global.NewEpoch(latestEpoch.Id.AddRaw(1), now)
	newEpoch.TotalFee = getTotalFee(ctx, keeper)
	distributeReward(ctx, keeper, newEpoch)
	keeper.globalKeeper.SetLatestEpoch(ctx, newEpoch)
}

// Calculate total fee will be collected in current epoch and reset subscription requestcount
func getTotalFee(ctx sdk.Context, keeper Keeper) sdk.Int {
	costPerEpoch := keeper.globalKeeper.CostPerEpoch(ctx)
	totalFee := sdk.ZeroInt()
	keeper.IterateSubscriptions(ctx, func(subscription Subscription) (stop bool) {
		subscription.RequestCount = 0
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

// Calculate reward distribution to each delegator
func distributeReward(ctx sdk.Context, keeper Keeper, epoch global.Epoch) {
	validators := keeper.validatorKeeper.GetValidators(ctx)
	totalStake := sdk.ZeroInt()
	var candidates []validator.Candidate

	for _, validator := range validators {
		ethAddr := validator.Description.Identity
		candidate := keeper.validatorKeeper.GetCandidate(ctx, ethAddr)
		totalStake = totalStake.Add(candidate.TotalStake)
		candidates = append(candidates, candidate)
	}

	for _, candidate := range candidates {
		for _, delegator := range candidate.Delegators {
			reward, _ := keeper.validatorKeeper.GetReward(ctx, delegator.EthAddress)
			delegatorFee := epoch.TotalFee.Mul(delegator.Stake).Quo(totalStake)
			reward.ServiceReward = reward.ServiceReward.Add(delegatorFee)
			keeper.validatorKeeper.SetReward(ctx, delegator.EthAddress, reward)
		}
	}
}
