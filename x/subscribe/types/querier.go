package types

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
)

const (
	QuerySubscrption = "subscription"
	QueryRequest     = "request"
)

type QuerySubscrptionParams struct {
	EthAddress string
}

func NewQuerySubscrptionParams(addr string) QuerySubscrptionParams {
	return QuerySubscrptionParams{
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
