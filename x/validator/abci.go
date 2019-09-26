package validator

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const pullerDuration = 10
const pusherDuration = 10

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, req abci.RequestEndBlock, keeper Keeper) {
	setPuller(ctx, req, keeper)
	setPusher(ctx, keeper)
}

// Update puller for every pullerDuration
func setPuller(ctx sdk.Context, req abci.RequestEndBlock, keeper Keeper) {
	puller := keeper.GetPuller(ctx)
	validators := keeper.GetValidators(ctx)
	vIdx := uint(req.Height) / pullerDuration % uint(len(validators))

	if puller.ValidatorIdx != vIdx || puller.ValidatorAddr.Empty() {
		puller = NewPuller(vIdx, sdk.AccAddress(validators[vIdx].OperatorAddress))
		keeper.SetPuller(ctx, puller)
	}
}

// Update pusher for every pusherDuration
func setPusher(ctx sdk.Context, keeper Keeper) {
	pusher := keeper.GetPusher(ctx)
	validators := keeper.GetValidators(ctx)
	latestBlock := keeper.globalKeeper.GetLatestBlock(ctx)
	vIdx := uint(latestBlock.Number) / pusherDuration % uint(len(validators))

	if pusher.ValidatorIdx != vIdx || pusher.ValidatorAddr.Empty() {
		pusher = NewPusher(vIdx, sdk.AccAddress(validators[vIdx].OperatorAddress))
		keeper.SetPusher(ctx, pusher)
	}
}
