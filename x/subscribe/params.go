package subscribe

import (
	"github.com/celer-network/sgn/x/subscribe/types"
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

// RequestGuardCount - number of handlers to handle the request
func (k Keeper) RequestGuardCount(ctx sdk.Context) (res uint64) {
	k.paramstore.Get(ctx, types.KeyRequestGuardCount, &res)
	return
}

// Get all parameteras as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.RequestGuardCount(ctx),
	)
}

// set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}
