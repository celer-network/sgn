package keeper

import (
	"github.com/celer-network/sgn/x/sync/types"
	"github.com/celer-network/sgn/x/validator"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper defines the sync module Keeper
type Keeper struct {
	// The reference to the Paramstore to get and set sync specific params
	paramSpace types.ParamSubspace

	vk validator.Keeper

	// The (unexposed) keys used to access the stores from the Context.
	storeKey sdk.StoreKey

	// The codec codec for binary encoding/decoding.
	cdc *codec.Codec
}

// NewKeeper returns a sync keeper. It handles:
// - submitting sync changes
// - depositing funds into changes, and activating upon sufficient funds being deposited
// - users voting on changes, with weight proportional to stake in the system
// - and tallying the result of the vote.
//
// CONTRACT: the parameter Subspace must have the param key table already initialized
func NewKeeper(
	cdc *codec.Codec, key sdk.StoreKey, paramSpace types.ParamSubspace,
	vk validator.Keeper,
) Keeper {

	// It is vital to seal the sync change router here as to not allow
	// further handlers to be registered after the keeper is created since this
	// could create invalid or non-deterministic behavior.
	rtr.Seal()

	return Keeper{
		storeKey:   key,
		paramSpace: paramSpace,
		vk:         vk,
		cdc:        cdc,
	}
}

// Router returns the sync Keeper's Router
func (keeper Keeper) Router() types.Router {
	return keeper.router
}

func (keeper Keeper) GetValidators(ctx sdk.Context) []staking.Validator {
	return keeper.vk.GetValidators(ctx)
}
