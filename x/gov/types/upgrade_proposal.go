package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/upgrade"
)

const (
	ProposalTypeSoftwareUpgrade string = "SoftwareUpgrade"
)

// Software Upgrade Proposals
type SoftwareUpgradeProposal struct {
	Title       string       `json:"title" yaml:"title"`
	Description string       `json:"description" yaml:"description"`
	Plan        upgrade.Plan `json:"plan" yaml:"plan"`
}

func NewSoftwareUpgradeProposal(title, description string, plan upgrade.Plan) Content {
	return &SoftwareUpgradeProposal{title, description, plan}
}

// Implements Proposal Interface
var _ Content = SoftwareUpgradeProposal{}

func init() {
	RegisterProposalType(ProposalTypeSoftwareUpgrade)
}

func (sup SoftwareUpgradeProposal) GetTitle() string { return sup.Title }

func (sup SoftwareUpgradeProposal) GetDescription() string { return sup.Description }

func (sup SoftwareUpgradeProposal) ProposalRoute() string { return RouterKey }

func (sup SoftwareUpgradeProposal) ProposalType() string { return ProposalTypeSoftwareUpgrade }

func (sup SoftwareUpgradeProposal) ValidateBasic() error {
	if err := sup.Plan.ValidateBasic(); err != nil {
		return err
	}
	return ValidateAbstract(sup)
}

func (sup SoftwareUpgradeProposal) String() string {
	return fmt.Sprintf(`Software Upgrade Proposal:
  Title:       %s
  Description: %s
`, sup.Title, sup.Description)
}
