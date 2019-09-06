package global

import (
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey  sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc       *codec.Codec // The wire codec for binary encoding/decoding.
	ethClient *mainchain.EthClient
}

// NewKeeper creates new instances of the global Keeper
func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec, ethClient *mainchain.EthClient) Keeper {
	return Keeper{
		storeKey:  storeKey,
		cdc:       cdc,
		ethClient: ethClient,
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
