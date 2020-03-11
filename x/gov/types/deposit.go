package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Deposit defines an amount deposited by an account address to an active proposal
type Deposit struct {
	ProposalID uint64         `json:"proposal_id" yaml:"proposal_id"` //  proposalID of the proposal
	Depositor  sdk.AccAddress `json:"depositor" yaml:"depositor"`     //  Address of the depositor
	Amount     sdk.Int        `json:"amount" yaml:"amount"`           //  Deposit amount
}

// NewDeposit creates a new Deposit instance
func NewDeposit(proposalID uint64, depositor sdk.AccAddress, amount sdk.Int) Deposit {
	return Deposit{proposalID, depositor, amount}
}

func (d Deposit) String() string {
	return fmt.Sprintf("deposit by %s on Proposal %d is for the amount %s",
		d.Depositor, d.ProposalID, d.Amount)
}

// Deposits is a collection of Deposit objects
type Deposits []Deposit

func (d Deposits) String() string {
	if len(d) == 0 {
		return "[]"
	}
	out := fmt.Sprintf("Deposits for Proposal %d:", d[0].ProposalID)
	for _, dep := range d {
		out += fmt.Sprintf("\n  %s: %s", dep.Depositor, dep.Amount)
	}
	return out
}

// Equals returns whether two deposits are equal.
func (d Deposit) Equals(comp Deposit) bool {
	return d.Depositor.Equals(comp.Depositor) && d.ProposalID == comp.ProposalID && d.Amount.Equal(comp.Amount)
}

// Empty returns whether a deposit is empty.
func (d Deposit) Empty() bool {
	return d.Equals(Deposit{})
}

type Depositor struct {
	Amount     sdk.Int   `json:"amount" yaml:"amount"`           //  Deposit amount
	MutedUntil time.Time `json:"muted_until" yaml:"muted_until"` // timestamp depositor cannot vote until
}

// NewDepositor creates a new Depositor instance
func NewDepositor() Depositor {
	return Depositor{
		Amount:     sdk.ZeroInt(),
		MutedUntil: time.Unix(0, 0),
	}
}

func (d Depositor) String() string {
	return fmt.Sprintf("Amount: %s, MutedUntil: %s", d.Amount, d.MutedUntil)
}
