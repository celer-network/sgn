package keeper

import (
	"fmt"
	"time"

	"github.com/celer-network/sgn/x/sync/types"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// SubmitChange create new change given a content
func (keeper Keeper) SubmitChange(ctx sdk.Context, changeType string, data []byte, initiatorAddr sdk.AccAddress) (types.Change, error) {
	changeID, err := keeper.GetChangeID(ctx)
	if err != nil {
		return types.Change{}, err
	}

	submitTime := ctx.BlockHeader().Time
	votingPeriod := keeper.GetVotingParams(ctx).VotingPeriod
	change := types.NewChange(changeID, changeType, data, submitTime, submitTime.Add(votingPeriod), initiatorAddr)

	keeper.SetChange(ctx, change)
	keeper.InsertActiveChangeQueue(ctx, changeID, change.VotingEndTime)
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

// ApproveChange adds a vote on a specific change
func (keeper Keeper) ApproveChange(ctx sdk.Context, changeID uint64, voterAddr sdk.AccAddress) error {
	change, ok := keeper.GetChange(ctx, changeID)
	if !ok {
		return sdkerrors.Wrapf(types.ErrUnknownChange, "%d", changeID)
	}
	if change.Status != types.StatusVotingPeriod {
		return sdkerrors.Wrapf(types.ErrInactiveChange, "%d", changeID)
	}

	change.Voters = append(change.Voters, sdk.ValAddress(voterAddr))
	keeper.SetChange(ctx, change)

	return nil
}

// GetChange get change from store by ChangeID
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
	store.Set(types.ChangeKey(change.ChangeID), bz)
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

// GetChangesFiltered retrieves changes filtered by a given set of params which
// include pagination parameters along with voter and depositor addresses and a
// change status. The voter address will filter changes by whether or not
// that address has voted on changes. The depositor address will filter changes
// by whether or not that address has deposited to them. Finally, status will filter
// changes by status.
//
// NOTE: If no filters are provided, all changes will be returned in paginated
// form.
func (keeper Keeper) GetChangesFiltered(ctx sdk.Context, params types.QueryChangesParams) []types.Change {
	changes := keeper.GetChanges(ctx)
	filteredChanges := make([]types.Change, 0, len(changes))

	for _, p := range changes {
		matchStatus := true

		// match status (if supplied/valid)
		if types.ValidChangeStatus(params.ChangeStatus) {
			matchStatus = p.Status == params.ChangeStatus
		}

		if matchStatus {
			filteredChanges = append(filteredChanges, p)
		}
	}

	start, end := client.Paginate(len(filteredChanges), params.Page, params.Limit, 100)
	if start < 0 || end < 0 {
		filteredChanges = []types.Change{}
	} else {
		filteredChanges = filteredChanges[start:end]
	}

	return filteredChanges
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

// InsertActiveChangeQueue inserts a ChangeID into the active change queue at endTime
func (keeper Keeper) InsertActiveChangeQueue(ctx sdk.Context, changeID uint64, endTime time.Time) {
	store := ctx.KVStore(keeper.storeKey)
	bz := types.GetChangeIDBytes(changeID)
	store.Set(types.ActiveChangeQueueKey(changeID, endTime), bz)
}

// RemoveFromActiveChangeQueue removes a changeID from the Active Change Queue
func (keeper Keeper) RemoveFromActiveChangeQueue(ctx sdk.Context, changeID uint64, endTime time.Time) {
	store := ctx.KVStore(keeper.storeKey)
	store.Delete(types.ActiveChangeQueueKey(changeID, endTime))
}

// IterateActiveChangesQueue iterates over the changes in the active change queue
// and performs a callback function
func (keeper Keeper) IterateActiveChangesQueue(ctx sdk.Context, endTime time.Time, cb func(change types.Change) (stop bool)) {
	iterator := keeper.ActiveChangeQueueIterator(ctx, endTime)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		changeID, _ := types.SplitActiveChangeQueueKey(iterator.Key())
		change, found := keeper.GetChange(ctx, changeID)
		if !found {
			panic(fmt.Sprintf("change %d does not exist", changeID))
		}

		if cb(change) {
			break
		}
	}
}

// ActiveChangeQueueIterator returns an sdk.Iterator for all the changes in the Active Queue that expire by endTime
func (keeper Keeper) ActiveChangeQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)
	return store.Iterator(types.ActiveChangeQueuePrefix, sdk.PrefixEndBytes(types.ActiveChangeByTimeKey(endTime)))
}