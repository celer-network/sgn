package global

import (
	"github.com/celer-network/sgn/x/global/types"
)

const (
	ModuleName = types.ModuleName
	RouterKey  = types.RouterKey
	StoreKey   = types.StoreKey
)

var (
	ModuleCdc      = types.ModuleCdc
	RegisterCodec  = types.RegisterCodec
	NewBlock       = types.NewBlock
	LatestBlockKey = types.LatestBlockKey
)

type (
	Block = types.Block
)
