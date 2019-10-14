package slash

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/celer-network/sgn/x/validator"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/tendermint/tendermint/crypto"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey        sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc             *codec.Codec // The wire codec for binary encoding/decoding.
	validatorKeeper validator.Keeper
	paramstore      params.Subspace
}

// NewKeeper creates new instances of the slash Keeper
func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec, validatorKeeper validator.Keeper, paramstore params.Subspace) Keeper {
	return Keeper{
		storeKey:        storeKey,
		cdc:             cdc,
		validatorKeeper: validatorKeeper,
		paramstore:      paramstore.WithKeyTable(ParamKeyTable()),
	}
}

// HandleGuardFailure handles a validator fails to guard state.
func (k Keeper) HandleGuardFailure(ctx sdk.Context, guardAddr, reportAddr sdk.AccAddress) {
	guardValAddr := sdk.ValAddress(guardAddr)
	guardValidator, found := k.validatorKeeper.GetValidator(ctx, guardValAddr)
	if !found {
		return
	}

	reportValAddr := sdk.ValAddress(reportAddr)
	reportValidator, found := k.validatorKeeper.GetValidator(ctx, reportValAddr)
	if !found {
		return
	}

	var beneficiaries []AccountPercentPair
	beneficiaries = append(beneficiaries, NewAccountPercentPair(reportValidator.Description.Identity, k.SlashFractionGuardFailure(ctx)))

	k.Slash(ctx, AttributeValueGuardFailure, guardValidator, guardValidator.GetConsensusPower(), k.SlashFractionGuardFailure(ctx), []AccountPercentPair{})
}

// HandleDoubleSign handles a validator signing two blocks at the same height.
// power: power of the double-signing validator at the height of infraction
func (k Keeper) HandleDoubleSign(ctx sdk.Context, addr crypto.Address, power int64) {
	logger := ctx.Logger()
	consAddr := sdk.ConsAddress(addr)
	validator, found := k.validatorKeeper.GetValidatorByConsAddr(ctx, consAddr)
	if !found {
		return
	}

	logger.Info(fmt.Sprintf("Confirmed double sign from %s", consAddr))
	k.Slash(ctx, types.AttributeValueDoubleSign, validator, power, k.SlashFractionDoubleSign(ctx), []AccountPercentPair{})
}

// HandleValidatorSignature handles a validator signature, must be called once per validator per block.
func (k Keeper) HandleValidatorSignature(ctx sdk.Context, addr crypto.Address, power int64, signed bool) {
	logger := ctx.Logger()
	height := ctx.BlockHeight()
	consAddr := sdk.ConsAddress(addr)
	signInfo, found := k.GetValidatorSigningInfo(ctx, consAddr)
	if !found {
		signInfo = types.NewValidatorSigningInfo(
			consAddr,
			ctx.BlockHeight(),
			0,
			time.Unix(0, 0),
			false,
			0,
		)
	}

	// this is a relative index, so it counts blocks the validator *should* have signed
	// will use the 0-value default signing info if not present, except for start height
	signedBlocksWindow := k.SignedBlocksWindow(ctx)
	index := signInfo.IndexOffset % signedBlocksWindow
	signInfo.IndexOffset++

	// Update signed block bit array & counter
	// This counter just tracks the sum of the bit array
	// That way we avoid needing to read/write the whole array each time
	previous := k.GetValidatorMissedBlockBitArray(ctx, consAddr, index)
	missed := !signed
	switch {
	case !previous && missed:
		// Array value has changed from not missed to missed, increment counter
		k.SetValidatorMissedBlockBitArray(ctx, consAddr, index, true)
		signInfo.MissedBlocksCounter++
	case previous && !missed:
		// Array value has changed from missed to not missed, decrement counter
		k.SetValidatorMissedBlockBitArray(ctx, consAddr, index, false)
		signInfo.MissedBlocksCounter--
	default:
		// Array value at this index has not changed, no need to update counter
	}

	minHeight := signInfo.StartHeight + signedBlocksWindow
	maxMissed := signedBlocksWindow - k.MinSignedPerWindow(ctx).MulInt64(signedBlocksWindow).RoundInt64()

	// if we are past the minimum height and the validator has missed too many blocks, punish them
	if height > minHeight && signInfo.MissedBlocksCounter > maxMissed {
		validator, found := k.validatorKeeper.GetValidatorByConsAddr(ctx, consAddr)
		if !found {
			return
		}

		// Downtime confirmed: slash the validator
		logger.Info(fmt.Sprintf("Validator %s past min height of %d and above max miss threshold of %d",
			consAddr, minHeight, maxMissed))

		// We need to reset the counter & array so that the validator won't be immediately slashed for downtime upon rebonding.
		signInfo.MissedBlocksCounter = 0
		signInfo.IndexOffset = 0
		k.ClearValidatorMissedBlockBitArray(ctx, consAddr)
		k.Slash(ctx, types.AttributeValueMissingSignature, validator, power, k.SlashFractionDowntime(ctx), []AccountPercentPair{})
	}

	k.SetValidatorSigningInfo(ctx, signInfo)
}

