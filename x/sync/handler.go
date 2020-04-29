package sync

import (
	"github.com/celer-network/sgn/seal"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler creates an sdk.Handler for all the sync type messages
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		logEntry := seal.NewMsgLog()
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		var res *sdk.Result
		var err error

		switch msg := msg.(type) {
		case MsgSubmitChange:
			res, err = handleMsgSubmitChange(ctx, keeper, msg, logEntry)
		case MsgApprove:
			res, err = handleMsgApprove(ctx, keeper, msg, logEntry)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}

		if err != nil {
			logEntry.Error = append(logEntry.Error, err.Error())
		}

		seal.CommitMsgLog(logEntry)
		return res, err
	}
}

func handleMsgSubmitChange(ctx sdk.Context, keeper Keeper, msg MsgSubmitChange, logEntry *seal.MsgLog) (*sdk.Result, error) {
	logEntry.Type = msg.Type()
	logEntry.Sender = msg.Sender.String()
	logEntry.Change.Type = msg.ChangeType

	change, err := keeper.SubmitChange(ctx, msg.ChangeType, msg.Data, msg.Sender)
	if err != nil {
		return nil, err
	}

	logEntry.Change.Id = change.ID

	return &sdk.Result{
		Data:   GetChangeIDBytes(change.ID),
		Events: ctx.EventManager().Events(),
	}, nil
}

func handleMsgApprove(ctx sdk.Context, keeper Keeper, msg MsgApprove, logEntry *seal.MsgLog) (*sdk.Result, error) {
	logEntry.Type = msg.Type()
	logEntry.Sender = msg.Sender.String()
	logEntry.Change.Id = msg.ChangeID

	err := keeper.ApproveChange(ctx, msg.ChangeID, msg.Sender)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{}, nil
}
