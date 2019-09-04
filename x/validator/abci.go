package validator

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const pullerDuration = 10

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, req abci.RequestEndBlock, keeper Keeper) {
	puller := keeper.GetPuller(ctx)
	validators := keeper.stakingKeeper.GetBondedValidatorsByPower(ctx)
	vIdx := uint(req.Height) / pullerDuration % uint(len(validators))

	if puller.ValidatorIdx != vIdx || puller.ValidatorAddr.Empty() {
		puller = NewPuller(vIdx, sdk.AccAddress(validators[vIdx].OperatorAddress))
		keeper.SetPuller(ctx, puller)
	}
}
