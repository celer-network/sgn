package types

import (
	"bytes"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all staking state that must be provided at genesis
type GenesisState struct {
	StartingChangeID uint64        `json:"starting_change_id" yaml:"starting_change_id"`
	Changes          Changes       `json:"changes" yaml:"changes"`
	DepositParams    DepositParams `json:"deposit_params" yaml:"deposit_params"`
	VotingParams     VotingParams  `json:"voting_params" yaml:"voting_params"`
	TallyParams      TallyParams   `json:"tally_params" yaml:"tally_params"`
}

// NewGenesisState creates a new genesis state for the sync module
func NewGenesisState(startingChangeID uint64, dp DepositParams, vp VotingParams, tp TallyParams) GenesisState {
	return GenesisState{
		StartingChangeID: startingChangeID,
		DepositParams:    dp,
		VotingParams:     vp,
		TallyParams:      tp,
	}
}

// DefaultGenesisState defines the default sync genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(
		DefaultStartingChangeID,
		DefaultDepositParams(),
		DefaultVotingParams(),
		DefaultTallyParams(),
	)
}

// Equal checks whether two sync GenesisState structs are equivalent
func (data GenesisState) Equal(data2 GenesisState) bool {
	b1 := ModuleCdc.MustMarshalBinaryBare(data)
	b2 := ModuleCdc.MustMarshalBinaryBare(data2)
	return bytes.Equal(b1, b2)
}

// IsEmpty returns true if a GenesisState is empty
func (data GenesisState) IsEmpty() bool {
	return data.Equal(GenesisState{})
}

// ValidateGenesis checks if parameters are within valid ranges
func ValidateGenesis(data GenesisState) error {
	threshold := data.TallyParams.Threshold
	if threshold.IsNegative() || threshold.GT(sdk.OneDec()) {
		return fmt.Errorf("sync vote threshold should be positive and less or equal to one, is %s",
			threshold.String())
	}

	veto := data.TallyParams.Veto
	if veto.IsNegative() || veto.GT(sdk.OneDec()) {
		return fmt.Errorf("sync vote veto threshold should be positive and less or equal to one, is %s",
			veto.String())
	}

	if data.DepositParams.MinDeposit.IsNegative() {
		return fmt.Errorf("sync deposit amount must not be a negative amount, is %s",
			data.DepositParams.MinDeposit.String())
	}

	return nil
}
