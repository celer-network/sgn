package global

import (
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey   sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc        *codec.Codec // The wire codec for binary encoding/decoding.
	ethClient  *mainchain.EthClient
	paramstore params.Subspace
}

// NewKeeper creates new instances of the global Keeper
func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec, ethClient *mainchain.EthClient, paramstore params.Subspace) Keeper {
	return Keeper{
		storeKey:   storeKey,
		cdc:        cdc,
		ethClient:  ethClient,
		paramstore: paramstore.WithKeyTable(ParamKeyTable()),
	}
}

// Gets the lastest Block metadata
func (k Keeper) GetLatestBlock(ctx sdk.Context) Block {
	store := ctx.KVStore(k.storeKey)
	if !store.Has(LatestBlockKey) {
		return Block{}
	}

	var lastestBlock Block
	bz := store.Get(LatestBlockKey)
	k.cdc.MustUnmarshalBinaryBare(bz, &lastestBlock)
	return lastestBlock
}

// Gets the secure block number
func (k Keeper) GetSecureBlockNum(ctx sdk.Context) uint64 {
	latestBlock := k.GetLatestBlock(ctx)

	if latestBlock.Number < common.ConfirmationCount {
		return 0
	}

	return latestBlock.Number - common.ConfirmationCount
}

// Sync the lastest Block metadata
func (k Keeper) SyncBlock(ctx sdk.Context, blockNumber uint64) {
	store := ctx.KVStore(k.storeKey)
	newBlock := NewBlock(blockNumber)
	store.Set(LatestBlockKey, k.cdc.MustMarshalBinaryBare(newBlock))
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
