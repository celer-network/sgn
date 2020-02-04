package subscribe

import (
	"github.com/celer-network/sgn/x/global"
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
		keeper.validatorKeeper.DistributeServiceReward(ctx, latestEpoch.TotalFee)
	}

	newEpoch := global.NewEpoch(latestEpoch.Id.AddRaw(1), now)
	keeper.globalKeeper.SetLatestEpoch(ctx, newEpoch)
}
