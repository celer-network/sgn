package types

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// subscribe params default values
const (
	// Default epoch length based on seconds
	DefaultEpochLength int64 = 60

	// Default cost per epoch, 1 CELR token per epoch
	DefaultCostPerEpoch int64 = 1000000000000000000
)

// nolint - Keys for parameter access
var (
	KeyEpochLength  = []byte("EpochLength")
	KeyCostPerEpoch = []byte("KeyCostPerEpoch")
)

var _ params.ParamSet = (*Params)(nil)

// Params defines the high level settings for subscribe
type Params struct {
	EpochLength  int64   `json:"epochLength" yaml:"epochLength"`   // epoch length based on seconds
	CostPerEpoch sdk.Int `json:"costPerEpoch" yaml:"costPerEpoch"` // The fee will be charged for subscription per epoch
}

// NewParams creates a new Params instance
func NewParams(EpochLength int64, CostPerEpoch sdk.Int) Params {

	return Params{
		EpochLength:  EpochLength,
		CostPerEpoch: CostPerEpoch,
	}
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{KeyEpochLength, &p.EpochLength},
		{KeyCostPerEpoch, &p.CostPerEpoch},
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
	return NewParams(DefaultEpochLength, sdk.NewInt(DefaultCostPerEpoch))
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	return fmt.Sprintf(`Params:
  Max Validators:    %d
  Cost Per Epoch:       %d`,
		p.EpochLength, p.CostPerEpoch)
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
	if p.EpochLength == 0 {
		return fmt.Errorf("subscribe parameter EpochLength must be a positive integer")
	}
	return nil
}
