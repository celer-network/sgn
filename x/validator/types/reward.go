package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Reward struct {
	Amount sdk.Int `json:"amount"`
}

func NewReward() Reward {
	return Reward{}
}

// implement fmt.Stringer
func (w Reward) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Amount: %v`, w.Amount))
}
