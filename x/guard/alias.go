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
)

var (
	ModuleCdc                  = types.ModuleCdc
	RegisterCodec              = types.RegisterCodec
	NewMsgRequestGuard         = types.NewMsgRequestGuard
	NewSubscription            = types.NewSubscription
	NewRequest                 = types.NewRequest
	NewEpoch                   = types.NewEpoch
	NewQuerySubscriptionParams = types.NewQuerySubscriptionParams
	NewQueryRequestParams      = types.NewQueryRequestParams
	NewQueryEpochParams        = types.NewQueryEpochParams
	GetSubscriptionKey         = types.GetSubscriptionKey
	GetRequestKey              = types.GetRequestKey
	GetEpochKey                = types.GetEpochKey
	GetLatestEpochKey          = types.GetLatestEpochKey
	SubscriptionKeyPrefix      = types.SubscriptionKeyPrefix
	CLIQuerySubscription       = cli.QuerySubscription
	CLIQueryRequest            = cli.QueryRequest
	CLIQueryParams             = cli.QueryParams
	DefaultParams              = types.DefaultParams
)

type (
	MsgRequestGuard         = types.MsgRequestGuard
	Subscription            = types.Subscription
	Request                 = types.Request
	Epoch                   = types.Epoch
	Params                  = types.Params
	QuerySubscriptionParams = types.QuerySubscriptionParams
	QueryRequestParams      = types.QueryRequestParams
	QueryEpochParams        = types.QueryEpochParams
)
