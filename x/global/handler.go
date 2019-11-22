package global

import (
	"context"
	"fmt"
	"math/big"

	"github.com/celer-network/goutils/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "global" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSyncBlock:
			return handleMsgSyncBlock(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized global Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle a message to sync block
func handleMsgSyncBlock(ctx sdk.Context, keeper Keeper, msg MsgSyncBlock) sdk.Result {
	log.Infoln("Handle message to sync mainchain block number", msg.BlockNumber)

	latestBlock := keeper.GetLatestBlock(ctx)
	if msg.BlockNumber < latestBlock.Number {
		errMsg := fmt.Sprintf("Mainchain block number is smaller than current latest block, recv %d latest %d",
			msg.BlockNumber, latestBlock.Number)
		log.Error(errMsg)
		return sdk.ErrInternal(errMsg).Result()
	}

	head, err := keeper.ethClient.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to query mainchain header: %s", err)
		log.Error(errMsg)
		return sdk.ErrInternal(errMsg).Result()
	}

	blockDiff := new(big.Int).Sub(head.Number, new(big.Int).SetUint64(msg.BlockNumber))
	if blockDiff.CmpAbs(big.NewInt(keeper.MaxBlockDiff(ctx))) > 0 {
		errMsg := fmt.Sprintf("Mainchain block number is out of bound, recv %d head %d",
			msg.BlockNumber, head.Number.Uint64())
		log.Error(errMsg)
		return sdk.ErrInternal(errMsg).Result()
	}

	keeper.SyncBlock(ctx, msg.BlockNumber)
	return sdk.Result{}
}
