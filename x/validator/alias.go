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
	QueryDelegator         = types.QueryDelegator
	QueryCandidate         = types.QueryCandidate
	QueryReward            = types.QueryReward
	TypeMsgWithdrawReward  = types.TypeMsgWithdrawReward
	AttributeKeyEthAddress = types.AttributeKeyEthAddress
	ActionInitiateWithdraw = types.ActionInitiateWithdraw
	ServiceReward          = types.ServiceReward
	MiningReward           = types.MiningReward
)

var (
	NewMsgSetTransactors           = types.NewMsgSetTransactors
	NewMsgWithdrawReward           = types.NewMsgWithdrawReward
	NewMsgEditCandidateDescription = types.NewMsgEditCandidateDescription
	NewMsgSignReward               = types.NewMsgSignReward
	NewQueryRewardParams           = types.NewQueryRewardParams
	ModuleCdc                      = types.ModuleCdc
	RegisterCodec                  = types.RegisterCodec
	PullerKey                      = types.PullerKey
	CandidateKeyPrefix             = types.CandidateKeyPrefix
	GetDelegatorKey                = types.GetDelegatorKey
	GetDelegatorsKey               = types.GetDelegatorsKey
	GetCandidateKey                = types.GetCandidateKey
	GetRewardKey                   = types.GetRewardKey
	NewPuller                      = types.NewPuller
	NewDelegator                   = types.NewDelegator
	NewCandidate                   = types.NewCandidate
	NewReward                      = types.NewReward
	CLIQueryPuller                 = cli.QueryPuller
	CLIQueryCandidate              = cli.QueryCandidate
	CLIQueryReward                 = cli.QueryReward
	CLIQueryDelegator              = cli.QueryDelegator
	CLIQueryValidator              = cli.QueryValidator
	CLIQueryValidators             = cli.QueryValidators
	CLIQueryBondedValidators       = cli.QueryBondedValidators
	DefaultParams                  = types.DefaultParams
)

type (
	Puller                      = types.Puller
	Params                      = types.Params
	Delegator                   = types.Delegator
	Candidate                   = types.Candidate
	Reward                      = types.Reward
	RewardType                  = types.RewardType
	QueryDelegatorParams        = types.QueryDelegatorParams
	QueryCandidateParams        = types.QueryCandidateParams
	QueryRewardParams           = types.QueryRewardParams
	MsgSetTransactors           = types.MsgSetTransactors
	MsgEditCandidateDescription = types.MsgEditCandidateDescription
	MsgWithdrawReward           = types.MsgWithdrawReward
	MsgSignReward               = types.MsgSignReward
)
