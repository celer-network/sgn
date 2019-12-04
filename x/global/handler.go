package global

import (
	"context"
	"fmt"
	"math/big"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/seal"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "global" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		logEntry := seal.NewMsgLog()
		var res sdk.Result
		var err error
		switch msg := msg.(type) {
		case MsgSyncBlock:
			res, err = handleMsgSyncBlock(ctx, keeper, msg, logEntry)
		default:
			errMsg := fmt.Sprintf("Unrecognized global Msg type: %v", msg.Type())
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

// Handle a message to sync block
func handleMsgSyncBlock(ctx sdk.Context, keeper Keeper, msg MsgSyncBlock, logEntry *seal.MsgLog) (res sdk.Result, err error) {
	logEntry.Type = msg.Type()
	logEntry.Sender = msg.Sender.String()
	logEntry.BlockNum = msg.BlockNumber

	latestBlock := keeper.GetLatestBlock(ctx)
	if msg.BlockNumber < latestBlock.Number {
		return res, fmt.Errorf("smaller than latest blknum %d", latestBlock.Number)
	}

	block, err := keeper.ethClient.Client.HeaderByNumber(context.Background(), big.NewInt(int64(msg.BlockNumber)))
	if err != nil {
		return res, fmt.Errorf("Failed to query mainchain block by number: %s", err)
	}

	timeDiff := ctx.BlockTime().Unix() - int64(block.Time)
	if timeDiff < keeper.BlkTimeDiffLower(ctx) || timeDiff > keeper.BlkTimeDiffUpper(ctx) {
		return res, fmt.Errorf("out of bound diff with mainchain timestamp %d", ctx.BlockTime().Unix())
	}

	keeper.SyncBlock(ctx, msg.BlockNumber)
	return
}
