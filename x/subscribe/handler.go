package subscribe

import (
	"fmt"
	"math/big"

	"github.com/celer-network/sgn/chain"
	"github.com/celer-network/sgn/entity"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/protobuf/proto"
)

// NewHandler returns a handler for "subscribe" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSubscribe:
			return handleMsgSubscribe(ctx, keeper, msg)
		case MsgRequestGuard:
			return handleMsgRequestGuard(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized subscribe Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle a message to subscribe
func handleMsgSubscribe(ctx sdk.Context, keeper Keeper, msg MsgSubscribe) sdk.Result {
	expiration, err := keeper.ethClient.Guard.SubscriptionExpiration(&bind.CallOpts{
		BlockNumber: new(big.Int).SetUint64(keeper.globalKeeper.GetSecureBlockNum(ctx)),
	}, ethcommon.HexToAddress(msg.EthAddress))
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("Failed to query subscription expiration: %s", err)).Result()
	}

	keeper.Subscribe(ctx, msg.EthAddress, expiration.Uint64())

	return sdk.Result{}
}

// Handle a message to request guard
func handleMsgRequestGuard(ctx sdk.Context, keeper Keeper, msg MsgRequestGuard) sdk.Result {
	subscription, found := keeper.GetSubscription(ctx, msg.EthAddress)
	if !found {
		return sdk.ErrInternal("Cannot find subscription").Result()
	}

	latestBlk := keeper.globalKeeper.GetLatestBlock(ctx)
	// TODO: add a safe margin to ensure consistent validation and that guardians have enough time to submit tx
	if latestBlk.Number > subscription.Expiration {
		return sdk.ErrInternal("Subscription expired").Result()
	}

	var signedSimplexState chain.SignedSimplexState
	err := proto.Unmarshal(msg.SignedSimplexStateBytes, &signedSimplexState)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("Failed to unmarshal signedSimplexStateBytes: %s", err)).Result()
	}

	var simplexPaymentChannel entity.SimplexPaymentChannel
	err = proto.Unmarshal(signedSimplexState.SimplexState, &simplexPaymentChannel)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("Failed to unmarshal simplexState: %s", err)).Result()
	}

	request, err := getRequest(ctx, keeper, simplexPaymentChannel)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("Failed to get request: %s", err)).Result()
	}

	if simplexPaymentChannel.SeqNum < request.SeqNum {
		return sdk.ErrInternal("Seq Num must be larger than previous request").Result()
	}

	err = verifySignedSimplexStateSigs(request, signedSimplexState)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("Failed to verify sigs: %s", err)).Result()
	}

	request.SeqNum = simplexPaymentChannel.SeqNum
	request.SignedSimplexStateBytes = msg.SignedSimplexStateBytes
	keeper.SetRequest(ctx, simplexPaymentChannel.ChannelId, request)

	return sdk.Result{}
}
