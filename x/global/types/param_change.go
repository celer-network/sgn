package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ParamChange struct {
	Record   sdk.Int `json:"record"`
	NewValue sdk.Int `json:"newValue"`
}

func NewParamChange(record, newValue sdk.Int) ParamChange {
	return ParamChange{
		Record:   record,
		NewValue: newValue,
	}
}

// implement fmt.Stringer
func (p ParamChange) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Record: %v, NewValue: %v`, p.Record, p.NewValue))
}
