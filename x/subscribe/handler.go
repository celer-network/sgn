package subscribe

import (
	"fmt"
	"math/big"

	"github.com/celer-network/sgn/proto/chain"
	"github.com/celer-network/sgn/proto/entity"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	protobuf "github.com/golang/protobuf/proto"
)

// NewHandler returns a handler for "subscribe" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSubscribe:
			return handleMsgSubscribe(ctx, keeper, msg)
		case MsgRequestGuard:
			return handleMsgRequestGuard(ctx, keeper, msg)
		case MsgGuardProof:
			return handleMsgGuardProof(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized subscribe Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle a message to subscribe
func handleMsgSubscribe(ctx sdk.Context, keeper Keeper, msg MsgSubscribe) sdk.Result {
	deposit, err := keeper.ethClient.Guard.SubscriptionDeposits(&bind.CallOpts{
		BlockNumber: new(big.Int).SetUint64(keeper.globalKeeper.GetSecureBlockNum(ctx)),
	}, ethcommon.HexToAddress(msg.EthAddress))
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("Failed to query subscription desposit: %s", err)).Result()
	}

	subscription, found := keeper.GetSubscription(ctx, msg.EthAddress)
	if !found {
		subscription = NewSubscription(msg.EthAddress)
	}
	subscription.Deposit = sdk.NewIntFromBigInt(deposit)
	keeper.SetSubscription(ctx, subscription)
	return sdk.Result{}
}

// Handle a message to request guard
func handleMsgRequestGuard(ctx sdk.Context, keeper Keeper, msg MsgRequestGuard) sdk.Result {
	logger := ctx.Logger()

	subscription, found := keeper.GetSubscription(ctx, msg.EthAddress)
	if !found {
		return sdk.ErrInternal("Cannot find subscription").Result()
	}

	keeper.SetSubscription(ctx, subscription)

	var signedSimplexState chain.SignedSimplexState
	err := protobuf.Unmarshal(msg.SignedSimplexStateBytes, &signedSimplexState)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("Failed to unmarshal signedSimplexStateBytes: %s", err)).Result()
	}

	var simplexPaymentChannel entity.SimplexPaymentChannel
	err = protobuf.Unmarshal(signedSimplexState.SimplexState, &simplexPaymentChannel)
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
	logger.Info("Set request", "channelId", simplexPaymentChannel.ChannelId, "request seqNum", request.SeqNum)

	return sdk.Result{}
}

// Handle a message to submit guard proof
func handleMsgGuardProof(ctx sdk.Context, keeper Keeper, msg MsgGuardProof) sdk.Result {
	request, found := keeper.GetRequest(ctx, msg.ChannelId)
	if !found {
		return sdk.ErrInternal("Cannot find request").Result()
	}

	// TODO: need to validate tx hash
	request.TxHash = msg.TxHash
	keeper.SetRequest(ctx, msg.ChannelId, request)

	// TODO: punish corresponding guard correctly
	// for _, guard := range request.RequestGuards {
	// 	keeper.slashKeeper.HandleGuardFailure(ctx, msg.Sender, guard)
	// }

	return sdk.Result{}
}
