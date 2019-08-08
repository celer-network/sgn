package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	QueryEthAddress = "ethAddress"
)

type QueryEthAddressParams struct {
	Address sdk.AccAddress
}

func NewQueryEthAddressParams(addr sdk.AccAddress) QueryEthAddressParams {
	return QueryEthAddressParams{
		Address: addr,
	}
}
