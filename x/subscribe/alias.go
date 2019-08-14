package subscribe

import (
	"github.com/celer-network/sgn/x/subscribe/types"
)

const (
	ModuleName       = types.ModuleName
	RouterKey        = types.RouterKey
	StoreKey         = types.StoreKey
	QuerySubscrption = types.QuerySubscrption
)

var (
	NewMsgSubscribe = types.NewMsgSubscribe
	ModuleCdc       = types.ModuleCdc
	RegisterCodec   = types.RegisterCodec
	NewSubscription = types.NewSubscription
)

type (
	MsgSubscribe           = types.MsgSubscribe
	Subscription           = types.Subscription
	QuerySubscrptionParams = types.QuerySubscrptionParams
)
