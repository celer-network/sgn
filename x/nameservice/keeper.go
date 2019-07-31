package nameservice

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/bank"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	coinKeeper bank.Keeper

	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context

	cdc *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the nameservice Keeper
func NewKeeper(coinKeeper bank.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		coinKeeper: coinKeeper,
		storeKey:   storeKey,
		cdc:        cdc,
	}
}

// Gets the entire Whois metadata struct for a name
func (k Keeper) GetNumber(ctx sdk.Context) Number {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get([]byte("test"))
	var number Number
	k.cdc.MustUnmarshalBinaryBare(bz, &number)
	return number
}

// Sets the entire Whois metadata struct for a name
func (k Keeper) SetNumber(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)

	store.Set([]byte("test"), k.cdc.MustMarshalBinaryBare(Number{Value: 2}))
}
