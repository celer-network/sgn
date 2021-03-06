package gov

import (
	"fmt"

	"github.com/celer-network/sgn/seal"
	"github.com/celer-network/sgn/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler creates an sdk.Handler for all the gov type messages
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		logEntry := seal.NewMsgLog()
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		var res *sdk.Result
		var err error

		switch msg := msg.(type) {
		case MsgDeposit:
			res, err = handleMsgDeposit(ctx, keeper, msg, logEntry)

		case MsgSubmitProposal:
			res, err = handleMsgSubmitProposal(ctx, keeper, msg, logEntry)

		case MsgVote:
			res, err = handleMsgVote(ctx, keeper, msg, logEntry)

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

func handleMsgSubmitProposal(ctx sdk.Context, keeper Keeper, msg MsgSubmitProposal, logEntry *seal.MsgLog) (*sdk.Result, error) {
	logEntry.Type = msg.Type()
	logEntry.Sender = msg.Proposer.String()

	proposal, err := keeper.SubmitProposal(ctx, msg.Content)
	if err != nil {
		return nil, err
	}

	logEntry.Govern.ProposalId = proposal.ProposalID
	votingStarted, err := keeper.AddDeposit(ctx, proposal.ProposalID, msg.Proposer, msg.InitialDeposit)
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

	submitEvent := sdk.NewEvent(types.EventTypeSubmitProposal, sdk.NewAttribute(types.AttributeKeyProposalType, msg.Content.ProposalType()))
	if votingStarted {
		submitEvent = submitEvent.AppendAttributes(
			sdk.NewAttribute(types.AttributeKeyVotingPeriodStart, fmt.Sprintf("%d", proposal.ProposalID)),
		)
	}
	ctx.EventManager().EmitEvent(submitEvent)

	return &sdk.Result{
		Data:   GetProposalIDBytes(proposal.ProposalID),
		Events: ctx.EventManager().Events(),
	}, nil
}

func handleMsgDeposit(ctx sdk.Context, keeper Keeper, msg MsgDeposit, logEntry *seal.MsgLog) (*sdk.Result, error) {
	logEntry.Type = msg.Type()
	logEntry.Sender = msg.Depositor.String()
	logEntry.Govern.ProposalId = msg.ProposalID
	logEntry.Govern.Amount = msg.Amount.Uint64()

	votingStarted, err := keeper.AddDeposit(ctx, msg.ProposalID, msg.Depositor, msg.Amount)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Depositor.String()),
		),
	)

	if votingStarted {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeProposalDeposit,
				sdk.NewAttribute(types.AttributeKeyVotingPeriodStart, fmt.Sprintf("%d", msg.ProposalID)),
			),
		)
	}

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgVote(ctx sdk.Context, keeper Keeper, msg MsgVote, logEntry *seal.MsgLog) (*sdk.Result, error) {
	logEntry.Type = msg.Type()
	logEntry.Sender = msg.Voter.String()
	logEntry.Govern.ProposalId = msg.ProposalID
	logEntry.Govern.Option = uint32(msg.Option)

	err := keeper.AddVote(ctx, msg.ProposalID, msg.Voter, msg.Option)
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
