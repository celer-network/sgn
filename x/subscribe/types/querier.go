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
		EthAddress: mainchain.FormatAddrHex(addr),
	}
}

type QueryRequestParams struct {
	ChannelId []byte
	PeerFrom  string
}

func NewQueryRequestParams(channelId []byte, peerFrom string) QueryRequestParams {
	return QueryRequestParams{
		ChannelId: channelId,
		PeerFrom:  mainchain.FormatAddrHex(peerFrom),
	}
}
