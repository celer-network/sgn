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
	QueryPusher = types.QueryPusher
)

var (
	NewMsgClaimValidator = types.NewMsgClaimValidator
	NewMsgSyncValidator  = types.NewMsgSyncValidator
	ModuleCdc            = types.ModuleCdc
	RegisterCodec        = types.RegisterCodec
	PullerKey            = types.PullerKey
	PusherKey            = types.PusherKey
	NewPuller            = types.NewPuller
	NewPusher            = types.NewPusher
	CLIQueryPuller       = cli.QueryPuller
	CLIQueryPusher       = cli.QueryPusher
)

type (
	Puller            = types.Puller
	Pusher            = types.Pusher
	MsgClaimValidator = types.MsgClaimValidator
	MsgSyncValidator  = types.MsgSyncValidator
)
