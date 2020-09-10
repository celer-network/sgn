package validator

import (
	"github.com/celer-network/sgn/x/validator/client/cli"
	"github.com/celer-network/sgn/x/validator/types"
)

const (
	ModuleName               = types.ModuleName
	RouterKey                = types.RouterKey
	StoreKey                 = types.StoreKey
	QuerySyncer              = types.QuerySyncer
	QueryDelegator           = types.QueryDelegator
	QueryCandidate           = types.QueryCandidate
	QueryCandidateDelegators = types.QueryCandidateDelegators
	QueryReward              = types.QueryReward
	QueryRewardEpoch         = types.QueryRewardEpoch
	QueryParameters          = types.QueryParameters
	TypeMsgWithdrawReward    = types.TypeMsgWithdrawReward
	AttributeKeyEthAddress   = types.AttributeKeyEthAddress
	ActionInitiateWithdraw   = types.ActionInitiateWithdraw
)

var (
	NewMsgSetTransactors           = types.NewMsgSetTransactors
	NewMsgWithdrawReward           = types.NewMsgWithdrawReward
	NewMsgEditCandidateDescription = types.NewMsgEditCandidateDescription
	NewMsgSignReward               = types.NewMsgSignReward
	NewQueryRewardParams           = types.NewQueryRewardParams
	ModuleCdc                      = types.ModuleCdc
	RegisterCodec                  = types.RegisterCodec
	SyncerKey                      = types.SyncerKey
	CandidateKeyPrefix             = types.CandidateKeyPrefix
	GetDelegatorKey                = types.GetDelegatorKey
	GetDelegatorsKey               = types.GetDelegatorsKey
	GetCandidateKey                = types.GetCandidateKey
	GetRewardKey                   = types.GetRewardKey
	GetPendingRewardKey            = types.GetPendingRewardKey
	RewardEpochKey                 = types.RewardEpochKey
	NewSyncer                      = types.NewSyncer
	NewDelegator                   = types.NewDelegator
	NewCandidate                   = types.NewCandidate
	NewReward                      = types.NewReward
	NewRewardEpoch                 = types.NewRewardEpoch
	NewPendingReward               = types.NewPendingReward
	CLIQuerySyncer                 = cli.QuerySyncer
	CLIQueryCandidate              = cli.QueryCandidate
	CLIQueryCandidateDelegators    = cli.QueryCandidateDelegators
	CLIQueryReward                 = cli.QueryReward
	CLIQueryDelegator              = cli.QueryDelegator
	CLIQueryValidator              = cli.QueryValidator
	CLIQueryValidators             = cli.QueryValidators
	CLIQueryBondedValidators       = cli.QueryBondedValidators
	CLIQueryParams                 = cli.QueryParams
	DefaultParams                  = types.DefaultParams
)

type (
	Syncer                      = types.Syncer
	Params                      = types.Params
	Delegator                   = types.Delegator
	Candidate                   = types.Candidate
	Reward                      = types.Reward
	RewardEpoch                 = types.RewardEpoch
	PendingReward               = types.PendingReward
	QueryDelegatorParams        = types.QueryDelegatorParams
	QueryCandidateParams        = types.QueryCandidateParams
	QueryRewardParams           = types.QueryRewardParams
	MsgSetTransactors           = types.MsgSetTransactors
	MsgEditCandidateDescription = types.MsgEditCandidateDescription
	MsgWithdrawReward           = types.MsgWithdrawReward
	MsgSignReward               = types.MsgSignReward
)
