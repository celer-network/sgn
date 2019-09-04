package validator

import (
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/x/global"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey      sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc           *codec.Codec // The wire codec for binary encoding/decoding.
	ethClient     *mainchain.EthClient
	stakingKeeper staking.Keeper
	globalKeeper  global.Keeper
}

// NewKeeper creates new instances of the validator Keeper
func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec, ethClient *mainchain.EthClient, globalKeeper global.Keeper, stakingKeeper staking.Keeper) Keeper {
	return Keeper{
		storeKey:      storeKey,
		cdc:           cdc,
		ethClient:     ethClient,
		stakingKeeper: stakingKeeper,
		globalKeeper:  globalKeeper,
	}
}

// Gets the entire Puller metadata
func (k Keeper) GetPuller(ctx sdk.Context) Puller {
	store := ctx.KVStore(k.storeKey)

	if !store.Has(PullerKey) {
		return Puller{}
	}

	value := store.Get(PullerKey)
	var puller Puller
	k.cdc.MustUnmarshalBinaryBare(value, &puller)
	return puller
}

// Sets the entire Puller metadata
func (k Keeper) SetPuller(ctx sdk.Context, puller Puller) {
	store := ctx.KVStore(k.storeKey)
	store.Set(PullerKey, k.cdc.MustMarshalBinaryBare(puller))
}

// Gets the entire Pusher metadata
func (k Keeper) GetPusher(ctx sdk.Context) Pusher {
	store := ctx.KVStore(k.storeKey)

	if !store.Has(PusherKey) {
		return Pusher{}
	}

	value := store.Get(PusherKey)
	var pusher Pusher
	k.cdc.MustUnmarshalBinaryBare(value, &pusher)
	return pusher
}

// Sets the entire Pusher metadata
func (k Keeper) SetPusher(ctx sdk.Context, pusher Pusher) {
	store := ctx.KVStore(k.storeKey)
	store.Set(PusherKey, k.cdc.MustMarshalBinaryBare(pusher))
}
