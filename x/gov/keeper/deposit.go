package keeper

import (
	"fmt"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// GetDeposit gets the deposit of a specific depositor on a specific proposal
func (keeper Keeper) GetDeposit(ctx sdk.Context, proposalID uint64, depositorAddr sdk.AccAddress) (deposit types.Deposit, found bool) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(types.DepositKey(proposalID, depositorAddr))
	if bz == nil {
		return deposit, false
	}

	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &deposit)
	return deposit, true
}

// SetDeposit sets a Deposit to the gov store
func (keeper Keeper) SetDeposit(ctx sdk.Context, deposit types.Deposit) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(deposit)
	store.Set(types.DepositKey(deposit.ProposalID, deposit.Depositor), bz)
}

// GetAllDeposits returns all the deposits from the store
func (keeper Keeper) GetAllDeposits(ctx sdk.Context) (deposits types.Deposits) {
	keeper.IterateAllDeposits(ctx, func(deposit types.Deposit) bool {
		deposits = append(deposits, deposit)
		return false
	})
	return
}

// GetDeposits returns all the deposits from a proposal
func (keeper Keeper) GetDeposits(ctx sdk.Context, proposalID uint64) (deposits types.Deposits) {
	keeper.IterateDeposits(ctx, proposalID, func(deposit types.Deposit) bool {
		deposits = append(deposits, deposit)
		return false
	})
	return
}

// IterateAllDeposits iterates over the all the stored deposits and performs a callback function
func (keeper Keeper) IterateAllDeposits(ctx sdk.Context, cb func(deposit types.Deposit) (stop bool)) {
	store := ctx.KVStore(keeper.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.DepositsKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var deposit types.Deposit
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &deposit)

		if cb(deposit) {
			break
		}
	}
}

// IterateDeposits iterates over the all the proposals deposits and performs a callback function
func (keeper Keeper) IterateDeposits(ctx sdk.Context, proposalID uint64, cb func(deposit types.Deposit) (stop bool)) {
	store := ctx.KVStore(keeper.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.DepositsKey(proposalID))

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var deposit types.Deposit
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &deposit)

		if cb(deposit) {
			break
		}
	}
}

// GetAccTotalDeposit gets the AccTotalDeposit by address
func (keeper Keeper) GetAccTotalDeposit(ctx sdk.Context, depositorAddr sdk.AccAddress) (accTotalDeposit types.AccTotalDeposit, found bool) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(types.DepositorKey(depositorAddr))
	if bz == nil {
		return accTotalDeposit, false
	}

	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &accTotalDeposit)
	return accTotalDeposit, true
}

// SetAccTotalDeposit sets a AccTotalDeposit to the gov store
func (keeper Keeper) SetAccTotalDeposit(ctx sdk.Context, depositorAddr sdk.AccAddress, accTotalDeposit types.AccTotalDeposit) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(accTotalDeposit)
	store.Set(types.DepositorKey(depositorAddr), bz)
}

// AddDeposit adds or updates a deposit of a specific depositor on a specific proposal
// Activates voting period when appropriate
func (keeper Keeper) AddDeposit(ctx sdk.Context, proposalID uint64, depositorAddr sdk.AccAddress, depositAmount sdk.Int) (bool, error) {
	// Checks to see if proposal exists
	proposal, ok := keeper.GetProposal(ctx, proposalID)
	if !ok {
		return false, sdkerrors.Wrapf(types.ErrUnknownProposal, "%d", proposalID)
	}

	// Check if proposal is still depositable
	if (proposal.Status != types.StatusDepositPeriod) && (proposal.Status != types.StatusVotingPeriod) {
		return false, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "%d", proposalID)
	}

	validator, found := keeper.vk.GetValidator(ctx, sdk.ValAddress(depositorAddr))
	if !found {
		return false, sdkerrors.Wrapf(types.ErrUnknownProposal, "Invalid depositor addr %s", depositorAddr)
	}

	ethAddr := validator.Description.Identity
	selfDelegator, found := keeper.vk.GetDelegator(ctx, ethAddr, ethAddr)
	if !found {
		return false, sdkerrors.Wrapf(types.ErrUnknownProposal, "Invalid depositor addr %s, %s", depositorAddr, ethAddr)
	}

	accTotalDeposit, found := keeper.GetAccTotalDeposit(ctx, depositorAddr)
	if !found {
		accTotalDeposit = types.NewAccTotalDeposit()
	}

	accTotalDeposit.Amount = accTotalDeposit.Amount.Add(depositAmount)
	if accTotalDeposit.Amount.GT(selfDelegator.DelegatedStake) {
		return false, sdkerrors.Wrapf(
			sdkerrors.ErrInvalidRequest,
			"Depositor does not have enough stake to deposit, need %s, have %s", accTotalDeposit.Amount, selfDelegator.DelegatedStake)
	}
	keeper.SetAccTotalDeposit(ctx, depositorAddr, accTotalDeposit)

	// Update proposal
	proposal.TotalDeposit = proposal.TotalDeposit.Add(depositAmount)
	keeper.SetProposal(ctx, proposal)

	// Check if deposit has provided sufficient total funds to transition the proposal into the voting period
	activatedVotingPeriod := false
	if proposal.Status == types.StatusDepositPeriod && proposal.TotalDeposit.GTE(keeper.GetDepositParams(ctx).MinDeposit) {
		keeper.activateVotingPeriod(ctx, proposal)
		activatedVotingPeriod = true
	}

	// Add or update deposit object
	deposit, found := keeper.GetDeposit(ctx, proposalID, depositorAddr)
	if found {
		deposit.Amount = deposit.Amount.Add(depositAmount)
	} else {
		deposit = types.NewDeposit(proposalID, depositorAddr, depositAmount)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeProposalDeposit,
			sdk.NewAttribute(sdk.AttributeKeyAmount, depositAmount.String()),
			sdk.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
		),
	)

	keeper.SetDeposit(ctx, deposit)
	return activatedVotingPeriod, nil
}

// RefundDeposits refunds and deletes all the deposits on a specific proposal
func (keeper Keeper) RefundDeposits(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(keeper.storeKey)

	keeper.IterateDeposits(ctx, proposalID, func(deposit types.Deposit) bool {
		accTotalDeposit, found := keeper.GetAccTotalDeposit(ctx, deposit.Depositor)
		if !found {
			log.Errorf("deposit not found. %s", deposit)
			return false
		}
		accTotalDeposit.Amount = accTotalDeposit.Amount.Sub(deposit.Amount)
		keeper.SetAccTotalDeposit(ctx, deposit.Depositor, accTotalDeposit)

		store.Delete(types.DepositKey(proposalID, deposit.Depositor))
		return false
	})
}

// DeleteDeposits deletes all the deposits on a specific proposal without refunding them
func (keeper Keeper) DeleteDeposits(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(keeper.storeKey)

	keeper.IterateDeposits(ctx, proposalID, func(deposit types.Deposit) bool {
		// TODO: properly handle delete deposits
		accTotalDeposit, found := keeper.GetAccTotalDeposit(ctx, deposit.Depositor)
		if !found {
			log.Errorf("deposit not found. %s", deposit)
			return false
		}
		accTotalDeposit.Amount = accTotalDeposit.Amount.Sub(deposit.Amount)
		keeper.SetAccTotalDeposit(ctx, deposit.Depositor, accTotalDeposit)
		keeper.sk.HandleProposalDepositBurn(ctx, deposit.Depositor, deposit.Amount)

		store.Delete(types.DepositKey(proposalID, deposit.Depositor))
		return false
	})
}
