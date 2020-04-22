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

	threshold := keeper.GetTallyParams(ctx).Threshold.MulInt(totalToken).TruncateInt()
	// fetch active changes whose voting periods have ended (are passed the block time)
	keeper.IterateActiveChangesQueue(ctx, ctx.BlockHeader().Time, func(change Change) bool {
		tagValue := types.AttributeValueChangeFailed
		totalVote := sdk.ZeroInt()

		for _, voter := range change.Voters {
			totalVote = totalVote.Add(validatorsByAddr[voter.String()].Tokens)
		}

		change.Status = StatusFailed
		tagValue = types.AttributeValueChangeFailed

		if totalVote.GTE(threshold) {
			change.Status = StatusPassed
			tagValue = types.AttributeValueChangePassed
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
