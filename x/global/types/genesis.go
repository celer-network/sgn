package types

import (
	"bytes"
)

const (
	DefaultEthBlkNum = 0
)

// GenesisState - all staking state that must be provided at genesis
type GenesisState struct {
	EthBlkNum uint64 `json:"blk_num" yaml:"blk_num"`
}

// NewGenesisState creates a new genesis state for the global module
func NewGenesisState(ethBlkNum uint64) GenesisState {
	return GenesisState{
		EthBlkNum: ethBlkNum,
	}
}

// DefaultGenesisState defines the default global genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(DefaultEthBlkNum)
}

// Equal checks whether two global GenesisState structs are equivalent
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

	return nil
}
