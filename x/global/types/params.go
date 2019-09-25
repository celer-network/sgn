package types

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// global params default values
const (
	// Default epoch length based on seconds
	DefaultEpochLength int64 = 60

	// Default cost per epoch, 1 CELR token per epoch
	DefaultCostPerEpoch int64 = 1000000000000000000

	// Default max block diff accepted when sync block
	DefaultMaxBlockDiff int64 = 2

	// Default number of blocks to confirm a block is safe
	DefaultConfirmationCount uint64 = 5
)

// nolint - Keys for parameter access
var (
	KeyEpochLength       = []byte("EpochLength")
	KeyMaxBlockDiff      = []byte("KeyMaxBlockDiff")
	KeyConfirmationCount = []byte("KeyConfirmationCount")
	KeyCostPerEpoch      = []byte("KeyCostPerEpoch")
)

var _ params.ParamSet = (*Params)(nil)

// Params defines the high level settings for global
type Params struct {
	EpochLength       int64   `json:"epochLength" yaml:"epochLength"`             // epoch length based on seconds
	MaxBlockDiff      int64   `json:"maxBlockDiff" yaml:"maxBlockDiff"`           // Max block diff accepted when sync block
	ConfirmationCount uint64  `json:"confirmationCount" yaml:"confirmationCount"` // Number of blocks to confirm a block is safe
	CostPerEpoch      sdk.Int `json:"costPerEpoch" yaml:"costPerEpoch"`           // The fee will be charged for subscription per epoch
}

// NewParams creates a new Params instance
func NewParams(epochLength, maxBlockDiff int64, confirmationCount uint64, costPerEpoch sdk.Int) Params {

	return Params{
		EpochLength:       epochLength,
		MaxBlockDiff:      maxBlockDiff,
		ConfirmationCount: confirmationCount,
		CostPerEpoch:      costPerEpoch,
	}
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{KeyEpochLength, &p.EpochLength},
		{KeyMaxBlockDiff, &p.MaxBlockDiff},
		{KeyConfirmationCount, &p.ConfirmationCount},
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
	return NewParams(DefaultEpochLength, DefaultMaxBlockDiff, DefaultConfirmationCount, sdk.NewInt(DefaultCostPerEpoch))
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	return fmt.Sprintf(`Params:
  EpochLength:    %d
	MaxBlockDiff:   %d
	ConfirmationCount:   %d
	CostPerEpoch:   %v`,
		p.EpochLength, p.MaxBlockDiff, p.ConfirmationCount, p.CostPerEpoch)
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
	if p.EpochLength <= 0 {
		return fmt.Errorf("global parameter EpochLength must be a positive integer")
	}

	if p.MaxBlockDiff < 0 {
		return fmt.Errorf("global parameter EpochLength must be a positive integer")
	}

	if p.CostPerEpoch.LTE(sdk. ZeroInt()) {
		return fmt.Errorf("global parameter CostPerEpoch must be a positive integer")
	}

	return nil
}
