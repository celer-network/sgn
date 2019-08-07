package types

import (
	"fmt"
	"strings"
)

type EthAddress struct {
	Address string `json:"address"`
}

// Returns a new Number with the minprice as the price
func NewEthAddress(ethAddress string) EthAddress {
	return EthAddress{
		Address: ethAddress,
	}
}

// implement fmt.Stringer
func (e EthAddress) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Address: %s`, e.Address))
}
