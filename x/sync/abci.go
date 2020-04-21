package sync

import (
	"fmt"

	"github.com/celer-network/sgn/x/sync/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, keeper Keeper) {
	validators := keeper.GetValidators(ctx)
	totalToken := sdk.ZeroInt()
	validatorsByAddr := map[string]staking.Validator{}

	for _, validator := range validators {
		totalToken = totalToken.Add(validator.Tokens)
		validatorsByAddr[validator.OperatorAddress.String()] = validator
	}

	// fetch active changes whose voting periods have ended (are passed the block time)
	keeper.IterateActiveChangesQueue(ctx, ctx.BlockHeader().Time, func(change Change) bool {
		var tagValue string

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
		}

		keeper.SetChange(ctx, change)
		keeper.RemoveFromActiveChangeQueue(ctx, change.ChangeID, change.VotingEndTime)

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