// Slash a validator for an infraction
// Find the contributing stake and burn the specified slashFactor of it
func (k Keeper) Slash(ctx sdk.Context, reason string, validator staking.Validator, power int64, slashFactor sdk.Dec, beneficiaries []AccountPercentPair) {
	logger := ctx.Logger()

	if slashFactor.IsNegative() {
		panic(fmt.Errorf("attempted to slash with a negative slash factor: %v", slashFactor))
	}

	// Amount of slashing = slash slashFactor * power at time of infraction
	amount := sdk.TokensFromConsensusPower(power)
	slashAmount := amount.ToDec().Mul(slashFactor).TruncateInt()
	logger.Info(fmt.Sprintf(
		"validator %s slashed by %s with slash factor of %s",
		validator.GetOperator(), slashAmount, slashFactor.String()))

	candidate, found := k.validatorKeeper.GetCandidate(ctx, validator.Description.Identity)
	if !found {
		logger.Error("Cannot find candidate profile for validator", validator.Description.Identity)
	}

	penalty := NewPenalty(k.GetNextPenaltyNonce(ctx), reason, validator.Description.Identity)
	for _, delegator := range candidate.Delegators {
		penaltyAmt := slashAmount.Mul(delegator.DelegatedStake).Quo(candidate.StakingPool)
		accountAmtPair := NewAccountAmtPair(delegator.EthAddress, penaltyAmt)
		penalty.PenalizedDelegators = append(penalty.PenalizedDelegators, accountAmtPair)
	}

	penalty.Beneficiaries = beneficiaries
	penalty.GenerateProtoBytes()
	k.SetPenalty(ctx, penalty)

	// TODO: set penalty properly
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			EventTypeSlash,
			sdk.NewAttribute(sdk.AttributeKeyAction, ActionPenalty),
			sdk.NewAttribute(AttributeKeyNonce, sdk.NewUint(penalty.Nonce).String()),
			sdk.NewAttribute(types.AttributeKeyReason, reason),
		),
	)
}

// Gets the next Penalty nonce, and increment nonce by 1
func (k Keeper) GetNextPenaltyNonce(ctx sdk.Context) (nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(PenaltyNonceKey)

	if bz != nil {
		nonce = binary.BigEndian.Uint64(bz)
	}

	store.Set(PenaltyNonceKey, sdk.Uint64ToBigEndian(nonce+1))
	return
}

// Gets the entire Penalty metadata for a nonce
func (k Keeper) GetPenalty(ctx sdk.Context, nonce uint64) (penalty Penalty, found bool) {
	store := ctx.KVStore(k.storeKey)

	if !store.Has(GetPenaltyKey(nonce)) {
		return penalty, false
	}

	value := store.Get(GetPenaltyKey(nonce))
	k.cdc.MustUnmarshalBinaryBare(value, &penalty)
	return penalty, true
}

// Sets the entire Penalty metadata for a nonce
func (k Keeper) SetPenalty(ctx sdk.Context, penalty Penalty) {
	store := ctx.KVStore(k.storeKey)
	store.Set(GetPenaltyKey(penalty.Nonce), k.cdc.MustMarshalBinaryBare(penalty))
}

// Stored by *validator* address (not operator address)
func (k Keeper) GetValidatorSigningInfo(ctx sdk.Context, address sdk.ConsAddress) (info types.ValidatorSigningInfo, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetValidatorSigningInfoKey(address))
	if bz == nil {
		found = false
		return
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &info)
	found = true
	return
}

// Stored by *validator* address (not operator address)
func (k Keeper) SetValidatorSigningInfo(ctx sdk.Context, info types.ValidatorSigningInfo) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(info)
	store.Set(types.GetValidatorSigningInfoKey(info.Address), bz)
}

// Stored by *validator* address (not operator address)
func (k Keeper) GetValidatorMissedBlockBitArray(ctx sdk.Context, address sdk.ConsAddress, index int64) (missed bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetValidatorMissedBlockBitArrayKey(address, index))
	if bz == nil {
		// lazy: treat empty key as not missed
		missed = false
		return
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &missed)
	return
}

// Stored by *validator* address (not operator address)
func (k Keeper) SetValidatorMissedBlockBitArray(ctx sdk.Context, address sdk.ConsAddress, index int64, missed bool) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(missed)
	store.Set(types.GetValidatorMissedBlockBitArrayKey(address, index), bz)
}

// Stored by *validator* address (not operator address)
func (k Keeper) ClearValidatorMissedBlockBitArray(ctx sdk.Context, address sdk.ConsAddress) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetValidatorMissedBlockBitArrayPrefixKey(address))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
}
