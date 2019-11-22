package types

import (
	"github.com/celer-network/sgn/mainchain"
)

const (
	QuerySubscription = "subscription"
	QueryRequest      = "request"
	QueryParameters   = "parameters"
)

type QuerySubscriptionParams struct {
	EthAddress string
}

func NewQuerySubscriptionParams(addr string) QuerySubscriptionParams {
	return QuerySubscriptionParams{
		EthAddress: mainchain.Hex2AddrHex(addr),
	}
}

type QueryRequestParams struct {
	ChannelId []byte
}

func NewQueryRequestParams(channelId []byte) QueryRequestParams {
	return QueryRequestParams{
		ChannelId: channelId,
	}
}
