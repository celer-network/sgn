package global

// nolint

import (
	"github.com/celer-network/sgn/x/global/keeper"
	"github.com/celer-network/sgn/x/global/types"
)

const (
	ModuleName        = types.ModuleName
	StoreKey          = types.StoreKey
	RouterKey         = types.RouterKey
	QuerierRoute      = types.QuerierRoute
	DefaultParamspace = types.DefaultParamspace
	QueryParams       = types.QueryParams
)

var (
	// functions aliases
	NewKeeper               = keeper.NewKeeper
	NewQuerier              = keeper.NewQuerier
	RegisterCodec           = types.RegisterCodec
	RegisterChangeTypeCodec = types.RegisterChangeTypeCodec
	NewGenesisState         = types.NewGenesisState
	DefaultGenesisState     = types.DefaultGenesisState
	ValidateGenesis         = types.ValidateGenesis
	NewParams               = types.NewParams

	// variable aliases
	ModuleCdc    = types.ModuleCdc
	EthBlkNumKey = types.EthBlkNumKey
)

type (
	Keeper       = keeper.Keeper
	GenesisState = types.GenesisState
	Params       = types.Params
)
