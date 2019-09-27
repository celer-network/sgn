package validator

import (
	"github.com/celer-network/sgn/x/validator/client/cli"
	"github.com/celer-network/sgn/x/validator/types"
)

const (
	ModuleName             = types.ModuleName
	RouterKey              = types.RouterKey
	StoreKey               = types.StoreKey
	QueryPuller            = types.QueryPuller
	QueryPusher            = types.QueryPusher
	QueryDelegator         = types.QueryDelegator
	QueryCandidate         = types.QueryCandidate
	QueryReward            = types.QueryReward
	TypeMsgWithdrawReward  = types.TypeMsgWithdrawReward
	AttributeKeyEthAddress = types.AttributeKeyEthAddress
	MsgWithdrawRewardName  = types.MsgWithdrawRewardName
	ActionInitiateWithdraw = types.ActionInitiateWithdraw
)

var (
	NewMsgInitializeCandidate = types.NewMsgInitializeCandidate
	NewMsgClaimValidator      = types.NewMsgClaimValidator
	NewMsgSyncValidator       = types.NewMsgSyncValidator
	NewMsgSyncDelegator       = types.NewMsgSyncDelegator
	NewMsgWithdrawReward      = types.NewMsgWithdrawReward
	NewMsgSignReward          = types.NewMsgSignReward
	NewQueryRewardParams      = types.NewQueryRewardParams
	ModuleCdc                 = types.ModuleCdc
	RegisterCodec             = types.RegisterCodec
	PullerKey                 = types.PullerKey
	PusherKey                 = types.PusherKey
	GetDelegatorKey           = types.GetDelegatorKey
	GetDelegatorsKey          = types.GetDelegatorsKey
	GetCandidateKey           = types.GetCandidateKey
	GetRewardKey              = types.GetRewardKey
	NewPuller                 = types.NewPuller
	NewPusher                 = types.NewPusher
	NewDelegator              = types.NewDelegator
	NewCandidate              = types.NewCandidate
	NewReward                 = types.NewReward
	CLIQueryPuller            = cli.QueryPuller
	CLIQueryPusher            = cli.QueryPusher
	CLIQueryReward            = cli.QueryReward
)

type (
	Puller                 = types.Puller
	Pusher                 = types.Pusher
	Delegator              = types.Delegator
	Candidate              = types.Candidate
	Reward                 = types.Reward
	QueryDelegatorParams   = types.QueryDelegatorParams
	QueryCandidateParams   = types.QueryCandidateParams
	QueryRewardParams      = types.QueryRewardParams
	MsgInitializeCandidate = types.MsgInitializeCandidate
	MsgClaimValidator      = types.MsgClaimValidator
	MsgSyncValidator       = types.MsgSyncValidator
	MsgSyncDelegator       = types.MsgSyncDelegator
	MsgWithdrawReward      = types.MsgWithdrawReward
	MsgSignReward          = types.MsgSignReward
)
