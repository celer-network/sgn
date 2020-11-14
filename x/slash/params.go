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

// PenaltyDelegatorSize - sliding window for downtime slashing
func (k Keeper) PenaltyDelegatorSize(ctx sdk.Context) (res int64) {
	k.paramstore.Get(ctx, types.KeyPenaltyDelegatorSize, &res)
	return
}

// MinSignedPerWindow - minimum blocks signed per window
func (k Keeper) MinSignedPerWindow(ctx sdk.Context) (res sdk.Dec) {
	k.paramstore.Get(ctx, types.KeyMinSignedPerWindow, &res)
	return
}

// SlashFractionDoubleSign - fraction of power slashed in case of double sign
func (k Keeper) SlashFractionDoubleSign(ctx sdk.Context) (res sdk.Dec) {
	k.paramstore.Get(ctx, types.KeySlashFractionDoubleSign, &res)
	return
}

// SlashFractionDowntime - fraction of power slashed for downtime
func (k Keeper) SlashFractionDowntime(ctx sdk.Context) (res sdk.Dec) {
	k.paramstore.Get(ctx, types.KeySlashFractionDowntime, &res)
	return
}

// SlashFractionGuardFailure - fraction of power slashed for guard failure
func (k Keeper) SlashFractionGuardFailure(ctx sdk.Context) (res sdk.Dec) {
	k.paramstore.Get(ctx, types.KeySlashFractionGuardFailure, &res)
	return
}

// FallbackGuardReward - fraction of penalty for reward to the guard
func (k Keeper) FallbackGuardReward(ctx sdk.Context) (res sdk.Dec) {
	k.paramstore.Get(ctx, types.KeyFallbackGuardReward, &res)
	return
}

// SyncerReward - fraction of penalty for reward to the syncer
func (k Keeper) SyncerReward(ctx sdk.Context) (res sdk.Dec) {
	k.paramstore.Get(ctx, types.KeySyncerReward, &res)
	return
}

// Get all parameteras as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.SignedBlocksWindow(ctx),
		k.PenaltyDelegatorSize(ctx),
		k.MinSignedPerWindow(ctx),
		k.SlashFractionDoubleSign(ctx),
		k.SlashFractionDowntime(ctx),
		k.SlashFractionGuardFailure(ctx),
		k.FallbackGuardReward(ctx),
		k.SyncerReward(ctx),
	)
}

// set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}
