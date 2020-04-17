package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/celer-network/sgn/x/sync/types"
)

// SubmitChange create new change given a content
func (keeper Keeper) SubmitChange(ctx sdk.Context, content types.Content) (types.Change, error) {
	if !keeper.router.HasRoute(content.ChangeRoute()) {
		return types.Change{}, sdkerrors.Wrap(types.ErrNoChangeHandlerExists, content.ChangeRoute())
	}

	// Execute the change content in a cache-wrapped context to validate the
	// actual parameter changes before the change proceeds through the
	// sync process. State is not persisted.
	cacheCtx, _ := ctx.CacheContext()
	handler := keeper.router.GetRoute(content.ChangeRoute())
	if err := handler(cacheCtx, content); err != nil {
		return types.Change{}, sdkerrors.Wrap(types.ErrInvalidChangeContent, err.Error())
	}

	changeID, err := keeper.GetChangeID(ctx)
	if err != nil {
		return types.Change{}, err
	}

	submitTime := ctx.BlockHeader().Time
	depositPeriod := keeper.GetDepositParams(ctx).MaxDepositPeriod

	change := types.NewChange(content, changeID, submitTime, submitTime.Add(depositPeriod))

	keeper.SetChange(ctx, change)
	keeper.InsertInactiveChangeQueue(ctx, changeID, change.DepositEndTime)
	keeper.SetChangeID(ctx, changeID+1)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSubmitChange,
			sdk.NewAttribute(types.AttributeKeyChangeID, fmt.Sprintf("%d", changeID)),
		),
	)

	return change, nil
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

// DeleteChange deletes a change from store
func (keeper Keeper) DeleteChange(ctx sdk.Context, changeID uint64) {
	store := ctx.KVStore(keeper.storeKey)
	change, ok := keeper.GetChange(ctx, changeID)
	if !ok {
		panic(fmt.Sprintf("couldn't find change with id#%d", changeID))
	}
	keeper.RemoveFromInactiveChangeQueue(ctx, changeID, change.DepositEndTime)
	keeper.RemoveFromActiveChangeQueue(ctx, changeID, change.VotingEndTime)
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
		matchVoter, matchDepositor, matchStatus := true, true, true

		// match status (if supplied/valid)
		if types.ValidChangeStatus(params.ChangeStatus) {
			matchStatus = p.Status == params.ChangeStatus
		}

		// match voter address (if supplied)
		if len(params.Voter) > 0 {
			_, matchVoter = keeper.GetVote(ctx, p.ChangeID, params.Voter)
		}

		// match depositor (if supplied)
		if len(params.Depositor) > 0 {
			_, matchDepositor = keeper.GetDeposit(ctx, p.ChangeID, params.Depositor)
		}

		if matchVoter && matchDepositor && matchStatus {
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

func (keeper Keeper) activateVotingPeriod(ctx sdk.Context, change types.Change) {
	change.VotingStartTime = ctx.BlockHeader().Time
	votingPeriod := keeper.GetVotingParams(ctx).VotingPeriod
	change.VotingEndTime = change.VotingStartTime.Add(votingPeriod)
	change.Status = types.StatusVotingPeriod
	keeper.SetChange(ctx, change)

	keeper.RemoveFromInactiveChangeQueue(ctx, change.ChangeID, change.DepositEndTime)
	keeper.InsertActiveChangeQueue(ctx, change.ChangeID, change.VotingEndTime)
}
