package validator

import (
	"github.com/celer-network/sgn/x/validator/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Default parameter namespace
const (
	DefaultParamspace = types.ModuleName
)

// ParamTable for validator module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&types.Params{})
}

// PullerDuration - puller duration
func (k Keeper) PullerDuration(ctx sdk.Context) (res uint) {
	k.paramstore.Get(ctx, types.KeyPullerDuration, &res)
	return
}

// MiningReward - mining reward
func (k Keeper) MiningReward(ctx sdk.Context) (res sdk.Int) {
	k.paramstore.Get(ctx, types.KeyMiningReward, &res)
	return
}

// PullerReward - puller reward
func (k Keeper) PullerReward(ctx sdk.Context) (res sdk.Int) {
	k.paramstore.Get(ctx, types.KeyPullerReward, &res)
	return
}

// Get all parameteras as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.PullerDuration(ctx),
		k.MiningReward(ctx),
		k.PullerReward(ctx),
	)
}

// set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}
