package slash

import (
	"github.com/celer-network/sgn/x/slash/types"
)

const (
	ModuleName      = types.ModuleName
	RouterKey       = types.RouterKey
	StoreKey        = types.StoreKey
	QueryPenalty    = types.QueryPenalty
	QueryParameters = types.QueryParameters
)

var (
	ModuleCdc             = types.ModuleCdc
	RegisterCodec         = types.RegisterCodec
	NewPenalty            = types.NewPenalty
	NewQueryPenaltyParams = types.NewQueryPenaltyParams
	GetPenaltyKey         = types.GetPenaltyKey
	DefaultParams         = types.DefaultParams
)

type (
	Penalty            = types.Penalty
	Params             = types.Params
	QueryPenaltyParams = types.QueryPenaltyParams
)
