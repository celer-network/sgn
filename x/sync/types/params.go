package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	params "github.com/cosmos/cosmos-sdk/x/params/subspace"
)

// Default period for  voting
const (
	DefaultPeriod time.Duration = 15 * time.Second // 15 seconds
)

// Default sync params
var (
	DefaultThreshold = sdk.NewDecWithPrec(668, 3)
)

// Parameter store key
var (
	ParamStoreKeyVotingParams = []byte("votingparams")
	ParamStoreKeyTallyParams  = []byte("tallyparams")
)

// ParamKeyTable - Key declaration for parameters
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable(
		params.NewParamSetPair(ParamStoreKeyVotingParams, VotingParams{}, validateVotingParams),
		params.NewParamSetPair(ParamStoreKeyTallyParams, TallyParams{}, validateTallyParams),
	)
}

// TallyParams defines the params around Tallying votes in sync
type TallyParams struct {
	Threshold sdk.Dec `json:"threshold,omitempty" yaml:"threshold,omitempty"` //  Minimum proportion of Yes votes for change to pass. Initial value: 0.668
}

// NewTallyParams creates a new TallyParams object
func NewTallyParams(threshold sdk.Dec) TallyParams {
	return TallyParams{
		Threshold: threshold,
	}
}

// DefaultTallyParams default parameters for tallying
func DefaultTallyParams() TallyParams {
	return NewTallyParams(DefaultThreshold)
}

// String implements stringer insterface
func (tp TallyParams) String() string {
	return fmt.Sprintf(`Tally Params:
  Threshold:          %s`,
		tp.Threshold)
}

func validateTallyParams(i interface{}) error {
	v, ok := i.(TallyParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if !v.Threshold.IsPositive() {
		return fmt.Errorf("vote threshold must be positive: %s", v.Threshold)
	}
	if v.Threshold.GT(sdk.OneDec()) {
		return fmt.Errorf("vote threshold too large: %s", v)
	}

	return nil
}

// VotingParams defines the params around Voting in sync
type VotingParams struct {
	VotingPeriod time.Duration `json:"voting_period,omitempty" yaml:"voting_period,omitempty"` //  Length of the voting period.
}

// NewVotingParams creates a new VotingParams object
func NewVotingParams(votingPeriod time.Duration) VotingParams {
	return VotingParams{
		VotingPeriod: votingPeriod,
	}
}

// DefaultVotingParams default parameters for voting
func DefaultVotingParams() VotingParams {
	return NewVotingParams(DefaultPeriod)
}

// String implements stringer interface
func (vp VotingParams) String() string {
	return fmt.Sprintf(`Voting Params:
  Voting Period:      %s`, vp.VotingPeriod)
}

func validateVotingParams(i interface{}) error {
	v, ok := i.(VotingParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.VotingPeriod <= 0 {
		return fmt.Errorf("voting period must be positive: %s", v.VotingPeriod)
	}

	return nil
}

// Params returns all of the sync params
type Params struct {
	VotingParams VotingParams `json:"voting_params" yaml:"voting_params"`
	TallyParams  TallyParams  `json:"tally_params" yaml:"tally_params"`
}

func (gp Params) String() string {
	return gp.VotingParams.String() + "\n" +
		gp.TallyParams.String()
}

// NewParams creates a new sync Params instance
func NewParams(vp VotingParams, tp TallyParams) Params {
	return Params{
		VotingParams: vp,
		TallyParams:  tp,
	}
}

// DefaultParams default sync params
func DefaultParams() Params {
	return NewParams(DefaultVotingParams(), DefaultTallyParams())
}
