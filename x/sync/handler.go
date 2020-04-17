package sync

import (
	"fmt"

	"github.com/celer-network/sgn/x/sync/types"
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

		case MsgVote:
			res, err = handleMsgVote(ctx, keeper, msg)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}

		return res, err
	}
}

func handleMsgSubmitChange(ctx sdk.Context, keeper Keeper, msg MsgSubmitChange) (*sdk.Result, error) {
	change, err := keeper.SubmitChange(ctx, msg.Content)
	if err != nil {
		return nil, err
	}

	votingStarted, err := keeper.AddDeposit(ctx, change.ChangeID, msg.Proposer, msg.InitialDeposit)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Proposer.String()),
		),
	)

	submitEvent := sdk.NewEvent(types.EventTypeSubmitChange, sdk.NewAttribute(types.AttributeKeyChangeType, msg.Content.ChangeType()))
	if votingStarted {
		submitEvent = submitEvent.AppendAttributes(
			sdk.NewAttribute(types.AttributeKeyVotingPeriodStart, fmt.Sprintf("%d", change.ChangeID)),
		)
	}
	ctx.EventManager().EmitEvent(submitEvent)

	return &sdk.Result{
		Data:   GetChangeIDBytes(change.ChangeID),
		Events: ctx.EventManager().Events(),
	}, nil
}

func handleMsgVote(ctx sdk.Context, keeper Keeper, msg MsgVote) (*sdk.Result, error) {
	err := keeper.AddVote(ctx, msg.ChangeID, msg.Voter, msg.Option)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Voter.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
