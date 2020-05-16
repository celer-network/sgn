package global

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey   sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc        *codec.Codec // The wire codec for binary encoding/decoding.
	paramstore params.Subspace
}

// NewKeeper creates new instances of the global Keeper
func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec, paramstore params.Subspace) Keeper {
	return Keeper{
		storeKey:   storeKey,
		cdc:        cdc,
		paramstore: paramstore.WithKeyTable(ParamKeyTable()),
	}
}

// Gets the entire Epoch metadata for a epochId
func (k Keeper) GetEpoch(ctx sdk.Context, epochId sdk.Int) (epoch Epoch, found bool) {
	store := ctx.KVStore(k.storeKey)

	if !store.Has(GetEpochKey(epochId)) {
		return epoch, false
	}

	value := store.Get(GetEpochKey(epochId))
	k.cdc.MustUnmarshalBinaryBare(value, &epoch)
	return epoch, true
}

// Sets the entire Epoch metadata for a epochId
func (k Keeper) SetEpoch(ctx sdk.Context, epoch Epoch) {
	store := ctx.KVStore(k.storeKey)
	store.Set(GetEpochKey(epoch.Id), k.cdc.MustMarshalBinaryBare(epoch))
}

// Gets the entire latest Epoch metadata
func (k Keeper) GetLatestEpoch(ctx sdk.Context) (epoch Epoch) {
	store := ctx.KVStore(k.storeKey)

	if !store.Has(GetLatestEpochKey()) {
		epoch = NewEpoch(sdk.NewInt(1), ctx.BlockTime().Unix())
		k.SetLatestEpoch(ctx, epoch)
		return
	}

	value := store.Get(GetLatestEpochKey())
	k.cdc.MustUnmarshalBinaryBare(value, &epoch)
	return
}

// Sets the entire LatestEpoch metadata
func (k Keeper) SetLatestEpoch(ctx sdk.Context, epoch Epoch) {
	store := ctx.KVStore(k.storeKey)
	store.Set(GetLatestEpochKey(), k.cdc.MustMarshalBinaryBare(epoch))
	k.SetEpoch(ctx, epoch)
}
