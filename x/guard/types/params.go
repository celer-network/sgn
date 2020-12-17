package types

import (
	"bytes"
	"fmt"

	"github.com/celer-network/sgn/mainchain"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	ec "github.com/ethereum/go-ethereum/common"
)

// guard params default values
const (
	// Default number of guards for guarding request
	DefaultRequestGuardCount uint64 = 3

	// Default minimal channel dispute timeout in mainchain blocks
	DefaultMinDisputeTimeout uint64 = 80000
)

var (
	// Default cost per request
	DefaultRequestCost = sdk.NewInt(10000000000000000)

	// Default Ledger address in a zero address that should not be used
	DefaultLedgerAddress string = mainchain.ZeroAddrHex
)

// nolint - Keys for parameter access
var (
	KeyRequestGuardCount = []byte("RequestGuardCount")
	KeyRequestCost       = []byte("RequestCost")
	KeyMinDisputeTimeout = []byte("MinDisputeTimeout")
	KeyLedgerAddress     = []byte("LedgerAddress")
)

var _ params.ParamSet = (*Params)(nil)

// Params defines the high level settings for guard
type Params struct {
	RequestGuardCount uint64  `json:"request_guard_count" yaml:"request_guard_count"` // request guard count
	RequestCost       sdk.Int `json:"request_cost" yaml:"request_cost"`               // request cost
	MinDisputeTimeout uint64  `json:"min_dispute_timeout" yaml:"min_dispute_timeout"` // minimal channel dispute timeout in mainchain blocks
	LedgerAddress     string  `json:"ledger_address" yaml:"ledger_address"`           // ledger contract address
}

// NewParams creates a new Params instance
func NewParams(requestGuardCount uint64, requestCost sdk.Int, minDisputeTimeout uint64, ledgerAddress string) Params {

	return Params{
		RequestGuardCount: requestGuardCount,
		RequestCost:       requestCost,
		MinDisputeTimeout: minDisputeTimeout,
		LedgerAddress:     ledgerAddress,
	}
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(KeyRequestGuardCount, &p.RequestGuardCount, validateRequestGuardCount),
		params.NewParamSetPair(KeyRequestCost, &p.RequestCost, validateRequestCost),
		params.NewParamSetPair(KeyMinDisputeTimeout, &p.MinDisputeTimeout, validateMinDisputeTimeout),
		params.NewParamSetPair(KeyLedgerAddress, &p.LedgerAddress, validateLedgerAddress),
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
	return NewParams(DefaultRequestGuardCount, DefaultRequestCost, DefaultMinDisputeTimeout, DefaultLedgerAddress)
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	return fmt.Sprintf(`Params:
  RequestGuardCount: %d,
  RequestCost:       %s
  MinDisputeTimeout: %d`, p.RequestGuardCount, p.RequestCost, p.MinDisputeTimeout)
}

// unmarshal the current guard params value from store key or panic
func MustUnmarshalParams(cdc *codec.Codec, value []byte) Params {
	params, err := UnmarshalParams(cdc, value)
	if err != nil {
		panic(err)
	}
	return params
}

// unmarshal the current guard params value from store key
func UnmarshalParams(cdc *codec.Codec, value []byte) (params Params, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &params)
	if err != nil {
		return
	}
	return
}

// validate a set of params
func (p Params) Validate() error {
	if err := validateRequestGuardCount(p.RequestGuardCount); err != nil {
		return err
	}

	if err := validateRequestCost(p.RequestCost); err != nil {
		return err
	}

	if err := validateMinDisputeTimeout(p.MinDisputeTimeout); err != nil {
		return err
	}

	return nil
}

func validateRequestGuardCount(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("guard parameter RequestGuardCount must be positive: %d", v)
	}

	return nil
}

func validateRequestCost(i interface{}) error {
	v, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("guard parameter RequestCost cannot be negative: %s", v)
	}

	return nil
}

func validateMinDisputeTimeout(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("guard parameter MinDisputeTimeout must be positive: %d", v)
	}

	return nil
}

func validateLedgerAddress(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if !ec.IsHexAddress(v) {
		return fmt.Errorf("invalid LedgerAddress: %s", v)
	}
	return nil
}
