package slash

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "slash" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSignPenalty:
			return handleMsgSignPenalty(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized slash Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle a message to sign penalty
func handleMsgSignPenalty(ctx sdk.Context, keeper Keeper, msg MsgSignPenalty) sdk.Result {
	validator, found := keeper.validatorKeeper.GetValidator(ctx, sdk.ValAddress(msg.Sender))
	if !found {
		return sdk.ErrInternal("Sender is not validator").Result()
	}

	penalty, found := keeper.GetPenalty(ctx, msg.Nonce)
	if !found {
		return sdk.ErrInternal("Penalty does not exist").Result()
	}

	err := penalty.AddSig(msg.Sig, validator.Description.Identity)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("Failed to add sig: %s", err)).Result()
	}

	keeper.SetPenalty(ctx, penalty)
	return sdk.Result{}
}
