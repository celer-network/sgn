package types

import (
	"fmt"
	"strings"
)

type Subscription struct {
	Expiration uint64 `json:"expiration"`
}

// Returns a new Number with the minprice as the price
func NewSubscription(expiration uint64) Subscription {
	return Subscription{
		Expiration: expiration,
	}
}

// implement fmt.Stringer
func (s Subscription) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Expiration: %d`, s.Expiration))
}
