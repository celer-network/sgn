package keeper

import (
	"fmt"
	"time"

	"github.com/celer-network/sgn/x/sync/types"
	"github.com/celer-network/sgn/x/slash"
	"github.com/celer-network/sgn/x/validator"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/tendermint/libs/log"
)

// Keeper defines the sync module Keeper
type Keeper struct {
	// The reference to the Paramstore to get and set sync specific params
	paramSpace types.ParamSubspace

	vk validator.Keeper

	sk slash.Keeper

	// The (unexposed) keys used to access the stores from the Context.
	storeKey sdk.StoreKey

	// The codec codec for binary encoding/decoding.
	cdc *codec.Codec

	// Change router
	router types.Router
}

// NewKeeper returns a sync keeper. It handles:
// - submitting sync changes
// - depositing funds into changes, and activating upon sufficient funds being deposited
// - users voting on changes, with weight proportional to stake in the system
// - and tallying the result of the vote.
//
// CONTRACT: the parameter Subspace must have the param key table already initialized
func NewKeeper(
	cdc *codec.Codec, key sdk.StoreKey, paramSpace types.ParamSubspace,
	vk validator.Keeper, sk slash.Keeper, rtr types.Router,
) Keeper {

	// It is vital to seal the sync change router here as to not allow
	// further handlers to be registered after the keeper is created since this
	// could create invalid or non-deterministic behavior.
	rtr.Seal()

	return Keeper{
		storeKey:   key,
		paramSpace: paramSpace,
		vk:         vk,
		sk:         sk,
		cdc:        cdc,
		router:     rtr,
	}
}

// Logger returns a module-specific logger.
func (keeper Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// Router returns the sync Keeper's Router
func (keeper Keeper) Router() types.Router {
	return keeper.router
}

// ChangeQueues

// InsertActiveChangeQueue inserts a ChangeID into the active change queue at endTime
func (keeper Keeper) InsertActiveChangeQueue(ctx sdk.Context, changeID uint64, endTime time.Time) {
	store := ctx.KVStore(keeper.storeKey)
	bz := types.GetChangeIDBytes(changeID)
	store.Set(types.ActiveChangeQueueKey(changeID, endTime), bz)
}

// RemoveFromActiveChangeQueue removes a changeID from the Active Change Queue
func (keeper Keeper) RemoveFromActiveChangeQueue(ctx sdk.Context, changeID uint64, endTime time.Time) {
	store := ctx.KVStore(keeper.storeKey)
	store.Delete(types.ActiveChangeQueueKey(changeID, endTime))
}

// InsertInactiveChangeQueue Inserts a ChangeID into the inactive change queue at endTime
func (keeper Keeper) InsertInactiveChangeQueue(ctx sdk.Context, changeID uint64, endTime time.Time) {
	store := ctx.KVStore(keeper.storeKey)
	bz := types.GetChangeIDBytes(changeID)
	store.Set(types.InactiveChangeQueueKey(changeID, endTime), bz)
}

// RemoveFromInactiveChangeQueue removes a changeID from the Inactive Change Queue
func (keeper Keeper) RemoveFromInactiveChangeQueue(ctx sdk.Context, changeID uint64, endTime time.Time) {
	store := ctx.KVStore(keeper.storeKey)
	store.Delete(types.InactiveChangeQueueKey(changeID, endTime))
}

// Iterators

// IterateActiveChangesQueue iterates over the changes in the active change queue
// and performs a callback function
func (keeper Keeper) IterateActiveChangesQueue(ctx sdk.Context, endTime time.Time, cb func(change types.Change) (stop bool)) {
	iterator := keeper.ActiveChangeQueueIterator(ctx, endTime)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		changeID, _ := types.SplitActiveChangeQueueKey(iterator.Key())
		change, found := keeper.GetChange(ctx, changeID)
		if !found {
			panic(fmt.Sprintf("change %d does not exist", changeID))
		}

		if cb(change) {
			break
		}
	}
}

// IterateInactiveChangesQueue iterates over the changes in the inactive change queue
// and performs a callback function
func (keeper Keeper) IterateInactiveChangesQueue(ctx sdk.Context, endTime time.Time, cb func(change types.Change) (stop bool)) {
	iterator := keeper.InactiveChangeQueueIterator(ctx, endTime)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		changeID, _ := types.SplitInactiveChangeQueueKey(iterator.Key())
		change, found := keeper.GetChange(ctx, changeID)
		if !found {
			panic(fmt.Sprintf("change %d does not exist", changeID))
		}

		if cb(change) {
			break
		}
	}
}

// ActiveChangeQueueIterator returns an sdk.Iterator for all the changes in the Active Queue that expire by endTime
func (keeper Keeper) ActiveChangeQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)
	return store.Iterator(types.ActiveChangeQueuePrefix, sdk.PrefixEndBytes(types.ActiveChangeByTimeKey(endTime)))
}

// InactiveChangeQueueIterator returns an sdk.Iterator for all the changes in the Inactive Queue that expire by endTime
func (keeper Keeper) InactiveChangeQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)
	return store.Iterator(types.InactiveChangeQueuePrefix, sdk.PrefixEndBytes(types.InactiveChangeByTimeKey(endTime)))
}
