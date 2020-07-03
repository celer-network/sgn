package types

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// validator params default values
const (
	DefaultPullerDuration uint = 10
	DefaultPusherDuration uint = 10
)

var (
	DefaultMiningReward = sdk.NewInt(10000000000000)
	DefaultPullerReward = sdk.NewInt(500000000000)
)

// nolint - Keys for parameter access
var (
	KeyPullerDuration = []byte("PullerDuration")
	KeyPusherDuration = []byte("PusherDuration")
	KeyMiningReward   = []byte("MiningReward")
	KeyPullerReward   = []byte("PullerReward")
)

var _ params.ParamSet = (*Params)(nil)

type Params struct {
	PullerDuration uint    `json:"pullerDuration" yaml:"pullerDuration"`
	PusherDuration uint    `json:"pusherDuration" yaml:"pusherDuration"`
	MiningReward   sdk.Int `json:"miningReward" yaml:"miningReward"`
	PullerReward   sdk.Int `json:"pullerReward" yaml:"pullerReward"`
}

// NewParams creates a new Params instance
func NewParams(pullerDuration uint, pusherDuration uint, miningReward, pullerReward sdk.Int) Params {

	return Params{
		PullerDuration: pullerDuration,
		PusherDuration: pusherDuration,
		MiningReward:   miningReward,
		PullerReward:   pullerReward,
	}
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(KeyPullerDuration, &p.PullerDuration, validatePullerDuration),
		params.NewParamSetPair(KeyPusherDuration, &p.PusherDuration, validatePusherDuration),
		params.NewParamSetPair(KeyMiningReward, &p.MiningReward, validateMiningReward),
		params.NewParamSetPair(KeyPullerReward, &p.PullerReward, validatePullerReward),
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
	return NewParams(DefaultPullerDuration, DefaultPusherDuration, DefaultMiningReward, DefaultPullerReward)
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	return fmt.Sprintf(`Params:
  PullerDuration:    %d,
  PusherDuration:    %d,
	MiningReward:    %s
	PullerReward:    %s`,
		p.PullerDuration, p.PusherDuration, p.MiningReward, p.PullerReward)
}

// unmarshal the current validator params value from store key or panic
func MustUnmarshalParams(cdc *codec.Codec, value []byte) Params {
	params, err := UnmarshalParams(cdc, value)
	if err != nil {
		panic(err)
	}
	return params
}

// unmarshal the current validator params value from store key
func UnmarshalParams(cdc *codec.Codec, value []byte) (params Params, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &params)
	if err != nil {
		return
	}
	return
}

// validate a set of params
func (p Params) Validate() error {
	if p.PullerDuration == 0 {
		return fmt.Errorf("validator parameter PullerDuration must be a positive integer")
	}

	if p.PusherDuration == 0 {
		return fmt.Errorf("validator parameter PusherDuration must be a positive integer")
	}

	if !p.MiningReward.IsPositive() {
		return fmt.Errorf("validator parameter MiningReward must be a positive integer")
	}

	if !p.PullerReward.IsPositive() {
		return fmt.Errorf("validator parameter PullerReward must be a positive integer")
	}

	return nil
}

func validatePullerDuration(i interface{}) error {
	v, ok := i.(uint)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("validator parameter PullerDuration must be positive: %d", v)
	}

	return nil
}

func validatePusherDuration(i interface{}) error {
	v, ok := i.(uint)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("validator parameter PusherDuration must be positive: %d", v)
	}

	return nil
}

func validateMiningReward(i interface{}) error {
	v, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("guard parameter MiningReward cannot be negative: %s", v)
	}

	return nil
}

func validatePullerReward(i interface{}) error {
	v, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("guard parameter PullerReward cannot be negative: %s", v)
	}

	return nil
}
