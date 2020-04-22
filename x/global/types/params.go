package types

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// global params default values
const (
	// Default epoch length based on seconds
	DefaultEpochLength int64 = 60

	// Default number of blocks to confirm a block is safe
	DefaultConfirmationCount uint64 = 5
)

// nolint - Keys for parameter access
var (
	KeyEpochLength       = []byte("EpochLength")
	KeyConfirmationCount = []byte("KeyConfirmationCount")
)

var _ params.ParamSet = (*Params)(nil)

// Params defines the high level settings for global
type Params struct {
	EpochLength       int64  `json:"epochLength" yaml:"epochLength"`             // epoch length based on seconds
	ConfirmationCount uint64 `json:"confirmationCount" yaml:"confirmationCount"` // Number of blocks to confirm a block is safe
}

// NewParams creates a new Params instance
func NewParams(epochLength int64, confirmationCount uint64) Params {
	return Params{
		EpochLength:       epochLength,
		ConfirmationCount: confirmationCount,
	}
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(KeyEpochLength, &p.EpochLength, validateEpochLength),
		params.NewParamSetPair(KeyConfirmationCount, &p.ConfirmationCount, validateConfirmationCount),
	}
}

// Equal returns a boolean determining if two Param types are identical.
func (p Params) Equal(p2 Params) bool {
	bz1 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p)
	bz2 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(DefaultEpochLength, DefaultConfirmationCount)
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	return fmt.Sprintf(`Params:
  EpochLength:    %d
	ConfirmationCount:   %d`,
		p.EpochLength, p.ConfirmationCount)
}

// unmarshal the current global params value from store key or panic
func MustUnmarshalParams(cdc *codec.Codec, value []byte) Params {
	params, err := UnmarshalParams(cdc, value)
	if err != nil {
		panic(err)
	}
	return params
}

// unmarshal the current global params value from store key
func UnmarshalParams(cdc *codec.Codec, value []byte) (params Params, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &params)
	if err != nil {
		return
	}
	return
}

// validate a set of params
func (p Params) Validate() error {
	if err := validateEpochLength(p.EpochLength); err != nil {
		return err
	}

	if err := validateConfirmationCount(p.ConfirmationCount); err != nil {
		return err
	}

	return nil
}

func validateEpochLength(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("global parameter EpochLength must be positive: %d", v)
	}

	return nil
}

func validateConfirmationCount(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 0 {
		return fmt.Errorf("global parameter ConfirmationCount cannot be negative: %d", v)
	}

	return nil
}
