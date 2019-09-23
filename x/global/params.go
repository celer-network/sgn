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

// ParamTable for subscribe module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&types.Params{})
}

// EpochLength - Epoch length based on seconds
func (k Keeper) EpochLength(ctx sdk.Context) (res int64) {
	k.paramstore.Get(ctx, types.KeyEpochLength, &res)
	return
}

// CostPerEpoch - Cost per epoch
func (k Keeper) CostPerEpoch(ctx sdk.Context) (res sdk.Int) {
	k.paramstore.Get(ctx, types.KeyCostPerEpoch, &res)
	return
}

// Get all parameteras as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.EpochLength(ctx),
		k.CostPerEpoch(ctx),
	)
}

// set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}
