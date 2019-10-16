package subscribe

import (
	"github.com/celer-network/sgn/x/subscribe/client/cli"
	"github.com/celer-network/sgn/x/subscribe/types"
)

const (
	ModuleName        = types.ModuleName
	RouterKey         = types.RouterKey
	StoreKey          = types.StoreKey
	QuerySubscription = types.QuerySubscription
	QueryRequest      = types.QueryRequest
	QueryParameters   = types.QueryParameters
)

var (
	NewMsgSubscribe            = types.NewMsgSubscribe
	NewMsgRequestGuard         = types.NewMsgRequestGuard
	NewMsgGuardProof           = types.NewMsgGuardProof
	ModuleCdc                  = types.ModuleCdc
	RegisterCodec              = types.RegisterCodec
	NewSubscription            = types.NewSubscription
	NewRequest                 = types.NewRequest
	NewQuerySubscriptionParams = types.NewQuerySubscriptionParams
	NewQueryRequestParams      = types.NewQueryRequestParams
	GetSubscriptionKey         = types.GetSubscriptionKey
	GetRequestKey              = types.GetRequestKey
	SubscriptionKeyPrefix      = types.SubscriptionKeyPrefix
	RequestGuardIdKey          = types.RequestGuardIdKey
	CLIQuerySubscription       = cli.QuerySubscription
	CLIQueryRequest            = cli.QueryRequest
	DefaultParams              = types.DefaultParams
)

type (
	MsgSubscribe            = types.MsgSubscribe
	MsgRequestGuard         = types.MsgRequestGuard
	MsgGuardProof           = types.MsgGuardProof
	Subscription            = types.Subscription
	Request                 = types.Request
	Params                  = types.Params
	QuerySubscriptionParams = types.QuerySubscriptionParams
	QueryRequestParams      = types.QueryRequestParams
)
