package types

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
)

const (
	QuerySubscrption = "subscription"
)

type QuerySubscrptionParams struct {
	EthAddress string
}

func NewQuerySubscrptionParams(addr string) QuerySubscrptionParams {
	return QuerySubscrptionParams{
		EthAddress: ethcommon.HexToAddress(addr).String(),
	}
}
