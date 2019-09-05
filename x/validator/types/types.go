package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Puller struct {
	ValidatorIdx  uint           `json:"validatorIdx"`
	ValidatorAddr sdk.AccAddress `json:"validatorAddr"`
}

// Returns a new Number with the minprice as the price
func NewPuller(validatorIdx uint, validatorAddr sdk.AccAddress) Puller {
	return Puller{
		ValidatorIdx:  validatorIdx,
		ValidatorAddr: validatorAddr,
	}
}

// implement fmt.Stringer
func (r Puller) String() string {
	return strings.TrimSpace(fmt.Sprintf(`ValidatorIdx: %d, ValidatorAddr: %x`, r.ValidatorIdx, r.ValidatorAddr))
}

type Pusher struct {
	ValidatorIdx  uint           `json:"validatorIdx"`
	ValidatorAddr sdk.AccAddress `json:"validatorAddr"`
}

// Returns a new Number with the minprice as the price
func NewPusher(validatorIdx uint, validatorAddr sdk.AccAddress) Pusher {
	return Pusher{
		ValidatorIdx:  validatorIdx,
		ValidatorAddr: validatorAddr,
	}
}

// implement fmt.Stringer
func (r Pusher) String() string {
	return strings.TrimSpace(fmt.Sprintf(`ValidatorIdx: %d, ValidatorAddr: %x`, r.ValidatorIdx, r.ValidatorAddr))
}
