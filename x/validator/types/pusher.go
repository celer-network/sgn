package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Pusher struct {
	ValidatorIdx  uint           `json:"validatorIdx"`
	ValidatorAddr sdk.AccAddress `json:"validatorAddr"`
}

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
