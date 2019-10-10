package slash

import (
	"fmt"
	"time"

	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/x/global"
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
	ethClient       *mainchain.EthClient
	globalKeeper    global.Keeper
	validatorKeeper validator.Keeper
	paramstore      params.Subspace
}

// NewKeeper creates new instances of the slash Keeper
func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec, ethClient *mainchain.EthClient,
	globalKeeper global.Keeper, validatorKeeper validator.Keeper, paramstore params.Subspace) Keeper {
	return Keeper{
		storeKey:        storeKey,
		cdc:             cdc,
		ethClient:       ethClient,
		globalKeeper:    globalKeeper,
		validatorKeeper: validatorKeeper,
		paramstore:      paramstore.WithKeyTable(ParamKeyTable()),
	}
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
	k.Slash(ctx, validator, power, k.SlashFractionDoubleSign(ctx), types.AttributeValueDoubleSign)
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
	maxMissed := signedBlocksWindow - k.MinSignedPerWindow(ctx)

	// if we are past the minimum height and the validator has missed too many blocks, punish them
	if height > minHeight && signInfo.MissedBlocksCounter > maxMissed {
		validator, found := k.validatorKeeper.GetValidatorByConsAddr(ctx, consAddr)
		if !found {
			return
		}

		// Downtime confirmed: slash the validator
		logger.Info(fmt.Sprintf("Validator %s past min height of %d and below signed blocks threshold of %d",
			consAddr, minHeight, k.MinSignedPerWindow(ctx)))

		// We need to reset the counter & array so that the validator won't be immediately slashed for downtime upon rebonding.
		signInfo.MissedBlocksCounter = 0
		signInfo.IndexOffset = 0
		k.ClearValidatorMissedBlockBitArray(ctx, consAddr)
		k.Slash(ctx, validator, power, k.SlashFractionDowntime(ctx), types.AttributeValueMissingSignature)
	}

	k.SetValidatorSigningInfo(ctx, signInfo)
}

// Slash a validator for an infraction
// Find the contributing stake and burn the specified slashFactor of it
func (k Keeper) Slash(ctx sdk.Context, validator staking.Validator, power int64, slashFactor sdk.Dec, reason string) {
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

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSlash,
			sdk.NewAttribute(types.AttributeKeyReason, reason),
		),
	)
}

// Gets the entire Penalty metadata for a nounce
func (k Keeper) GetPenalty(ctx sdk.Context, nounce uint64) (Penalty, bool) {
	store := ctx.KVStore(k.storeKey)

	if !store.Has(GetPenaltyKey(nounce)) {
		return Penalty{}, false
	}

	value := store.Get(GetPenaltyKey(nounce))
	var request Penalty
	k.cdc.MustUnmarshalBinaryBare(value, &request)
	return request, true
}

// Sets the entire Penalty metadata for a nounce
func (k Keeper) SetPenalty(ctx sdk.Context, nounce uint64, request Penalty) {
	store := ctx.KVStore(k.storeKey)
	store.Set(GetPenaltyKey(nounce), k.cdc.MustMarshalBinaryBare(request))
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
