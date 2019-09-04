package validator

import (
	"github.com/celer-network/sgn/x/validator/client/cli"
	"github.com/celer-network/sgn/x/validator/types"
)

const (
	ModuleName  = types.ModuleName
	RouterKey   = types.RouterKey
	StoreKey    = types.StoreKey
	QueryPuller = types.QueryPuller
)

var (
	NewMsgClaimValidator = types.NewMsgClaimValidator
	ModuleCdc            = types.ModuleCdc
	RegisterCodec        = types.RegisterCodec
	PullerKey            = types.PullerKey
	NewPuller            = types.NewPuller
	CLIQueryPuller       = cli.QueryPuller
)

type (
	Puller            = types.Puller
	MsgClaimValidator = types.MsgClaimValidator
)
