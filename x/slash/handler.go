package slash

import (
	"fmt"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/seal"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "slash" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		logEntry := seal.NewMsgLog()
		var res sdk.Result
		var err error
		switch msg := msg.(type) {
		case MsgSignPenalty:
			res, err = handleMsgSignPenalty(ctx, keeper, msg, logEntry)
		default:
			errMsg := fmt.Sprintf("Unrecognized slash Msg type: %v", msg.Type())
			log.Error(errMsg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}

		if err != nil {
			logEntry.Error = append(logEntry.Error, err.Error())
			seal.CommitMsgLog(logEntry)
			return sdk.ErrInternal(err.Error()).Result()
		}
		seal.CommitMsgLog(logEntry)
		return res
	}
}

// Handle a message to sign penalty
func handleMsgSignPenalty(ctx sdk.Context, keeper Keeper, msg MsgSignPenalty, logEntry *seal.MsgLog) (sdk.Result, error) {
	logEntry.Type = msg.Type()
	logEntry.Sender = msg.Sender.String()
	logEntry.SlashNonce = msg.Nonce

	res := sdk.Result{}
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

	err := penalty.AddSig(msg.Sig, validator.Description.Identity)
	if err != nil {
		return res, fmt.Errorf(fmt.Sprintf("Failed to add sig: %s", err))
	}

	keeper.SetPenalty(ctx, penalty)
	return res, nil
}
