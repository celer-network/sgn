package subscribe

import (
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/x/global"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey     sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc          *codec.Codec // The wire codec for binary encoding/decoding.
	ethClient    *mainchain.EthClient
	globalKeeper global.Keeper
}

// NewKeeper creates new instances of the subscribe Keeper
func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec, ethClient *mainchain.EthClient, globalKeeper global.Keeper) Keeper {
	return Keeper{
		storeKey:     storeKey,
		cdc:          cdc,
		ethClient:    ethClient,
		globalKeeper: globalKeeper,
	}
}

// Gets the entire Subscription metadata for a ethAddress
func (k Keeper) GetSubscription(ctx sdk.Context, ethAddress string) (subscription Subscription, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get([]byte(ethAddress))

	if value == nil {
		return subscription, false
	}

	k.cdc.MustUnmarshalBinaryBare(value, &subscription)
	return subscription, true
}

// Sets the entire Subscription metadata for a ethAddress
func (k Keeper) SetSubscription(ctx sdk.Context, ethAddress string, subscription Subscription) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(ethAddress), k.cdc.MustMarshalBinaryBare(subscription))
}

// Sets the entire Subscription metadata for a ethAddress
func (k Keeper) Subscribe(ctx sdk.Context, ethAddress string, expiration uint64) {
	k.SetSubscription(ctx, ethAddress, NewSubscription(expiration))
}
