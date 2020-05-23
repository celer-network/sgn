package sync

import (
	"fmt"

	"github.com/celer-network/goutils/log"
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
		totalVote := sdk.ZeroInt()
		for _, voter := range change.Voters {
			validator, ok := validatorsByAddr[voter.String()]
			if !ok {
				continue
			}
			totalVote = totalVote.Add(validator.Tokens)
		}

		tagValue := types.AttributeValueChangeFailed
		change.Status = StatusFailed
		log.Infoln("Change type", change.Type, totalVote, threshold)
		if totalVote.GTE(threshold) {
			err := keeper.ApplyChange(ctx, change)
			if err != nil {
				log.Errorln("Apply change err:", err)
			} else {
				change.Status = StatusPassed
				tagValue = types.AttributeValueChangePassed
			}
		}

		keeper.SetChange(ctx, change)
		keeper.RemoveFromActiveChangeQueue(ctx, change.ID, change.VotingEndTime)

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeActiveChange,
				sdk.NewAttribute(types.AttributeKeyChangeID, fmt.Sprintf("%d", change.ID)),
				sdk.NewAttribute(types.AttributeKeyChangeResult, tagValue),
			),
		)

		return false
	})
}
