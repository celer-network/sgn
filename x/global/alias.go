package global

import (
	"github.com/celer-network/sgn/x/global/client/cli"
	"github.com/celer-network/sgn/x/global/types"
)

const (
	ModuleName       = types.ModuleName
	RouterKey        = types.RouterKey
	StoreKey         = types.StoreKey
	QueryLatestBlock = types.QueryLatestBlock
	QueryEpoch       = types.QueryEpoch
	QueryParameters  = types.QueryParameters
)

var (
	ModuleCdc           = types.ModuleCdc
	RegisterCodec       = types.RegisterCodec
	NewBlock            = types.NewBlock
	NewEpoch            = types.NewEpoch
	NewMsgSyncBlock     = types.NewMsgSyncBlock
	NewQueryEpochParams = types.NewQueryEpochParams
	LatestBlockKey      = types.LatestBlockKey
	GetEpochKey         = types.GetEpochKey
	GetLatestEpochKey   = types.GetLatestEpochKey
	CLIQueryLatestBlock = cli.QueryLatestBlock
	DefaultParams       = types.DefaultParams
)

type (
	Block            = types.Block
	Epoch            = types.Epoch
	Params           = types.Params
	MsgSyncBlock     = types.MsgSyncBlock
	QueryEpochParams = types.QueryEpochParams
)
