package types

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
)

const (
	QueryPenalty    = "penalty"
	QueryParameters = "parameters"
)

type QueryPenaltyParams struct {
	EthAddress string
}

func NewQueryPenaltyParams(addr string) QueryPenaltyParams {
	return QueryPenaltyParams{
		EthAddress: ethcommon.HexToAddress(addr).String(),
	}
}
