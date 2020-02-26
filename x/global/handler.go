package global

import (
	"context"
	"fmt"
	"math/big"

	"github.com/celer-network/sgn/seal"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler returns a handler for "global" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		logEntry := seal.NewMsgLog()
		var res *sdk.Result
		var err error
		switch msg := msg.(type) {
		case MsgSyncBlock:
			res, err = handleMsgSyncBlock(ctx, keeper, msg, logEntry)
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

// Handle a message to sync block
func handleMsgSyncBlock(ctx sdk.Context, keeper Keeper, msg MsgSyncBlock, logEntry *seal.MsgLog) (*sdk.Result, error) {
	logEntry.Type = msg.Type()
	logEntry.Sender = msg.Sender.String()
	logEntry.BlockNum = msg.BlockNumber

	latestBlock := keeper.GetLatestBlock(ctx)
	if msg.BlockNumber < latestBlock.Number {
		return nil, fmt.Errorf("smaller than latest blknum %d", latestBlock.Number)
	}

	block, err := keeper.ethClient.Client.HeaderByNumber(context.Background(), big.NewInt(int64(msg.BlockNumber)))
	if err != nil {
		return nil, fmt.Errorf("Failed to query mainchain block by number: %s", err)
	}

	timeDiff := ctx.BlockTime().Unix() - int64(block.Time)
	if timeDiff < keeper.BlkTimeDiffLower(ctx) || timeDiff > keeper.BlkTimeDiffUpper(ctx) {
		return nil, fmt.Errorf("out of bound diff with mainchain timestamp %d and sidechain timestamp %d", block.Time, ctx.BlockTime().Unix())
	}

	keeper.SyncBlock(ctx, msg.BlockNumber)
	return &sdk.Result{}, nil
}
