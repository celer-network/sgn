package global

import (
	"context"
	"time"

	"github.com/celer-network/sgn/mainchain"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const BlockTimeout = 15

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
func (k Keeper) GetLatestBlock(ctx sdk.Context) (Block, error) {
	store := ctx.KVStore(k.storeKey)
	var lastestBlock Block

	if store.Has(LatestBlockKey) {
		bz := store.Get(LatestBlockKey)
		k.cdc.MustUnmarshalBinaryBare(bz, &lastestBlock)
	}

	// block is time out try to reload block
	if uint64(time.Now().Unix())-lastestBlock.Time > BlockTimeout {
		head, err := k.ethClient.Client.HeaderByNumber(context.Background(), nil)
		if err != nil {
			return Block{}, err
		}

		lastestBlock = NewBlock(head.Number.Uint64(), head.Time)
		store.Set(LatestBlockKey, k.cdc.MustMarshalBinaryBare(lastestBlock))
	}

	return lastestBlock, nil
}
