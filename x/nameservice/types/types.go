package types

import (
	"fmt"
	"strings"
)

type Number struct {
	Value uint `json:"value"`
}

// Returns a new Number with the minprice as the price
func NewNumber() Number {
	return Number{
		Value: 2,
	}
}

// implement fmt.Stringer
func (w Number) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Value: %s`, w.Value))
}
