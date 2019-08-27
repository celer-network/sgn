package bridge

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context

	cdc *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the bridge Keeper
func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		storeKey: storeKey,
		cdc:      cdc,
	}
}

// Gets the entire EthAddress metadata for a cosmos address
func (k Keeper) GetEthAddress(ctx sdk.Context, address sdk.AccAddress) EthAddress {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(address.Bytes())

	var ethAddress EthAddress
	k.cdc.MustUnmarshalBinaryBare(bz, &ethAddress)
	return ethAddress
}

// Sets the entire EthAddress metadata for a cosmos address
func (k Keeper) SetEthAddress(ctx sdk.Context, sender sdk.AccAddress, ethAddress string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(sender.Bytes(), k.cdc.MustMarshalBinaryBare(NewEthAddress(ethAddress)))
}
