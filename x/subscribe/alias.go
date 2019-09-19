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
	NewEpoch                   = types.NewEpoch
	NewQuerySubscriptionParams = types.NewQuerySubscriptionParams
	NewQueryRequestParams      = types.NewQueryRequestParams
	GetSubscriptionKey         = types.GetSubscriptionKey
	GetRequestKey              = types.GetRequestKey
	GetEpochKey                = types.GetEpochKey
	GetLatestEpochKey          = types.GetLatestEpochKey
	SubscriptionKey            = types.SubscriptionKey
	CLIQueryRequest            = cli.QueryRequest
)

type (
	MsgSubscribe            = types.MsgSubscribe
	MsgRequestGuard         = types.MsgRequestGuard
	MsgGuardProof           = types.MsgGuardProof
	Subscription            = types.Subscription
	Request                 = types.Request
	Epoch                   = types.Epoch
	QuerySubscriptionParams = types.QuerySubscriptionParams
	QueryRequestParams      = types.QueryRequestParams
)
