package cron

import (
	"time"

	"github.com/celer-network/sgn/x/validator"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey        sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc             *codec.Codec // The wire codec for binary encoding/decoding.
	bankKeeper      bank.Keeper
	validatorKeeper validator.Keeper
}

// NewKeeper creates new instances of the cron Keeper
func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec,
	bankKeeper bank.Keeper, validatorKeeper validator.Keeper) Keeper {
	return Keeper{
		storeKey:        storeKey,
		cdc:             cdc,
		bankKeeper:      bankKeeper,
		validatorKeeper: validatorKeeper,
	}
}

// Get the DailyTimestamp
func (k Keeper) GetDailyTimestamp(ctx sdk.Context) (t time.Time) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(DailyTimestampKey)
	if bz == nil {
		return
	}

	t.UnmarshalBinary(bz)
	return
}

// Set the DailyTimestamp
func (k Keeper) SetDailyTimestamp(ctx sdk.Context, t time.Time) {
	store := ctx.KVStore(k.storeKey)
	timeBytes, _ := t.MarshalBinary()
	store.Set(DailyTimestampKey, timeBytes)
}
