package types

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
)

const (
	QueryGuardian = "guardian"
	QueryRequest  = "request"
)

type QueryGuardianParams struct {
	EthAddress string
}

func NewQueryGuardianParams(addr string) QueryGuardianParams {
	return QueryGuardianParams{
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
