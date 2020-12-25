package global

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler creates an sdk.Handler for all the global type messages
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (res *sdk.Result, err error) {
		return
		// logEntry := seal.NewMsgLog()
		// ctx = ctx.WithEventManager(sdk.NewEventManager())

		// switch msg := msg.(type) {
		// default:
		// 	return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		// }

		// if err != nil {
		// 	logEntry.Error = append(logEntry.Error, err.Error())
		// }

		// seal.CommitMsgLog(logEntry)
		// return res, err
	}
}
