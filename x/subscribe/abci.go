package subscribe

import (
	log "github.com/celer-network/sgn/clog"
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

	if latestEpoch.TotalFee.IsPositive() {
		distributeReward(ctx, keeper, latestEpoch)
	}

	newEpoch := global.NewEpoch(latestEpoch.Id.AddRaw(1), now)
	keeper.globalKeeper.SetLatestEpoch(ctx, newEpoch)
}

// Calculate reward distribution to each delegator
func distributeReward(ctx sdk.Context, keeper Keeper, epoch global.Epoch) {
	validators := keeper.validatorKeeper.GetValidators(ctx)
	totalStake := sdk.ZeroInt()
	var candidates []validator.Candidate

	for _, validator := range validators {
		ethAddr := validator.Description.Identity
		if ethAddr == "" {
			log.Errorf("Miss eth address for validator %x", validator.OperatorAddress)
			continue
		}
		candidate, found := keeper.validatorKeeper.GetCandidate(ctx, ethAddr)

		if found {
			totalStake = totalStake.Add(candidate.StakingPool)
			candidates = append(candidates, candidate)
		}
	}

	for _, candidate := range candidates {
		totalReward := epoch.TotalFee.Mul(candidate.StakingPool).Quo(totalStake)
		keeper.validatorKeeper.HandleServiceReward(ctx, candidate, totalReward)
	}
}
