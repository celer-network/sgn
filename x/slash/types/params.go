package types

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

const (
	DefaultSignedBlocksWindow = 100
)

// slash params default values
var (
	DefaultMinSignedPerWindow        = sdk.NewDecWithPrec(5, 1)
	DefaultSlashFractionDoubleSign   = sdk.NewDec(1).Quo(sdk.NewDec(20))
	DefaultSlashFractionDowntime     = sdk.NewDec(1).Quo(sdk.NewDec(100))
	DefaultSlashFractionGuardFailure = sdk.NewDec(1).Quo(sdk.NewDec(100))
	DefaultFallbackGuardReward       = sdk.NewDec(1).Quo(sdk.NewDec(2))
)

// nolint - Keys for parameter access
var (
	KeySignedBlocksWindow        = []byte("SignedBlocksWindow")
	KeyMinSignedPerWindow        = []byte("MinSignedPerWindow")
	KeySlashFractionDoubleSign   = []byte("SlashFractionDoubleSign")
	KeySlashFractionDowntime     = []byte("SlashFractionDowntime")
	KeySlashFractionGuardFailure = []byte("SlashFractionGuardFailure")
	KeyFallbackGuardReward       = []byte("FallbackGuardReward")
)

var _ params.ParamSet = (*Params)(nil)

// Params defines the high level settings for slash
type Params struct {
	SignedBlocksWindow        int64   `json:"signed_blocks_window" yaml:"signed_blocks_window"`
	MinSignedPerWindow        sdk.Dec `json:"min_signed_per_window" yaml:"min_signed_per_window"`
	SlashFractionDoubleSign   sdk.Dec `json:"slashFractionDoubleSign" yaml:"slashFractionDoubleSign"`
	SlashFractionDowntime     sdk.Dec `json:"slashFractionDowntime" yaml:"slashFractionDowntime"`
	SlashFractionGuardFailure sdk.Dec `json:"slashFractionGuardFailure" yaml:"slashFractionGuardFailure"`
	FallbackGuardReward       sdk.Dec `json:"fallbackGuardReward" yaml:"fallbackGuardReward"`
}

// NewParams creates a new Params instance
func NewParams(signedBlocksWindow int64, minSignedPerWindow,
	slashFractionDoubleSign, slashFractionDowntime, slashFractionGuardFailure, fallbackGuardReward sdk.Dec) Params {
	return Params{
		SignedBlocksWindow:        signedBlocksWindow,
		MinSignedPerWindow:        minSignedPerWindow,
		SlashFractionDoubleSign:   slashFractionDoubleSign,
		SlashFractionDowntime:     slashFractionDowntime,
		SlashFractionGuardFailure: slashFractionGuardFailure,
		FallbackGuardReward:       fallbackGuardReward,
	}
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{KeySignedBlocksWindow, &p.SignedBlocksWindow},
		{KeyMinSignedPerWindow, &p.MinSignedPerWindow},
		{KeySlashFractionDoubleSign, &p.SlashFractionDoubleSign},
		{KeySlashFractionDowntime, &p.SlashFractionDowntime},
		{KeySlashFractionGuardFailure, &p.SlashFractionGuardFailure},
		{KeyFallbackGuardReward, &p.FallbackGuardReward},
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
	return NewParams(DefaultSignedBlocksWindow, DefaultMinSignedPerWindow, DefaultSlashFractionDoubleSign, DefaultSlashFractionDowntime, DefaultSlashFractionGuardFailure, DefaultFallbackGuardReward)
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	return fmt.Sprintf(`Params:
  SignedBlocksWindow:    %d,
  MinSignedPerWindow:    %s,
  SlashFractionDoubleSign:    %s,
	SlashFractionDowntime:    %s
	SlashFractionGuardFailure:    %s
	FallbackGuardReward:    %s`,
		p.SignedBlocksWindow, p.MinSignedPerWindow,
		p.SlashFractionDoubleSign, p.SlashFractionDowntime, p.SlashFractionGuardFailure, p.FallbackGuardReward)
}

// unmarshal the current slash params value from store key or panic
func MustUnmarshalParams(cdc *codec.Codec, value []byte) Params {
	params, err := UnmarshalParams(cdc, value)
	if err != nil {
		panic(err)
	}
	return params
}

// unmarshal the current slash params value from store key
func UnmarshalParams(cdc *codec.Codec, value []byte) (params Params, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &params)
	if err != nil {
		return
	}
	return
}

// validate a set of params
func (p Params) Validate() error {
	if p.SignedBlocksWindow == 0 {
		return fmt.Errorf("slash parameter SignedBlocksWindow must be positive")
	}

	if p.MinSignedPerWindow.IsNegative() {
		return fmt.Errorf("slash parameter MinSignedPerWindow must be positive")
	}

	if p.MinSignedPerWindow.GT(sdk.OneDec()) {
		return fmt.Errorf("slash parameter MinSignedPerWindow must be less or equal than 1")
	}

	if p.SlashFractionDoubleSign.IsNegative() {
		return fmt.Errorf("slash parameter SlashFractionDoubleSign must be positive")
	}

	if p.SlashFractionDowntime.IsNegative() {
		return fmt.Errorf("slash parameter SlashFractionDowntime must be positive")
	}

	if p.SlashFractionGuardFailure.IsNegative() {
		return fmt.Errorf("slash parameter SlashFractionGuardFailure must be positive")
	}

	if p.FallbackGuardReward.IsNegative() {
		return fmt.Errorf("slash parameter FallbackGuardReward must be positive")
	}

	return nil
}
