package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Delegator struct {
	EthAddress string  `json:"ethAddress"`
	Stake      sdk.Int `json:"stake"`
}

func NewDelegator(ethAddress string) Delegator {
	return Delegator{
		EthAddress: ethAddress,
	}
}

// implement fmt.Stringer
func (c Delegator) String() string {
	return strings.TrimSpace(fmt.Sprintf(`EthAddress: %s, Stake: %v`, c.EthAddress, c.Stake))
}

type Candidate struct {
	EthAddress string      `json:"ethAddress"`
	TotalStake sdk.Int     `json:"totalStake"`
	Delegators []Delegator `json:"delegators"`
}

func NewCandidate(ethAddress string, seq sdk.Int) Candidate {
	return Candidate{
		EthAddress: ethAddress,
	}
}

// implement fmt.Stringer
func (c Candidate) String() string {
	return strings.TrimSpace(fmt.Sprintf(`EthAddress: %s, TotalStake: %v`, c.EthAddress, c.TotalStake))
}
