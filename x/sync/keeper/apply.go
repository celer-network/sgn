package keeper

import (
	"errors"

	"github.com/celer-network/sgn/x/global"
	"github.com/celer-network/sgn/x/sync/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (keeper Keeper) ApplyChange(ctx sdk.Context, change types.Change) error {
	switch change.Type {
	case types.SyncBlock:
		return keeper.SyncBlock(ctx, change.Data)
	default:
		return errors.New("Invalid change type")
	}
}

func (keeper Keeper) SyncBlock(ctx sdk.Context, data []byte) error {
	var block global.Block
	keeper.cdc.MustUnmarshalBinaryBare(data, &block)

	keeper.gk.SyncBlock(ctx, block.Number)
	return nil
}
