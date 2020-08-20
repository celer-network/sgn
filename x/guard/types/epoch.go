package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Epoch struct {
	Id        sdk.Int `json:"id"`
	Timestamp int64   `json:"timestamp"`
	TotalFee  sdk.Int `json:"total_fee"`
}

func NewEpoch(id sdk.Int, timestamp int64) Epoch {
	return Epoch{
		Id:        id,
		Timestamp: timestamp,
		TotalFee:  sdk.ZeroInt(),
	}
}

// implement fmt.Stringer
func (e Epoch) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Id: %v, Timestamp: %d, TotalFee: %v`, e.Id, e.Timestamp, e.TotalFee))
}
