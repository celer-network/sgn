package types

import (
	"github.com/celer-network/sgn/mainchain"
)

const (
	QuerySubscription = "subscription"
	QueryRequest      = "request"
	QueryEpoch        = "epoch"
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
	Owner     string
}

func NewQueryRequestParams(channelId []byte, owner string) QueryRequestParams {
	return QueryRequestParams{
		ChannelId: channelId,
		Owner:     mainchain.FormatAddrHex(owner),
	}
}

type QueryEpochParams struct {
	EpochId int64
}

func NewQueryEpochParams(epochId int64) QueryEpochParams {
	return QueryEpochParams{
		EpochId: epochId,
	}
}
