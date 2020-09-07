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
	DefaultSyncerDuration uint = 10
	DefaultEpochLength    uint = 5
)

var (
	DefaultMiningReward = sdk.NewInt(10000000000000)
	DefaultPullerReward = sdk.NewInt(500000000000)
)

// nolint - Keys for parameter access
var (
	KeySyncerDuration = []byte("SyncerDuration")
	KeyEpochLength    = []byte("EpochLength")
	KeyMiningReward   = []byte("MiningReward")
	KeyPullerReward   = []byte("PullerReward")
)

var _ params.ParamSet = (*Params)(nil)

type Params struct {
	SyncerDuration uint    `json:"syncer_duration" yaml:"syncer_duration"`
	EpochLength    uint    `json:"epoch_length" yaml:"epoch_length"`
	MiningReward   sdk.Int `json:"mining_reward" yaml:"mining_reward"`
	PullerReward   sdk.Int `json:"puller_reward" yaml:"puller_reward"`
}

// NewParams creates a new Params instance
func NewParams(syncerDuration, epochLength uint, miningReward, pullerReward sdk.Int) Params {

	return Params{
		SyncerDuration: syncerDuration,
		MiningReward:   miningReward,
		PullerReward:   pullerReward,
		EpochLength:    epochLength,
	}
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(KeySyncerDuration, &p.SyncerDuration, validateSyncerDuration),
		params.NewParamSetPair(KeyEpochLength, &p.EpochLength, validateEpochLength),
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
	return NewParams(DefaultSyncerDuration, DefaultEpochLength, DefaultMiningReward, DefaultPullerReward)
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	return fmt.Sprintf(`Params:
  SyncerDuration:    %d,
     EpochLength:    %d,
    MiningReward:    %s,
    PullerReward:    %s`,
		p.SyncerDuration, p.EpochLength, p.MiningReward, p.PullerReward)
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
	if p.SyncerDuration == 0 {
		return fmt.Errorf("validator parameter SyncerDuration must be a positive integer")
	}

	if p.EpochLength == 0 {
		return fmt.Errorf("validator parameter EpochLength must be a positive integer")
	}

	if !p.MiningReward.IsPositive() {
		return fmt.Errorf("validator parameter MiningReward must be a positive integer")
	}

	if !p.PullerReward.IsPositive() {
		return fmt.Errorf("validator parameter PullerReward must be a positive integer")
	}

	return nil
}

func validateSyncerDuration(i interface{}) error {
	v, ok := i.(uint)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("validator parameter SyncerDuration must be positive: %d", v)
	}

	return nil
}

func validateEpochLength(i interface{}) error {
	v, ok := i.(uint)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("validator parameter EpochLength must be positive: %d", v)
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
