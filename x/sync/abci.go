package sync

import (
	"fmt"

	"github.com/celer-network/sgn/x/sync/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, keeper Keeper) {
	logger := keeper.Logger(ctx)

	// delete inactive change from store and its deposits
	keeper.IterateInactiveChangesQueue(ctx, ctx.BlockHeader().Time, func(change Change) bool {
		keeper.DeleteChange(ctx, change.ChangeID)
		keeper.DeleteDeposits(ctx, change.ChangeID)

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeInactiveChange,
				sdk.NewAttribute(types.AttributeKeyChangeID, fmt.Sprintf("%d", change.ChangeID)),
				sdk.NewAttribute(types.AttributeKeyChangeResult, types.AttributeValueChangeDropped),
			),
		)

		logger.Info(
			fmt.Sprintf("change %d (%s) didn't meet minimum deposit of %s (had only %s); deleted",
				change.ChangeID,
				change.GetTitle(),
				keeper.GetDepositParams(ctx).MinDeposit,
				change.TotalDeposit,
			),
		)
		return false
	})

	// fetch active changes whose voting periods have ended (are passed the block time)
	keeper.IterateActiveChangesQueue(ctx, ctx.BlockHeader().Time, func(change Change) bool {
		var tagValue, logMsg string

		passes, burnDeposits, tallyResults := keeper.Tally(ctx, change)

		if burnDeposits {
			keeper.DeleteDeposits(ctx, change.ChangeID)
		} else {
			keeper.RefundDeposits(ctx, change.ChangeID)
		}

		if passes {
			handler := keeper.Router().GetRoute(change.ChangeRoute())
			cacheCtx, writeCache := ctx.CacheContext()

			// The change handler may execute state mutating logic depending
			// on the change content. If the handler fails, no state mutation
			// is written and the error message is logged.
			err := handler(cacheCtx, change.Content)
			if err == nil {
				change.Status = StatusPassed
				tagValue = types.AttributeValueChangePassed
				logMsg = "passed"

				// The cached context is created with a new EventManager. However, since
				// the change handler execution was successful, we want to track/keep
				// any events emitted, so we re-emit to "merge" the events into the
				// original Context's EventManager.
				ctx.EventManager().EmitEvents(cacheCtx.EventManager().Events())

				// write state to the underlying multi-store
				writeCache()
			} else {
				change.Status = StatusFailed
				tagValue = types.AttributeValueChangeFailed
				logMsg = fmt.Sprintf("passed, but failed on execution: %s", err)
			}
		} else {
			change.Status = StatusRejected
			tagValue = types.AttributeValueChangeRejected
			logMsg = "rejected"
		}

		change.FinalTallyResult = tallyResults

		keeper.SetChange(ctx, change)
		keeper.RemoveFromActiveChangeQueue(ctx, change.ChangeID, change.VotingEndTime)

		logger.Info(
			fmt.Sprintf(
				"change %d (%s) tallied; result: %s",
				change.ChangeID, change.GetTitle(), logMsg,
			),
		)

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeActiveChange,
				sdk.NewAttribute(types.AttributeKeyChangeID, fmt.Sprintf("%d", change.ChangeID)),
				sdk.NewAttribute(types.AttributeKeyChangeResult, tagValue),
			),
		)
		return false
	})
}
