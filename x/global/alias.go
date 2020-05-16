package global

import (
	"github.com/celer-network/sgn/x/global/client/cli"
	"github.com/celer-network/sgn/x/global/types"
)

const (
	ModuleName      = types.ModuleName
	RouterKey       = types.RouterKey
	StoreKey        = types.StoreKey
	QueryEpoch      = types.QueryEpoch
	QueryParameters = types.QueryParameters
)

var (
	ModuleCdc           = types.ModuleCdc
	RegisterCodec       = types.RegisterCodec
	NewBlock            = types.NewBlock
	NewEpoch            = types.NewEpoch
	NewQueryEpochParams = types.NewQueryEpochParams
	GetEpochKey         = types.GetEpochKey
	GetLatestEpochKey   = types.GetLatestEpochKey
	DefaultParams       = types.DefaultParams
	CLIQueryParams      = cli.QueryParams
)

type (
	Block            = types.Block
	Epoch            = types.Epoch
	Params           = types.Params
	QueryEpochParams = types.QueryEpochParams
)
