package global

import (
	"context"
	"errors"
	"math"

	"github.com/celer-network/sgn/mainchain"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const maxBlockInterval = 2

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

// sync the lastest Block metadata
func (k Keeper) SyncBlock(ctx sdk.Context, blockNumber uint64) error {
	lastestBlock := k.GetLatestBlock(ctx)
	if blockNumber < lastestBlock.Number {
		return errors.New("Block number is smaller than current latest block")
	}

	store := ctx.KVStore(k.storeKey)
	head, err := k.ethClient.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return err
	}

	if math.Abs(float64(blockNumber-head.Number.Uint64())) > maxBlockInterval {
		return errors.New("Block number is out of bound")
	}

	newBlock := NewBlock(blockNumber)
	store.Set(LatestBlockKey, k.cdc.MustMarshalBinaryBare(newBlock))
	return nil
}
