package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Delegator struct {
	EthAddress     string  `json:"ethAddress"`
	DelegatedStake sdk.Int `json:"delegatedStake"`
}

func NewDelegator(ethAddress string) Delegator {
	return Delegator{
		EthAddress: ethAddress,
	}
}

// implement fmt.Stringer
func (c Delegator) String() string {
	return strings.TrimSpace(fmt.Sprintf(`EthAddress: %s, DelegatedStake: %v`, c.EthAddress, c.DelegatedStake))
}

type Candidate struct {
	StakingPool sdk.Int     `json:"stakingPool"`
	Delegators  []Delegator `json:"delegators"`
}

func NewCandidate() Candidate {
	return Candidate{}
}

// implement fmt.Stringer
func (c Candidate) String() string {
	return strings.TrimSpace(fmt.Sprintf(`StakingPool: %v`, c.StakingPool))
}
