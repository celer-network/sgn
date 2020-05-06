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
	ActionInitiateWithdraw = types.ActionInitiateWithdraw
)

var (
	NewMsgSetTransactors     = types.NewMsgSetTransactors
	NewMsgClaimValidator     = types.NewMsgClaimValidator
	NewMsgSyncValidator      = types.NewMsgSyncValidator
	NewMsgWithdrawReward     = types.NewMsgWithdrawReward
	NewMsgSignReward         = types.NewMsgSignReward
	NewQueryRewardParams     = types.NewQueryRewardParams
	ModuleCdc                = types.ModuleCdc
	RegisterCodec            = types.RegisterCodec
	PullerKey                = types.PullerKey
	PusherKey                = types.PusherKey
	CandidateKeyPrefix       = types.CandidateKeyPrefix
	GetDelegatorKey          = types.GetDelegatorKey
	GetDelegatorsKey         = types.GetDelegatorsKey
	GetCandidateKey          = types.GetCandidateKey
	GetRewardKey             = types.GetRewardKey
	NewPuller                = types.NewPuller
	NewPusher                = types.NewPusher
	NewDelegator             = types.NewDelegator
	NewCandidate             = types.NewCandidate
	NewReward                = types.NewReward
	CLIQueryPuller           = cli.QueryPuller
	CLIQueryPusher           = cli.QueryPusher
	CLIQueryCandidate        = cli.QueryCandidate
	CLIQueryReward           = cli.QueryReward
	CLIQueryDelegator        = cli.QueryDelegator
	CLIQueryValidator        = cli.QueryValidator
	CLIQueryValidators       = cli.QueryValidators
	CLIQueryBondedValidators = cli.QueryBondedValidators
	DefaultParams            = types.DefaultParams
)

type (
	Puller               = types.Puller
	Pusher               = types.Pusher
	Params               = types.Params
	Delegator            = types.Delegator
	Candidate            = types.Candidate
	Reward               = types.Reward
	QueryDelegatorParams = types.QueryDelegatorParams
	QueryCandidateParams = types.QueryCandidateParams
	QueryRewardParams    = types.QueryRewardParams
	MsgSetTransactors    = types.MsgSetTransactors
	MsgClaimValidator    = types.MsgClaimValidator
	MsgSyncValidator     = types.MsgSyncValidator
	MsgWithdrawReward    = types.MsgWithdrawReward
	MsgSignReward        = types.MsgSignReward
)
