package guardianmanager

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
	err := keeper.RequestGuard(ctx, msg.EthAddress, msg.SignedSimplexStateBytes)
	if err != nil {
		ctx.Logger().Error(err.Error())
		return err.Result()
	}

	return sdk.Result{}
}
