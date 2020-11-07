package validator

import (
	"fmt"
	"math"

	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/seal"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler returns a handler for "validator" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		logEntry := seal.NewMsgLog()
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		var res *sdk.Result
		var err error

		switch msg := msg.(type) {
		case MsgSetTransactors:
			res, err = handleMsgSetTransactors(ctx, keeper, msg, logEntry)
		case MsgEditCandidateDescription:
			res, err = handleMsgEditCandidateDescription(ctx, keeper, msg, logEntry)
		case MsgWithdrawReward:
			res, err = handleMsgWithdrawReward(ctx, keeper, msg, logEntry)
		case MsgSignReward:
			res, err = handleMsgSignReward(ctx, keeper, msg, logEntry)
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

// Handle a message to set transactors
func handleMsgSetTransactors(ctx sdk.Context, keeper Keeper, msg MsgSetTransactors, logEntry *seal.MsgLog) (*sdk.Result, error) {
	logEntry.Type = msg.Type()
	logEntry.Sender = msg.Sender.String()

	validator, found := keeper.GetValidator(ctx, sdk.ValAddress(msg.Sender))
	if !found {
		return nil, fmt.Errorf("Sender is not a validator")
	}

	candidate, found := keeper.GetCandidate(ctx, validator.Description.Identity)
	if !found {
		return nil, fmt.Errorf("Candidate does not exist")
	}

	dedup := make(map[string]bool)
	oldTransactors := candidate.Transactors
	candidate.Transactors = []sdk.AccAddress{}
	for _, transactor := range msg.Transactors {
		if !transactor.Equals(candidate.ValAccount) {
			if _, exist := dedup[transactor.String()]; !exist {
				logEntry.Transactor = append(logEntry.Transactor, transactor.String())
				candidate.Transactors = append(candidate.Transactors, transactor)
				dedup[transactor.String()] = true
				keeper.InitAccount(ctx, transactor)

				err := keeper.bankKeeper.SetCoins(ctx, candidate.ValAccount, sdk.NewCoins(sdk.NewCoin(common.QuotaCoinName, sdk.NewInt(math.MaxInt64))))
				if err != nil {
					return nil, fmt.Errorf("Failed to SetCoins: %s", err)
				}
			}
		}
	}

	for _, transactor := range oldTransactors {
		if _, exist := dedup[transactor.String()]; !exist {
			keeper.RemoveAccount(ctx, transactor)
			err := keeper.bankKeeper.SetCoins(ctx, candidate.ValAccount, sdk.NewCoins(sdk.NewCoin(common.QuotaCoinName, sdk.ZeroInt())))
			if err != nil {
				return nil, fmt.Errorf("Failed to SetCoins: %s", err)
			}
		}
	}

	keeper.SetCandidate(ctx, candidate)
	return &sdk.Result{}, nil
}

// Handle a message to edit candidate description
func handleMsgEditCandidateDescription(ctx sdk.Context, keeper Keeper, msg MsgEditCandidateDescription, logEntry *seal.MsgLog) (*sdk.Result, error) {
	logEntry.Type = msg.Type()
	logEntry.Sender = msg.Sender.String()
	logEntry.EthAddress = msg.EthAddress

	candidate, found := keeper.GetCandidate(ctx, msg.EthAddress)
	if !found {
		return nil, fmt.Errorf("Candidate does not exist")
	}

	description, err := candidate.Description.UpdateDescription(msg.Description)
	if err != nil {
		return nil, err
	}

	candidate.Description = description
	keeper.SetCandidate(ctx, candidate)
	return &sdk.Result{}, nil
}

// Handle a message to withdraw reward
func handleMsgWithdrawReward(ctx sdk.Context, keeper Keeper, msg MsgWithdrawReward, logEntry *seal.MsgLog) (*sdk.Result, error) {
	logEntry.Type = msg.Type()
	logEntry.Sender = msg.Sender.String()
	logEntry.EthAddress = msg.EthAddress

	reward, found := keeper.GetReward(ctx, msg.EthAddress)
	if !found {
		return nil, fmt.Errorf("Reward does not exist")
	}

	if ctx.BlockTime().Before(reward.LastWithdrawTime.Add(keeper.WithdrawWindow(ctx))) {
		logEntry.Warn = append(logEntry.Warn, "Request too fast")
		return &sdk.Result{}, nil
	}

	reward.InitateWithdraw(ctx.BlockTime())
	keeper.SetReward(ctx, reward)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			ModuleName,
			sdk.NewAttribute(sdk.AttributeKeyAction, ActionInitiateWithdraw),
			sdk.NewAttribute(AttributeKeyEthAddress, msg.EthAddress),
		),
	)
	return &sdk.Result{
		Events: ctx.EventManager().Events(),
	}, nil
}

// Handle a message to sign reward
func handleMsgSignReward(ctx sdk.Context, keeper Keeper, msg MsgSignReward, logEntry *seal.MsgLog) (*sdk.Result, error) {
	logEntry.Type = msg.Type()
	logEntry.Sender = msg.Sender.String()
	logEntry.EthAddress = msg.EthAddress

	validator, found := keeper.stakingKeeper.GetValidator(ctx, sdk.ValAddress(msg.Sender))
	if !found {
		return nil, fmt.Errorf("Sender is not validator")
	}
	if validator.Status != sdk.Bonded {
		return nil, fmt.Errorf("Validator is not bonded")
	}

	reward, found := keeper.GetReward(ctx, msg.EthAddress)
	if !found {
		return nil, fmt.Errorf("Reward does not exist")
	}

	err := reward.AddSig(msg.Sig, validator.Description.Identity)
	if err != nil {
		return nil, fmt.Errorf("Failed to add sig: %s", err)
	}

	keeper.SetReward(ctx, reward)
	return &sdk.Result{}, nil
}
