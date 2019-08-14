package types

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
)

const (
	QueryGuardian = "guardian"
)

type QueryGuardianParams struct {
	EthAddress string
}

func NewQueryGuardianParams(addr string) QueryGuardianParams {
	return QueryGuardianParams{
		EthAddress: ethcommon.HexToAddress(addr).String(),
	}
}
