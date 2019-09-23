package types

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// subscribe params default values
const (
	// Default request handler count
	DefaultRequestHandlerCount uint64 = 3
)

// nolint - Keys for parameter access
var (
	KeyRequestHandlerCount = []byte("RequestHandlerCount")
)

var _ params.ParamSet = (*Params)(nil)

// Params defines the high level settings for subscribe
type Params struct {
	RequestHandlerCount uint64 `json:"requestHandlerCount" yaml:"requestHandlerCount"` // epoch length based on seconds
}

// NewParams creates a new Params instance
func NewParams(requestHandlerCount uint64) Params {

	return Params{
		RequestHandlerCount: requestHandlerCount,
	}
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{KeyRequestHandlerCount, &p.RequestHandlerCount},
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
	return NewParams(DefaultRequestHandlerCount)
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	return fmt.Sprintf(`Params:
  RequestHandlerCount:    %d`,
		p.RequestHandlerCount)
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
	if p.RequestHandlerCount == 0 {
		return fmt.Errorf("subscribe parameter RequestHandlerCount must be a positive integer")
	}
	return nil
}
