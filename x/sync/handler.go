package sync

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler creates an sdk.Handler for all the sync type messages
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		var res *sdk.Result
		var err error

		switch msg := msg.(type) {

		case MsgSubmitChange:
			res, err = handleMsgSubmitChange(ctx, keeper, msg)

		case MsgApprove:
			res, err = handleMsgApprove(ctx, keeper, msg)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}

		return res, err
	}
}

func handleMsgSubmitChange(ctx sdk.Context, keeper Keeper, msg MsgSubmitChange) (*sdk.Result, error) {
	change, err := keeper.SubmitChange(ctx, msg.ChangeType, msg.Data, msg.Sender)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{
		Data:   GetChangeIDBytes(change.ChangeID),
		Events: ctx.EventManager().Events(),
	}, nil
}

func handleMsgApprove(ctx sdk.Context, keeper Keeper, msg MsgApprove) (*sdk.Result, error) {
	err := keeper.ApproveChange(ctx, msg.ChangeID, msg.Sender)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{}, nil
}
