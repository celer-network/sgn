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
	NewMsgInitializeCandidate = types.NewMsgInitializeCandidate
	NewMsgClaimValidator      = types.NewMsgClaimValidator
	NewMsgSyncValidator       = types.NewMsgSyncValidator
	NewMsgSyncDelegator       = types.NewMsgSyncDelegator
	ModuleCdc                 = types.ModuleCdc
	RegisterCodec             = types.RegisterCodec
	PullerKey                 = types.PullerKey
	PusherKey                 = types.PusherKey
	GetDelegatorKey           = types.GetDelegatorKey
	GetDelegatorsKey          = types.GetDelegatorsKey
	GetCandidateKey     = types.GetCandidateKey
	GetCandidateSnapshotKey           = types.GetCandidateSnapshotKey
	NewPuller                 = types.NewPuller
	NewPusher                 = types.NewPusher
	NewDelegator              = types.NewDelegator
	NewCandidate              = types.NewCandidate
	CLIQueryPuller            = cli.QueryPuller
	CLIQueryPusher            = cli.QueryPusher
)

type (
	Puller                 = types.Puller
	Pusher                 = types.Pusher
	Delegator              = types.Delegator
	Candidate              = types.Candidate
	MsgInitializeCandidate = types.MsgInitializeCandidate
	MsgClaimValidator      = types.MsgClaimValidator
	MsgSyncValidator       = types.MsgSyncValidator
	MsgSyncDelegator       = types.MsgSyncDelegator
)
