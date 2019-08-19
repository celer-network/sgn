package guardianmanager

import (
	"fmt"

	"github.com/celer-network/sgn/chain"
	"github.com/celer-network/sgn/entity"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/x/subscribe"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/protobuf/proto"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	coinKeeper      bank.Keeper
	subscribeKeeper subscribe.Keeper
	storeKey        sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc             *codec.Codec // The wire codec for binary encoding/decoding.
	ethClient       *mainchain.EthClient
}

// NewKeeper creates new instances of the guardianmanager Keeper
func NewKeeper(coinKeeper bank.Keeper, subscribeKeeper subscribe.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec, ethClient *mainchain.EthClient) Keeper {
	return Keeper{
		coinKeeper:      coinKeeper,
		subscribeKeeper: subscribeKeeper,
		storeKey:        storeKey,
		cdc:             cdc,
		ethClient:       ethClient,
	}
}

// Gets the entire Guardian metadata for a ethAddress
func (k Keeper) GetGuardian(ctx sdk.Context, ethAddress string) Guardian {
	store := ctx.KVStore(k.storeKey)

	if !store.Has([]byte(ethAddress)) {
		return NewGuardian()
	}

	value := store.Get([]byte(ethAddress))
	var guardian Guardian
	k.cdc.MustUnmarshalBinaryBare(value, &guardian)
	return guardian
}

// Sets the entire Guardian metadata for a ethAddress
func (k Keeper) SetGuardian(ctx sdk.Context, ethAddress string, guardian Guardian) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(ethAddress), k.cdc.MustMarshalBinaryBare(guardian))
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

// Sets the entire Subscription metadata for a ethAddress
func (k Keeper) RequestGuard(ctx sdk.Context, ethAddress string, signedSimplexStateBytes []byte) sdk.Error {
	subscription, found := k.subscribeKeeper.GetSubscription(ctx, ethAddress)
	if !found {
		return sdk.ErrInternal("Cannot find subscription")
	}

	latestBlkNum, err := k.ethClient.GetLatestBlkNum()
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("Failed to query latest block number: %s", err))
	}
	// TODO: add a safe margin to ensure consistent validation and that guardians have enough time to submit tx
	if latestBlkNum > subscription.Expiration {
		return sdk.ErrInternal("Subscription expired")
	}

	var signedSimplexState chain.SignedSimplexState
	err = proto.Unmarshal(signedSimplexStateBytes, &signedSimplexState)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("Failed to unmarshal signedSimplexStateBytes: %s", err))
	}

	var simplexPaymentChannel entity.SimplexPaymentChannel
	err = proto.Unmarshal(signedSimplexState.SimplexState, &simplexPaymentChannel)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("Failed to unmarshal simplexState: %s", err))
	}

	// TODO: add extra validation for the msg
	if simplexPaymentChannel.SeqNum < subscription.SeqNum {
		return sdk.ErrInternal("Seq Num must be larger than previous request")
	}

	subscription.SeqNum = simplexPaymentChannel.SeqNum
	subscription.ChannelId = simplexPaymentChannel.ChannelId
	subscription.SignedSimplexStateBytes = signedSimplexStateBytes
	k.subscribeKeeper.SetSubscription(ctx, ethAddress, subscription)
	return nil
}
