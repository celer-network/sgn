package types

import (
	"fmt"
	"strings"

	"github.com/celer-network/sgn/mainchain"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Subscription struct {
	EthAddress string  `json:"eth_address"`
	Deposit    sdk.Int `json:"deposit"`
	Spend      sdk.Int `json:"spend"`
}

func NewSubscription(ethAddress string) Subscription {
	return Subscription{
		EthAddress: mainchain.FormatAddrHex(ethAddress),
		Deposit:    sdk.ZeroInt(),
		Spend:      sdk.ZeroInt(),
	}
}

// implement fmt.Stringer
func (s Subscription) String() string {
	return strings.TrimSpace(fmt.Sprintf(`EthAddress: %s, Deposit: %v, Spend: %v`, s.EthAddress, s.Deposit, s.Spend))
}
