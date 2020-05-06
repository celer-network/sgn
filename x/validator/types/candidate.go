package types

import (
	"fmt"
	"strings"

	"github.com/celer-network/sgn/mainchain"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Delegator struct {
	CandidateAddr  string  `json:"candidateAddr"`
	DelegatorAddr  string  `json:"delegatorAddr"`
	DelegatedStake sdk.Int `json:"delegatedStake"`
}

func NewDelegator(candidateAddr, delegatorAddr string) Delegator {
	return Delegator{
		CandidateAddr: mainchain.FormatAddrHex(candidateAddr),
		DelegatorAddr: mainchain.FormatAddrHex(delegatorAddr),
	}
}

// implement fmt.Stringer
func (d Delegator) String() string {
	return strings.TrimSpace(fmt.Sprintf(`CandidateAddr: %s, CandidateAddr: %s, DelegatedStake: %v`,
		d.CandidateAddr, d.DelegatorAddr, d.DelegatedStake))
}

// operator will be used for running validator node, and transactor will be used for running gateway
type Candidate struct {
	EthAddress   string           `json:"ethAddress"`
	Operator     sdk.AccAddress   `json:"operator"`
	Transactors  []sdk.AccAddress `json:"transactors"`
	Delegators   []Delegator      `json:"delegators"`
	StakingPool  sdk.Int          `json:"stakingPool"`
	RequestCount sdk.Int          `json:"requestCount"`
}

func NewCandidate(ethAddress string, operator sdk.AccAddress) Candidate {
	return Candidate{
		EthAddress: mainchain.FormatAddrHex(ethAddress),
		Operator:   operator,
	}
}

// implement fmt.Stringer
func (c Candidate) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Operator: %s, StakingPool: %v`, c.Operator, c.StakingPool))
}
