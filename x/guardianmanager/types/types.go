package types

import (
	"fmt"
	"strings"
)

type Guardian struct {
	Balance uint64 `json:"balance"`
}

// Returns a new Number with the minprice as the price
func NewGuardian() Guardian {
	return Guardian{}
}

// implement fmt.Stringer
func (g Guardian) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Balance: %d`, g.Balance))
}
