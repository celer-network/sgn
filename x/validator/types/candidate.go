package types

import (
	"fmt"
	"strings"

	"github.com/celer-network/sgn/mainchain"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Delegator struct {
	EthAddress     string  `json:"ethAddress"`
	DelegatedStake sdk.Int `json:"delegatedStake"`
}

func NewDelegator(ethAddress string) Delegator {
	return Delegator{
		EthAddress: mainchain.FormatAddrHex(ethAddress),
	}
}

// implement fmt.Stringer
func (c Delegator) String() string {
	return strings.TrimSpace(fmt.Sprintf(`EthAddress: %s, DelegatedStake: %v`, c.EthAddress, c.DelegatedStake))
}

// operator will be used for running validator node, and transactor will be used for running gateway
type Candidate struct {
	Operator    sdk.AccAddress   `json:"operator"`
	Transactors []sdk.AccAddress `json:"transactors"`
	StakingPool sdk.Int          `json:"stakingPool"`
	Delegators  []Delegator      `json:"delegators"`
}

func NewCandidate(operator sdk.AccAddress) Candidate {
	return Candidate{
		Operator: operator,
	}
}

// implement fmt.Stringer
func (c Candidate) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Operator: %s, StakingPool: %v`, c.Operator, c.StakingPool))
}
