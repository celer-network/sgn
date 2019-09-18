package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Epoch struct {
	Id          sdk.Int `json:"id"`
	BlockNumber uint64  `json:"blockNumber"`
	TotalFee    sdk.Int `json:"totalFee"`
}

// Returns a new Number with the minprice as the price
func NewEpoch(id sdk.Int, blockNumber uint64) Epoch {
	return Epoch{
		Id:          id,
		BlockNumber: blockNumber,
		TotalFee:    sdk.NewInt(0),
	}
}

// implement fmt.Stringer
func (e Epoch) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Id: %v, BlockNumber: %d, TotalFee: %v`, e.Id, e.BlockNumber, e.TotalFee))
}
