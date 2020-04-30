package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Subscription struct {
	EthAddress string  `json:"ethAddress"`
	Deposit    sdk.Int `json:"deposit"`
	Spend      sdk.Int `json:"spend"`
}

func NewSubscription(ethAddress string) Subscription {
	return Subscription{
		EthAddress: ethAddress,
		Deposit:    sdk.ZeroInt(),
		Spend:      sdk.ZeroInt(),
	}
}

// implement fmt.Stringer
func (s Subscription) String() string {
	return strings.TrimSpace(fmt.Sprintf(`EthAddress: %s, Deposit: %v, Spend: %v`, s.EthAddress, s.Deposit, s.Spend))
}
