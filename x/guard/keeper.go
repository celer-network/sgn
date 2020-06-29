package guard

import (
	"fmt"

	"github.com/celer-network/sgn/x/validator"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey        sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc             *codec.Codec // The wire codec for binary encoding/decoding.
	validatorKeeper validator.Keeper
	paramstore      params.Subspace
}

// NewKeeper creates new instances of the guard Keeper
func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec,
	validatorKeeper validator.Keeper, paramstore params.Subspace) Keeper {
	return Keeper{
		storeKey:        storeKey,
		cdc:             cdc,
		validatorKeeper: validatorKeeper,
		paramstore:      paramstore.WithKeyTable(ParamKeyTable()),
	}
}

// Gets the entire Subscription metadata for a ethAddress
func (k Keeper) GetSubscription(ctx sdk.Context, ethAddress string) (subscription Subscription, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(GetSubscriptionKey(ethAddress))

	if value == nil {
		return subscription, false
	}

	k.cdc.MustUnmarshalBinaryBare(value, &subscription)
	return subscription, true
}

// Sets the entire Subscription metadata for a ethAddress
func (k Keeper) SetSubscription(ctx sdk.Context, subscription Subscription) {
	store := ctx.KVStore(k.storeKey)
	store.Set(GetSubscriptionKey(subscription.EthAddress), k.cdc.MustMarshalBinaryBare(subscription))
}

// IterateSubscriptions iterates over the stored ValidatorSigningInfo
func (k Keeper) IterateSubscriptions(ctx sdk.Context,
	handler func(subscription Subscription) (stop bool)) {

	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, SubscriptionKeyPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var subscription Subscription
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &subscription)
		if handler(subscription) {
			break
		}
	}
}

// Charge the fee for request
func (k Keeper) ChargeRequestFee(ctx sdk.Context, ethAddr string) error {
	subscription, found := k.GetSubscription(ctx, ethAddr)
	if !found {
		return fmt.Errorf("cannot find subscription")
	}

	requestCost := k.RequestCost(ctx)
	if subscription.Spend.Add(requestCost).GT(subscription.Deposit) {
		return fmt.Errorf("not enough balance, total deposit %s spend %s, fee %s",
			subscription.Deposit, subscription.Spend, requestCost)
	}

	subscription.Spend = subscription.Spend.Add(k.RequestCost(ctx))
	k.SetSubscription(ctx, subscription)

	latestEpoch := k.GetLatestEpoch(ctx)
	latestEpoch.TotalFee = latestEpoch.TotalFee.Add(requestCost)
	k.SetLatestEpoch(ctx, latestEpoch)

	return nil
}

// Gets the entire Request metadata for a channelId
func (k Keeper) GetRequest(ctx sdk.Context, channelId []byte, receiver string) (Request, bool) {
	store := ctx.KVStore(k.storeKey)

	if !store.Has(GetRequestKey(channelId, receiver)) {
		return Request{}, false
	}

	value := store.Get(GetRequestKey(channelId, receiver))
	var request Request
	k.cdc.MustUnmarshalBinaryBare(value, &request)
	return request, true
}

// Sets the entire Request metadata for a channelId
func (k Keeper) SetRequest(ctx sdk.Context, request Request) {
	store := ctx.KVStore(k.storeKey)
	store.Set(GetRequestKey(request.ChannelId, request.GetReceiverAddress()), k.cdc.MustMarshalBinaryBare(request))
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
