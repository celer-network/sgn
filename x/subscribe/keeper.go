package subscribe

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	coinKeeper bank.Keeper

	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context

	cdc *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the subscribe Keeper
func NewKeeper(coinKeeper bank.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		coinKeeper: coinKeeper,
		storeKey:   storeKey,
		cdc:        cdc,
	}
}

// Gets the entire Subscription metadata for a ethAddress
func (k Keeper) GetSubscription(ctx sdk.Context, ethAddress string) Subscription {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(ethAddress))

	var subscription Subscription
	k.cdc.MustUnmarshalBinaryBare(bz, &subscription)
	return subscription
}

// Sets the entire Subscription metadata for a ethAddress
func (k Keeper) Subscribe(ctx sdk.Context, ethAddress string) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(ethAddress), k.cdc.MustMarshalBinaryBare(NewSubscription(2)))
}
