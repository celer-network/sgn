package global

import (
	"github.com/celer-network/sgn/x/global/types"
)

const (
	ModuleName       = types.ModuleName
	RouterKey        = types.RouterKey
	StoreKey         = types.StoreKey
	QueryLatestBlock = types.QueryLatestBlock
)

var (
	ModuleCdc       = types.ModuleCdc
	RegisterCodec   = types.RegisterCodec
	NewBlock        = types.NewBlock
	NewMsgSyncBlock = types.NewMsgSyncBlock
	LatestBlockKey  = types.LatestBlockKey
)

type (
	Block        = types.Block
	MsgSyncBlock = types.MsgSyncBlock
)
