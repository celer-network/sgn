package guardianmanager

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

// Gets the entire Guardian metadata for a ethAddress
func (k Keeper) GetGuardian(ctx sdk.Context, ethAddress string) Guardian {
	store := ctx.KVStore(k.storeKey)

	if !store.Has([]byte(ethAddress)) {
		return NewGuardian()
	}

	bz := store.Get([]byte(ethAddress))
	var guardian Guardian
	k.cdc.MustUnmarshalBinaryBare(bz, &guardian)
	return guardian
}

// Sets the entire Guardian metadata for a ethAddress
func (k Keeper) SetGuardian(ctx sdk.Context, ethAddress string, guardian Guardian) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(ethAddress), k.cdc.MustMarshalBinaryBare(guardian))
}

// Sets the entire Guardian metadata for a ethAddress
func (k Keeper) Deposit(ctx sdk.Context, ethAddress string) {
	guardian := k.GetGuardian(ctx, ethAddress)
	guardian.Balance = 2

	k.SetGuardian(ctx, ethAddress, guardian)
}
