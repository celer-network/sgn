package types

import (
	"fmt"
	"strings"

	"github.com/celer-network/sgn/mainchain"
	sdk "github.com/cosmos/cosmos-sdk/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking"
)

type Delegator struct {
	CandidateAddr  string  `json:"candidate_addr"`
	DelegatorAddr  string  `json:"delegator_addr"`
	DelegatedStake sdk.Int `json:"delegated_stake"`
}

func NewDelegator(candidateAddr, delegatorAddr string) Delegator {
	return Delegator{
		CandidateAddr: mainchain.FormatAddrHex(candidateAddr),
		DelegatorAddr: mainchain.FormatAddrHex(delegatorAddr),
	}
}

// implement fmt.Stringer
func (d Delegator) String() string {
	return strings.TrimSpace(fmt.Sprintf(`CandidateAddr: %s, DelegatorAddr: %s, DelegatedStake: %v`,
		d.CandidateAddr, d.DelegatorAddr, d.DelegatedStake))
}

// valAccount will be used for running validator node, and transactors will be used for running gateway
type Candidate struct {
	EthAddress     string              `json:"eth_address"`
	ValAccount     sdk.AccAddress      `json:"val_account"`
	Transactors    []sdk.AccAddress    `json:"transactors"`
	Delegators     []Delegator         `json:"delegators"`
	StakingPool    sdk.Int             `json:"staking_pool"`
	CommissionRate sdk.Dec             `json:"commission_rate"`
	RequestCount   sdk.Int             `json:"request_count"`
	Description    staking.Description `json:"description"`
}

func NewCandidate(ethAddress string, acctAddress sdk.AccAddress) Candidate {
	return Candidate{
		EthAddress: mainchain.FormatAddrHex(ethAddress),
		ValAccount: acctAddress,
	}
}

// implement fmt.Stringer
func (c Candidate) String() string {
	return strings.TrimSpace(fmt.Sprintf(`ValAccount: %s, EthAddress: %s, StakingPool: %v`, c.ValAccount, c.EthAddress, c.StakingPool))
}
