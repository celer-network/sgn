package global

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "global" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSyncBlock:
			return handleMsgSyncBlock(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized global Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle a message to sync block
func handleMsgSyncBlock(ctx sdk.Context, keeper Keeper, msg MsgSyncBlock) sdk.Result {
	err := keeper.SyncBlock(ctx, msg.BlockNumber)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("Failed to sync block: %s", err)).Result()
	}

	return sdk.Result{}
}
