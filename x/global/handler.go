package global

import (
	"context"
	"fmt"
	"math/big"

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
	fmt.Sprintln("Handling MsgSyncBlock with block number of", msg.BlockNumber)

	lastestBlock := keeper.GetLatestBlock(ctx)
	if msg.BlockNumber < lastestBlock.Number {
		return sdk.ErrInternal("Block number is smaller than current latest block").Result()
	}

	head, err := keeper.ethClient.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("Failed to query mainchain header: %s", err)).Result()
	}

	blockDiff := new(big.Int).Sub(head.Number, new(big.Int).SetUint64(msg.BlockNumber))
	if blockDiff.CmpAbs(big.NewInt(keeper.MaxBlockDiff(ctx))) > 0 {
		return sdk.ErrInternal("Block number is out of bound").Result()
	}

	keeper.SyncBlock(ctx, msg.BlockNumber)
	return sdk.Result{}
}
