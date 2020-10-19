package slash

import (
	"fmt"

	"github.com/celer-network/sgn/seal"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler returns a handler for "slash" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		logEntry := seal.NewMsgLog()
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		var res *sdk.Result
		var err error
		switch msg := msg.(type) {
		case MsgSignPenalty:
			res, err = handleMsgSignPenalty(ctx, keeper, msg, logEntry)
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

// Handle a message to sign penalty
func handleMsgSignPenalty(ctx sdk.Context, keeper Keeper, msg MsgSignPenalty, logEntry *seal.MsgLog) (*sdk.Result, error) {
	logEntry.Type = msg.Type()
	logEntry.Sender = msg.Sender.String()
	logEntry.Penalty.Nonce = msg.Nonce

	res := &sdk.Result{}
	validator, found := keeper.validatorKeeper.GetValidator(ctx, sdk.ValAddress(msg.Sender))
	if !found {
		return res, fmt.Errorf("Sender is not validator")
	}
	if validator.Status != sdk.Bonded {
		return res, fmt.Errorf("Validator is not bonded")
	}

	penalty, found := keeper.GetPenalty(ctx, msg.Nonce)
	if !found {
		return res, fmt.Errorf("Penalty does not exist")
	}
	logEntry.Penalty.Validator = penalty.ValidatorAddr
	logEntry.Penalty.Reason = penalty.Reason

	err := penalty.AddSig(msg.Sig, validator.Description.Identity)
	if err != nil {
		return res, fmt.Errorf(fmt.Sprintf("Failed to add sig: %s", err))
	}

	keeper.SetPenalty(ctx, penalty)
	return res, nil
}
