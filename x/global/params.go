package global

import (
	"github.com/celer-network/sgn/x/global/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Default parameter namespace
const (
	DefaultParamspace = types.ModuleName
)

// ParamTable for global module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&types.Params{})
}

// EpochLength - Epoch length based on seconds
func (k Keeper) EpochLength(ctx sdk.Context) (res int64) {
	k.paramstore.Get(ctx, types.KeyEpochLength, &res)
	return
}

// MaxBlockNumDiff - Max block diff accepted when sync block
func (k Keeper) MaxBlockNumDiff(ctx sdk.Context) (res int64) {
	k.paramstore.Get(ctx, types.KeyMaxBlockNumDiff, &res)
	return
}

// MaxBlockTimeDiff - Max block diff accepted when sync block
func (k Keeper) MaxBlockTimeDiff(ctx sdk.Context) (res int64) {
	k.paramstore.Get(ctx, types.KeyMaxBlockTimeDiff, &res)
	return
}

// ConfirmationCount - Number of blocks to confirm a block is safe
func (k Keeper) ConfirmationCount(ctx sdk.Context) (res uint64) {
	k.paramstore.Get(ctx, types.KeyConfirmationCount, &res)
	return
}

// Get all parameteras as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.EpochLength(ctx),
		k.MaxBlockNumDiff(ctx),
		k.MaxBlockTimeDiff(ctx),
		k.ConfirmationCount(ctx),
	)
}

// set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}
