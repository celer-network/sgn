package guardianmanager

import (
	"fmt"

	"github.com/celer-network/sgn/chain"
	"github.com/celer-network/sgn/entity"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/protobuf/proto"
)

// NewHandler returns a handler for "guardianmanager" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgDeposit:
			return handleMsgDeposit(ctx, keeper, msg)
		case MsgRequestGuard:
			return handleMsgRequestGuard(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized guardianmanager Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle a message to deposit
func handleMsgDeposit(ctx sdk.Context, keeper Keeper, msg MsgDeposit) sdk.Result {
	err := keeper.Deposit(ctx, msg.EthAddress)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

// Handle a message to request guard
func handleMsgRequestGuard(ctx sdk.Context, keeper Keeper, msg MsgRequestGuard) sdk.Result {
	// TODO: need to validate signed simple state bytes
	subscription, found := keeper.subscribeKeeper.GetSubscription(ctx, msg.EthAddress)
	if !found {
		return sdk.ErrInternal("Cannot find subscription").Result()
	}

	latestBlkNum, err := keeper.ethClient.GetLatestBlkNum()
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("Failed to query latest block number: %s", err)).Result()
	}
	// TODO: add a safe margin to ensure consistent validation and that guardians have enough time to submit tx
	if latestBlkNum > subscription.Expiration {
		return sdk.ErrInternal("Subscription expired").Result()
	}

	var signedSimplexState chain.SignedSimplexState
	err = proto.Unmarshal(msg.SignedSimplexStateBytes, &signedSimplexState)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("Failed to unmarshal signedSimplexStateBytes: %s", err)).Result()
	}

	var simplexPaymentChannel entity.SimplexPaymentChannel
	err = proto.Unmarshal(signedSimplexState.SimplexState, &simplexPaymentChannel)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("Failed to unmarshal simplexState: %s", err)).Result()
	}

	// TODO: add extra validation for the msg
	if simplexPaymentChannel.SeqNum < subscription.SeqNum {
		return sdk.ErrInternal("Seq Num must be larger than previous request").Result()
	}

	subscription.SeqNum = simplexPaymentChannel.SeqNum
	subscription.ChannelId = simplexPaymentChannel.ChannelId
	subscription.SignedSimplexStateBytes = msg.SignedSimplexStateBytes
	keeper.subscribeKeeper.SetSubscription(ctx, msg.EthAddress, subscription)

	return sdk.Result{}
}
