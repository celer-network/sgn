package validator

import (
	"bytes"
	"sort"

	"github.com/celer-network/goutils/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	abci "github.com/tendermint/tendermint/abci/types"
)

// EndBlocker called every block, update validator set, distribute rewards
func EndBlocker(ctx sdk.Context, keeper Keeper) (updates []abci.ValidatorUpdate) {
	setSyncer(ctx, keeper)
	miningReward := keeper.MiningReward(ctx)
	keeper.AddEpochMiningReward(ctx, miningReward)
	keeper.DistributeReward(ctx)

	return applyAndReturnValidatorSetUpdates(ctx, keeper)
}

// Update syncer for every syncerDuration
func setSyncer(ctx sdk.Context, keeper Keeper) {
	syncer := keeper.GetSyncer(ctx)
	validators := keeper.GetBondedValidators(ctx)
	syncerDuration := keeper.SyncerDuration(ctx)
	vIdx := uint(ctx.BlockHeight()) / syncerDuration % uint(len(validators))

	if syncer.ValidatorIdx != vIdx || syncer.ValidatorAddr.Empty() {
		syncer = NewSyncer(vIdx, sdk.AccAddress(validators[vIdx].OperatorAddress))
		keeper.SetSyncer(ctx, syncer)
		log.Infof("set syncer to %s", syncer.ValidatorAddr)
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
	var valByPowerIndexKeysToRemove [][]byte
	for count := 0; iterator.Valid() && count < int(maxValidators); iterator.Next() {
		// everything that is iterated in this loop is becoming or already a
		// part of the bonded validator set

		valAddr := sdk.ValAddress(iterator.Value())
		validator, found := stakingKeeper.GetValidator(ctx, valAddr)
		if !found {
			log.Errorf("validator record not found for address: %X, skipping", valAddr)
			valByPowerIndexKeysToRemove = append(valByPowerIndexKeysToRemove, iterator.Key())
			continue
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
	store := ctx.KVStore(keeper.storeKey)
	for _, key := range valByPowerIndexKeysToRemove {
		store.Delete(key)
	}

	noLongerBonded := sortNoLongerBonded(last)
	for _, valAddrBytes := range noLongerBonded {
		valAddr := sdk.ValAddress(valAddrBytes)
		validator, found := stakingKeeper.GetValidator(ctx, valAddr)
		if !found {
			log.Errorf("validator record not found for address: %X\n", valAddr)
			stakingKeeper.DeleteLastValidatorPower(ctx, valAddr)
			continue
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
// returns the list of validators to be unbonded, sorted by validator account address
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
