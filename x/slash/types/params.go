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
	SlashFractionDoubleSign   sdk.Dec `json:"slash_fraction_double_sign" yaml:"slash_fraction_double_sign"`
	SlashFractionDowntime     sdk.Dec `json:"slash_fraction_downtime" yaml:"slash_fraction_downtime"`
	SlashFractionGuardFailure sdk.Dec `json:"slash_fraction_guard_failure" yaml:"slash_fraction_guard_failure"`
	FallbackGuardReward       sdk.Dec `json:"fallback_guard_reward" yaml:"fallback_guard_reward"`
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
		params.NewParamSetPair(KeySignedBlocksWindow, &p.SignedBlocksWindow, validateSignedBlocksWindow),
		params.NewParamSetPair(KeyMinSignedPerWindow, &p.MinSignedPerWindow, validateMinSignedPerWindow),
		params.NewParamSetPair(KeySlashFractionDoubleSign, &p.SlashFractionDoubleSign, validateSlashFractionDoubleSign),
		params.NewParamSetPair(KeySlashFractionDowntime, &p.SlashFractionDowntime, validateSlashFractionDowntime),
		params.NewParamSetPair(KeySlashFractionGuardFailure, &p.SlashFractionGuardFailure, validateSlashFractionGuardFailure),
		params.NewParamSetPair(KeyFallbackGuardReward, &p.FallbackGuardReward, validateFallbackGuardReward),
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
	if err := validateSignedBlocksWindow(p.SignedBlocksWindow); err != nil {
		return err
	}
	if err := validateMinSignedPerWindow(p.MinSignedPerWindow); err != nil {
		return err
	}
	if err := validateSlashFractionDoubleSign(p.SlashFractionDoubleSign); err != nil {
		return err
	}
	if err := validateSlashFractionDowntime(p.SlashFractionDowntime); err != nil {
		return err
	}
	if err := validateSlashFractionGuardFailure(p.SlashFractionGuardFailure); err != nil {
		return err
	}
	if err := validateFallbackGuardReward(p.FallbackGuardReward); err != nil {
		return err
	}

	return nil
}

func validateSignedBlocksWindow(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("slash parameter SignedBlocksWindow must be positive: %d", v)
	}

	return nil
}

func validateMinSignedPerWindow(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("slash parameter MinSignedPerWindow cannot be negative: %s", v)
	}

	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("slash parameter MinSignedPerWindow must be less or equal than 1: %s", v)
	}

	return nil
}

func validateSlashFractionDoubleSign(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("slash parameter SlashFractionDoubleSign cannot be negative: %s", v)
	}

	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("slash parameter SlashFractionDoubleSign must be less or equal than 1: %s", v)
	}

	return nil
}

func validateSlashFractionDowntime(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("slash parameter SlashFractionDowntime cannot be negative: %s", v)
	}

	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("slash parameter SlashFractionDowntime must be less or equal than 1: %s", v)
	}

	return nil
}

func validateSlashFractionGuardFailure(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("slash parameter SlashFractionGuardFailure cannot be negative: %s", v)
	}

	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("slash parameter SlashFractionGuardFailure must be less or equal than 1: %s", v)
	}

	return nil
}

func validateFallbackGuardReward(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("slash parameter FallbackGuardReward cannot be negative: %s", v)
	}

	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("slash parameter FallbackGuardReward must be less or equal than 1: %s", v)
	}

	return nil
}
