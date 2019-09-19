package types

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
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
		EthAddress: ethcommon.HexToAddress(addr).String(),
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

type QueryEpochParams struct {
	EpochId int64
}

func NewQueryEpochParams(epochId int64) QueryEpochParams {
	return QueryEpochParams{
		EpochId: epochId,
	}
}
