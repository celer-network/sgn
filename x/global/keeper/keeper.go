package keeper

import (
	"github.com/celer-network/sgn/x/global/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper defines the global module Keeper
type Keeper struct {
	// The reference to the Paramstore to get and set global specific params
	paramSpace types.ParamSubspace

	// The (unexposed) keys used to access the stores from the Context.
	storeKey sdk.StoreKey

	// The codec codec for binary encoding/decoding.
	cdc *codec.Codec
}

func NewKeeper(
	cdc *codec.Codec, key sdk.StoreKey, paramSpace types.ParamSubspace) Keeper {
	return Keeper{
		storeKey:   key,
		paramSpace: paramSpace,
		cdc:        cdc,
	}
}

func (keeper Keeper) GetEthBlkNum(ctx sdk.Context) (ethBlkNum uint64) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(types.EthBlkNumKey)
	if bz == nil {
		return
	}
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &ethBlkNum)
	return
}

func (keeper Keeper) SetEthBlkNum(ctx sdk.Context, ethBlkNum uint64) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(ethBlkNum)
	store.Set(types.EthBlkNumKey, bz)
}
