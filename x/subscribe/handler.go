package subscribe

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "subscribe" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSubscribe:
			return handleMsgSubscribe(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized subscribe Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle a message to subscribe
func handleMsgSubscribe(ctx sdk.Context, keeper Keeper, msg MsgSubscribe) sdk.Result {
	keeper.Subscribe(ctx, msg.EthAddress)
	return sdk.Result{}
}
