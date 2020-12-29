package keeper

import (
	"fmt"

	"github.com/celer-network/sgn/x/sync/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// SubmitChange create new change given a content
func (keeper Keeper) SubmitChange(ctx sdk.Context, changeType string, data []byte, blockNum uint64, initiatorAddr sdk.AccAddress) (types.Change, error) {
	changeID, err := keeper.GetChangeID(ctx)
	if err != nil {
		return types.Change{}, err
	}

	submitTime := ctx.BlockHeader().Time
	votingPeriod := keeper.GetVotingParams(ctx).VotingPeriod
	change := types.NewChange(changeID, changeType, data, blockNum, submitTime, submitTime.Add(votingPeriod), initiatorAddr)
	change.Rewardable = keeper.checkRewardable(ctx, change)

	keeper.SetChange(ctx, change)
	keeper.SetChangeID(ctx, changeID+1)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSync,
			sdk.NewAttribute(sdk.AttributeKeyAction, types.ActionSubmitChange),
			sdk.NewAttribute(types.AttributeKeyChangeID, fmt.Sprintf("%d", changeID)),
		),
	)

	return change, nil
}

func (keeper Keeper) checkRewardable(ctx sdk.Context, change types.Change) bool {
	syncer := keeper.validatorKeeper.GetSyncer(ctx)

	if !syncer.ValidatorAddr.Equals(change.Initiator) {
		return false
	}

	changeType := change.Type

	return changeType == types.UpdateSidechainAddr || changeType == types.GuardTrigger ||
		changeType == types.SyncDelegator || changeType == types.SyncValidator
}

// ApproveChange adds a vote on a specific change
func (keeper Keeper) ApproveChange(ctx sdk.Context, changeID uint64, voterAddr sdk.AccAddress) error {
	change, ok := keeper.GetChange(ctx, changeID)
	if !ok {
		return sdkerrors.Wrapf(types.ErrUnknownChange, "%d", changeID)
	}

	if change.Status != types.StatusActive {
		// Exit if the change has been approved or expired
		return nil
	}

	for _, voter := range change.Voters {
		if voter.Equals(voterAddr) {
			return sdkerrors.Wrapf(types.ErrDoubleVote, "%d", changeID)
		}
	}

	change.Voters = append(change.Voters, sdk.ValAddress(voterAddr))
	keeper.SetChange(ctx, change)

	return nil
}

// GetChange get change from store by ID
func (keeper Keeper) GetChange(ctx sdk.Context, changeID uint64) (change types.Change, ok bool) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(types.ChangeKey(changeID))
	if bz == nil {
		return
	}
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &change)
	return change, true
}

// SetChange set a change to store
func (keeper Keeper) SetChange(ctx sdk.Context, change types.Change) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(change)
	store.Set(types.ChangeKey(change.ID), bz)
}

// RemoveChange removes a change based on id
func (keeper Keeper) RemoveChange(ctx sdk.Context, changeID uint64) {
	store := ctx.KVStore(keeper.storeKey)
	store.Delete(types.ChangeKey(changeID))
}

// IterateChanges iterates over the all the changes and performs a callback function
func (keeper Keeper) IterateChanges(ctx sdk.Context, cb func(change types.Change) (stop bool)) {
	store := ctx.KVStore(keeper.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ChangesKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var change types.Change
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &change)

		if cb(change) {
			break
		}
	}
}

// GetChanges returns all the changes from store
func (keeper Keeper) GetChanges(ctx sdk.Context) (changes types.Changes) {
	keeper.IterateChanges(ctx, func(change types.Change) bool {
		changes = append(changes, change)
		return false
	})
	return
}

// GetChangeID gets the highest change ID
func (keeper Keeper) GetChangeID(ctx sdk.Context) (changeID uint64, err error) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(types.ChangeIDKey)
	if bz == nil {
		return 0, sdkerrors.Wrap(types.ErrInvalidGenesis, "initial change ID hasn't been set")
	}

	changeID = types.GetChangeIDFromBytes(bz)
	return changeID, nil
}

// SetChangeID sets the new change ID to the store
func (keeper Keeper) SetChangeID(ctx sdk.Context, changeID uint64) {
	store := ctx.KVStore(keeper.storeKey)
	store.Set(types.ChangeIDKey, types.GetChangeIDBytes(changeID))
}
