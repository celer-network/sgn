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
	var tagValue string
	validators := keeper.GetValidators(ctx)
	totalToken := sdk.ZeroInt()
	validatorsByAddr := map[string]staking.Validator{}

	for _, validator := range validators {
		totalToken = totalToken.Add(validator.Tokens)
		validatorsByAddr[validator.OperatorAddress.String()] = validator
	}

	threshold := keeper.GetTallyParams(ctx).Threshold.MulInt(totalToken).TruncateInt()

	activeChanges := keeper.GetActiveChanges(ctx)
	for _, change := range activeChanges {
		totalVote := sdk.ZeroInt()
		for _, voter := range change.Voters {
			validator, ok := validatorsByAddr[voter.String()]
			if !ok {
				continue
			}
			totalVote = totalVote.Add(validator.Tokens)
		}

		if totalVote.GTE(threshold) {
			log.Infoln("Change type", change.Type, totalVote, threshold)

			err := keeper.ApplyChange(ctx, change)
			if err != nil {
				log.Errorln("Apply change err:", err)
				change.Status = StatusFailed
				tagValue = types.AttributeValueChangeFailed
			} else {
				change.Status = StatusPassed
				tagValue = types.AttributeValueChangePassed
			}
		}

		if ctx.BlockTime().After(change.VotingEndTime) {
			change.Status = StatusFailed
			tagValue = types.AttributeValueChangeFailed
		}

		if change.Status != StatusActive {
			keeper.SetChange(ctx, change)
			keeper.RemoveFromActiveChangeQueue(ctx, change.ID)

			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeActiveChange,
					sdk.NewAttribute(types.AttributeKeyChangeID, fmt.Sprintf("%d", change.ID)),
					sdk.NewAttribute(types.AttributeKeyChangeResult, tagValue),
				),
			)
		}
	}
}
