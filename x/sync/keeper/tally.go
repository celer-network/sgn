package keeper

import (
	"github.com/celer-network/sgn/x/sync/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/exported"
)

// TODO: Break into several smaller functions for clarity

// Tally iterates over the votes and updates the tally of a change based on the voting power of the
// voters
func (keeper Keeper) Tally(ctx sdk.Context, change types.Change) (passes bool, burnDeposits bool, tallyResults types.TallyResult) {
	results := make(map[types.VoteOption]sdk.Int)
	results[types.OptionYes] = sdk.ZeroInt()
	results[types.OptionAbstain] = sdk.ZeroInt()
	results[types.OptionNo] = sdk.ZeroInt()
	results[types.OptionNoWithVeto] = sdk.ZeroInt()

	totalVotingPower := sdk.ZeroInt()
	totalBondedTokens := sdk.ZeroInt()
	currValidators := make(map[string]types.ValidatorGovInfo)

	// fetch all the bonded validators, insert them into currValidators
	keeper.vk.IterateBondedValidatorsByPower(ctx, func(index int64, validator exported.ValidatorI) (stop bool) {
		currValidators[validator.GetOperator().String()] = types.NewValidatorGovInfo(
			validator.GetOperator(),
			validator.GetBondedTokens(),
			types.OptionEmpty,
		)

		totalBondedTokens = totalBondedTokens.Add(validator.GetBondedTokens())
		return false
	})

	keeper.IterateVotes(ctx, change.ChangeID, func(vote types.Vote) bool {
		// if validator, just record it in the map
		valAddrStr := sdk.ValAddress(vote.Voter).String()
		if val, ok := currValidators[valAddrStr]; ok {
			val.Vote = vote.Option
			currValidators[valAddrStr] = val
		}

		keeper.deleteVote(ctx, vote.ChangeID, vote.Voter)
		return false
	})

	// iterate over the validators again to tally their voting power
	for _, val := range currValidators {
		if val.Vote == types.OptionEmpty {
			continue
		}

		results[val.Vote] = results[val.Vote].Add(val.BondedTokens)
		totalVotingPower = totalVotingPower.Add(val.BondedTokens)
	}

	tallyParams := keeper.GetTallyParams(ctx)
	tallyResults = types.NewTallyResultFromMap(results)

	// If there is not enough quorum of votes, the change fails
	percentVoting := totalVotingPower.ToDec().QuoInt(totalBondedTokens)
	if percentVoting.LT(tallyParams.Quorum) {
		return false, true, tallyResults
	}

	// If no one votes (everyone abstains), change fails
	if totalVotingPower.Equal(results[types.OptionAbstain]) {
		return false, false, tallyResults
	}

	// If more than 1/3 of voters veto, change fails
	if results[types.OptionNoWithVeto].ToDec().QuoInt(totalVotingPower).GT(tallyParams.Veto) {
		return false, true, tallyResults
	}

	// If more than 1/2 of non-abstaining voters vote Yes, change passes
	if results[types.OptionYes].ToDec().QuoInt(totalVotingPower.Sub(results[types.OptionAbstain])).GT(tallyParams.Threshold) {
		return true, false, tallyResults
	}

	// If more than 1/2 of non-abstaining voters vote No, change fails
	return false, false, tallyResults
}
