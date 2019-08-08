package subscribe

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "subscribe" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case NewMsgSubscribe:
			return handleNewMsgSubscribe(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized subscribe Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle a message to set name
func handleNewMsgSubscribe(ctx sdk.Context, keeper Keeper, msg NewMsgSubscribe) sdk.Result {
	keeper.Subscribe(ctx, msg.EthAddress) // If so, set the name to the value specified in the msg.
	return sdk.Result{}                   // return
}
