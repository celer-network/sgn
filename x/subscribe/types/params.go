package types

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// subscribe params default values
const (
	// Default request guard count
	DefaultRequestGuardCount uint64 = 3
)

// nolint - Keys for parameter access
var (
	KeyRequestGuardCount = []byte("RequestGuardCount")
)

var _ params.ParamSet = (*Params)(nil)

// Params defines the high level settings for subscribe
type Params struct {
	RequestGuardCount uint64 `json:"requestGuardCount" yaml:"requestGuardCount"` // epoch length based on seconds
}

// NewParams creates a new Params instance
func NewParams(requestGuardCount uint64) Params {

	return Params{
		RequestGuardCount: requestGuardCount,
	}
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{KeyRequestGuardCount, &p.RequestGuardCount},
	}
}

// Equal returns a boolean determining if two Param types are identical.
// TODO: This is slower than comparing struct fields directly
func (p Params) Equal(p2 Params) bool {
	bz1 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p)
	bz2 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(DefaultRequestGuardCount)
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	return fmt.Sprintf(`Params:
  RequestGuardCount:    %d`,
		p.RequestGuardCount)
}

// unmarshal the current subscribe params value from store key or panic
func MustUnmarshalParams(cdc *codec.Codec, value []byte) Params {
	params, err := UnmarshalParams(cdc, value)
	if err != nil {
		panic(err)
	}
	return params
}

// unmarshal the current subscribe params value from store key
func UnmarshalParams(cdc *codec.Codec, value []byte) (params Params, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &params)
	if err != nil {
		return
	}
	return
}

// validate a set of params
func (p Params) Validate() error {
	if p.RequestGuardCount == 0 {
		return fmt.Errorf("subscribe parameter RequestGuardCount must be a positive integer")
	}
	return nil
}
