package guard

import (
	"github.com/celer-network/sgn/x/guard/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Default parameter namespace
const (
	DefaultParamspace = types.ModuleName
)

// ParamTable for guard module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&types.Params{})
}

// RequestGuardCount - number of guards to handle the request
func (k Keeper) RequestGuardCount(ctx sdk.Context) (res uint64) {
	k.paramstore.Get(ctx, types.KeyRequestGuardCount, &res)
	return
}

// RequestCost - cost per request
func (k Keeper) RequestCost(ctx sdk.Context) (res sdk.Int) {
	k.paramstore.Get(ctx, types.KeyRequestCost, &res)
	return
}

// EpochLength - epoch length based on seconds
func (k Keeper) EpochLength(ctx sdk.Context) (res uint64) {
	k.paramstore.Get(ctx, types.KeyEpochLength, &res)
	return
}

// MinDisputeTimeout - minimal channel dispute timeout in mainchain blocks
func (k Keeper) MinDisputeTimeout(ctx sdk.Context) (res uint64) {
	k.paramstore.Get(ctx, types.KeyMinDisputeTimeout, &res)
	return
}

// Get all parameteras as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.RequestGuardCount(ctx),
		k.EpochLength(ctx),
		k.RequestCost(ctx),
		k.MinDisputeTimeout(ctx),
	)
}

// set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}
