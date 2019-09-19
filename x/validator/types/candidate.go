package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Delegator struct {
	Stake sdk.Int `json:"stake"`
}

func NewDelegator(stake sdk.Int) Delegator {
	return Delegator{
		Stake: stake,
	}
}

// implement fmt.Stringer
func (c Delegator) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Stake: %v`, c.Stake))
}

type Candidate struct {
	EthAddress string      `json:"ethAddress"`
	Seq        sdk.Int     `json:"seq"`
	TotalStake sdk.Int     `json:"totalStake"`
	Delegators []Delegator `json:"delegators"`
}

func NewCandidate(ethAddress string, seq sdk.Int) Candidate {
	return Candidate{
		EthAddress: ethAddress,
		Seq:        seq,
	}
}

// implement fmt.Stringer
func (c Candidate) String() string {
	return strings.TrimSpace(fmt.Sprintf(`EthAddress: %s, Seq: %v, TotalStake: %v`, c.EthAddress, c.Seq, c.TotalStake))
}
