package types

import (
	"fmt"

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

// Query Result Payload for a names query
type QueryResEthAddress struct {
	Address string `json:"address"`
}

// implement fmt.Stringer
func (q QueryResEthAddress) String() string {
	return fmt.Sprint(q.Address)
}
