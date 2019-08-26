package guardianmanager

import (
	"fmt"

	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/x/global"
	"github.com/celer-network/sgn/x/subscribe"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	globalKeeper    global.Keeper
	subscribeKeeper subscribe.Keeper
	storeKey        sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc             *codec.Codec // The wire codec for binary encoding/decoding.
	ethClient       *mainchain.EthClient
}

// NewKeeper creates new instances of the guardianmanager Keeper
func NewKeeper(globalKeeper global.Keeper, subscribeKeeper subscribe.Keeper,
	storeKey sdk.StoreKey, cdc *codec.Codec, ethClient *mainchain.EthClient) Keeper {
	return Keeper{
		globalKeeper:    globalKeeper,
		subscribeKeeper: subscribeKeeper,
		storeKey:        storeKey,
		cdc:             cdc,
		ethClient:       ethClient,
	}
}

// Gets the entire Guardian metadata for a ethAddress
func (k Keeper) GetGuardian(ctx sdk.Context, ethAddress string) Guardian {
	store := ctx.KVStore(k.storeKey)

	if !store.Has(GetGuardianKey(ethAddress)) {
		return NewGuardian()
	}

	value := store.Get(GetGuardianKey(ethAddress))
	var guardian Guardian
	k.cdc.MustUnmarshalBinaryBare(value, &guardian)
	return guardian
}

// Sets the entire Guardian metadata for a ethAddress
func (k Keeper) SetGuardian(ctx sdk.Context, ethAddress string, guardian Guardian) {
	store := ctx.KVStore(k.storeKey)
	store.Set(GetGuardianKey(ethAddress), k.cdc.MustMarshalBinaryBare(guardian))
}

// Gets the entire Request metadata for a channelId
func (k Keeper) GetRequest(ctx sdk.Context, channelId []byte) (Request, bool) {
	store := ctx.KVStore(k.storeKey)

	if !store.Has(GetRequestKey(channelId)) {
		return Request{}, false
	}

	value := store.Get(GetRequestKey(channelId))
	var request Request
	k.cdc.MustUnmarshalBinaryBare(value, &request)
	return request, true
}

// Sets the entire Request metadata for a channelId
func (k Keeper) SetRequest(ctx sdk.Context, channelId []byte, request Request) {
	store := ctx.KVStore(k.storeKey)
	store.Set(GetRequestKey(channelId), k.cdc.MustMarshalBinaryBare(request))
}

// Sets the entire Guardian metadata for a ethAddress
func (k Keeper) Deposit(ctx sdk.Context, ethAddress string) sdk.Error {
	deposit, err := k.ethClient.Guard.SecurityDeposit(&bind.CallOpts{}, ethcommon.HexToAddress(ethAddress))
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("Failed to query security deposit: %s", err))
	}

	guardian := k.GetGuardian(ctx, ethAddress)
	guardian.Balance = deposit.Uint64()
	k.SetGuardian(ctx, ethAddress, guardian)
	return nil
}
