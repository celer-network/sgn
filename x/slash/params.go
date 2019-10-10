package slash

import (
	"github.com/celer-network/sgn/x/slash/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Default parameter namespace
const (
	DefaultParamspace = types.ModuleName
)

// ParamTable for slash module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&types.Params{})
}

// SignedBlocksWindow - sliding window for downtime slashing
func (k Keeper) SignedBlocksWindow(ctx sdk.Context) (res int64) {
	k.paramstore.Get(ctx, types.KeySignedBlocksWindow, &res)
	return
}

// MinSignedPerWindow - minimum blocks signed per window
func (k Keeper) MinSignedPerWindow(ctx sdk.Context) (res int64) {
	k.paramstore.Get(ctx, types.KeyMinSignedPerWindow, &res)
	return
}

// SlashFractionDoubleSign - number of guards to handle the request
func (k Keeper) SlashFractionDoubleSign(ctx sdk.Context) (res sdk.Dec) {
	k.paramstore.Get(ctx, types.KeySlashFractionDoubleSign, &res)
	return
}

// SlashFractionDowntime - request limit per epoch
func (k Keeper) SlashFractionDowntime(ctx sdk.Context) (res sdk.Dec) {
	k.paramstore.Get(ctx, types.KeySlashFractionDowntime, &res)
	return
}

// Get all parameteras as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.SignedBlocksWindow(ctx),
		k.MinSignedPerWindow(ctx),
		k.SlashFractionDoubleSign(ctx),
		k.SlashFractionDowntime(ctx),
	)
}

// set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}
