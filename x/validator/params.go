package validator

import (
	"time"

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

// SyncerDuration - syncer duration
func (k Keeper) SyncerDuration(ctx sdk.Context) (res uint) {
	k.paramstore.Get(ctx, types.KeySyncerDuration, &res)
	return
}

// EpochLength - epoch length
func (k Keeper) EpochLength(ctx sdk.Context) (res uint) {
	k.paramstore.Get(ctx, types.KeyEpochLength, &res)
	return
}

// MaxValidatorDiff - max validator add
func (k Keeper) MaxValidatorDiff(ctx sdk.Context) (res uint) {
	k.paramstore.Get(ctx, types.KeyMaxValidatorDiff, &res)
	return
}

// WithdrawWindow - withdraw window
func (k Keeper) WithdrawWindow(ctx sdk.Context) (res time.Duration) {
	k.paramstore.Get(ctx, types.KeyWithdrawWindow, &res)
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

// ProportionalRewardFraction - The percentage reward based on validator stakepool proportionaly
func (k Keeper) ProportionalRewardFraction(ctx sdk.Context) (res sdk.Dec) {
	k.paramstore.Get(ctx, types.KeyProportionalRewardFraction, &res)
	return
}

// Get all parameteras as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.SyncerDuration(ctx),
		k.EpochLength(ctx),
		k.MaxValidatorDiff(ctx),
		k.WithdrawWindow(ctx),
		k.MiningReward(ctx),
		k.PullerReward(ctx),
		k.ProportionalRewardFraction(ctx),
	)
}

// set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}
