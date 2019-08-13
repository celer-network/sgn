package subscribe

import (
	"fmt"

	"github.com/celer-network/sgn/x/guardianmanager"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "subscribe" type messages.
func NewHandler(keeper Keeper, gmKeep guardianmanager.Keeper) sdk.Handler {
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
	err := keeper.Subscribe(ctx, msg.EthAddress)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

// Handle a message to request guard
func handleMsgRequestGuard(ctx sdk.Context, keeper Keeper, msg MsgRequestGuard) sdk.Result {
	// TODO: need to validate signed simple state bytes
	err := keeper.RequestGuard(ctx, msg.EthAddress, msg.SignedSimplexStateBytes)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{}
}
