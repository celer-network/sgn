package guard

import (
	"github.com/celer-network/sgn/x/guard/client/cli"
	"github.com/celer-network/sgn/x/guard/types"
)

const (
	ModuleName        = types.ModuleName
	RouterKey         = types.RouterKey
	StoreKey          = types.StoreKey
	QuerySubscription = types.QuerySubscription
	QueryRequest      = types.QueryRequest
	QueryEpoch        = types.QueryEpoch
	QueryParameters   = types.QueryParameters

	ChanStatus_Idle        = types.ChanStatus_Idle
	ChanStatus_Withdrawing = types.ChanStatus_Withdrawing
	ChanStatus_Settling    = types.ChanStatus_Settling
	ChanStatus_Settled     = types.ChanStatus_Settled
)

var (
	ModuleCdc                  = types.ModuleCdc
	RegisterCodec              = types.RegisterCodec
	NewMsgRequestGuard         = types.NewMsgRequestGuard
	NewSubscription            = types.NewSubscription
	NewRequest                 = types.NewRequest
	NewInitRequest             = types.NewInitRequest
	NewGuardTrigger            = types.NewGuardTrigger
	NewGuardProof              = types.NewGuardProof
	NewQuerySubscriptionParams = types.NewQuerySubscriptionParams
	NewQueryRequestParams      = types.NewQueryRequestParams
	NewQueryEpochParams        = types.NewQueryEpochParams
	GetSubscriptionKey         = types.GetSubscriptionKey
	GetRequestKey              = types.GetRequestKey
	SubscriptionKeyPrefix      = types.SubscriptionKeyPrefix
	CLIQuerySubscription       = cli.QuerySubscription
	CLIQueryRequest            = cli.QueryRequest
	CLIQueryParams             = cli.QueryParams
	DefaultParams              = types.DefaultParams
)

type (
	MsgRequestGuard         = types.MsgRequestGuard
	Subscription            = types.Subscription
	InitRequest             = types.InitRequest
	GuardTrigger            = types.GuardTrigger
	GuardProof              = types.GuardProof
	Request                 = types.Request
	ChanStatus              = types.ChanStatus
	Params                  = types.Params
	QuerySubscriptionParams = types.QuerySubscriptionParams
	QueryRequestParams      = types.QueryRequestParams
	QueryEpochParams        = types.QueryEpochParams
)
