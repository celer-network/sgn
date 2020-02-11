package validator

import (
	"bytes"
	"fmt"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	abci "github.com/tendermint/tendermint/abci/types"
)

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, req abci.RequestEndBlock, keeper Keeper) (updates []abci.ValidatorUpdate) {
	setPuller(ctx, req, keeper)
	setPusher(ctx, keeper)
	miningReward := keeper.MiningReward(ctx)
	keeper.DistributeReward(ctx, miningReward, MiningReward)

	return applyAndReturnValidatorSetUpdates(ctx, keeper)
}

// Update puller for every pullerDuration
func setPuller(ctx sdk.Context, req abci.RequestEndBlock, keeper Keeper) {
	puller := keeper.GetPuller(ctx)
	validators := keeper.GetValidators(ctx)
	pullerDuration := keeper.PullerDuration(ctx)
	vIdx := uint(req.Height) / pullerDuration % uint(len(validators))

	if puller.ValidatorIdx != vIdx || puller.ValidatorAddr.Empty() {
		puller = NewPuller(vIdx, sdk.AccAddress(validators[vIdx].OperatorAddress))
		keeper.SetPuller(ctx, puller)
	}
}

// Update pusher for every pusherDuration
func setPusher(ctx sdk.Context, keeper Keeper) {
	pusher := keeper.GetPusher(ctx)
	validators := keeper.GetValidators(ctx)
	latestBlock := keeper.globalKeeper.GetLatestBlock(ctx)
	pusherDuration := keeper.PusherDuration(ctx)
	vIdx := uint(latestBlock.Number) / pusherDuration % uint(len(validators))

	if pusher.ValidatorIdx != vIdx || pusher.ValidatorAddr.Empty() {
		pusher = NewPusher(vIdx, sdk.AccAddress(validators[vIdx].OperatorAddress))
		keeper.SetPusher(ctx, pusher)
	}
}

// Modify based on https://github.com/cosmos/cosmos-sdk/blob/v0.37.5/x/staking/keeper/val_state_change.go#L26
//
// Apply and return accumulated updates to the bonded validator set. Also,
// * Updates the active valset as keyed by LastValidatorPowerKey.
// * Updates the total power as keyed by LastTotalPowerKey.
// * Updates relevant indices.
//
//
// CONTRACT: Only validators with non-zero power or zero-power that were bonded
// at the previous block height or were removed from the validator set entirely
// are returned to Tendermint.
func applyAndReturnValidatorSetUpdates(ctx sdk.Context, keeper Keeper) (updates []abci.ValidatorUpdate) {
	stakingKeeper := keeper.stakingKeeper
	maxValidators := stakingKeeper.GetParams(ctx).MaxValidators
	totalPower := sdk.ZeroInt()

	// Retrieve the last validator set.
	// The persistent set is updated later in this function.
	// (see LastValidatorPowerKey).
	last := getLastValidatorsByAddr(ctx, stakingKeeper)

	// Iterate over validators, highest power to lowest.
	iterator := stakingKeeper.ValidatorsPowerStoreIterator(ctx)
	defer iterator.Close()
	for count := 0; iterator.Valid() && count < int(maxValidators); iterator.Next() {
		// everything that is iterated in this loop is becoming or already a
		// part of the bonded validator set

		valAddr := sdk.ValAddress(iterator.Value())
		validator, found := stakingKeeper.GetValidator(ctx, valAddr)
		if !found {
			panic(fmt.Sprintf("validator record not found for address: %X\n", valAddr))
		}

		if validator.Jailed {
			panic("should never retrieve a jailed validator from the power store")
		}

		// if we get to a zero-power validator (which we don't bond),
		// there are no more possible bonded validators
		if validator.PotentialConsensusPower() == 0 {
			break
		}

		// fetch the old power bytes
		var valAddrBytes [sdk.AddrLen]byte
		copy(valAddrBytes[:], valAddr[:])
		oldPowerBytes, found := last[valAddrBytes]

		newPower := validator.ConsensusPower()
		newPowerBytes := keeper.cdc.MustMarshalBinaryLengthPrefixed(newPower)

		// update the validator set if power has changed
		if !found || !bytes.Equal(oldPowerBytes, newPowerBytes) {
			updates = append(updates, validator.ABCIValidatorUpdate())

			stakingKeeper.SetLastValidatorPower(ctx, valAddr, newPower)
		}

		delete(last, valAddrBytes)

		count++
		totalPower = totalPower.Add(sdk.NewInt(newPower))
	}

	noLongerBonded := sortNoLongerBonded(last)
	for _, valAddrBytes := range noLongerBonded {
		validator, found := stakingKeeper.GetValidator(ctx, sdk.ValAddress(valAddrBytes))
		if !found {
			panic(fmt.Sprintf("validator record not found for address: %X\n", sdk.ValAddress(valAddrBytes)))
		}

		stakingKeeper.DeleteLastValidatorPower(ctx, validator.GetOperator())
		updates = append(updates, validator.ABCIValidatorUpdateZero())
	}

	// set total power on lookup index if there are any updates
	if len(updates) > 0 {
		stakingKeeper.SetLastTotalPower(ctx, totalPower)
	}

	return updates
}

type validatorsByAddr map[[sdk.AddrLen]byte][]byte

func getLastValidatorsByAddr(ctx sdk.Context, k staking.Keeper) validatorsByAddr {
	last := make(validatorsByAddr)
	iterator := k.LastValidatorsIterator(ctx)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var valAddr [sdk.AddrLen]byte
		// extract the validator address from the key (prefix is 1-byte)
		copy(valAddr[:], iterator.Key()[1:])
		powerBytes := iterator.Value()
		last[valAddr] = make([]byte, len(powerBytes))
		copy(last[valAddr], powerBytes)
	}

	return last
}

// given a map of remaining validators to previous bonded power
// returns the list of validators to be unbonded, sorted by operator address
func sortNoLongerBonded(last validatorsByAddr) [][]byte {
	// sort the map keys for determinism
	noLongerBonded := make([][]byte, len(last))
	index := 0
	for valAddrBytes := range last {
		valAddr := make([]byte, sdk.AddrLen)
		copy(valAddr, valAddrBytes[:])
		noLongerBonded[index] = valAddr
		index++
	}
	// sorted by address - order doesn't matter
	sort.SliceStable(noLongerBonded, func(i, j int) bool {
		// -1 means strictly less than
		return bytes.Compare(noLongerBonded[i], noLongerBonded[j]) == -1
	})
	return noLongerBonded
}
