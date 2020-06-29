package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/upgrade"
)

const (
	ProposalTypeSoftwareUpgrade string = "SoftwareUpgrade"
)

// Software Upgrade Proposals
type UpgradeProposal struct {
	Title       string       `json:"title" yaml:"title"`
	Description string       `json:"description" yaml:"description"`
	Plan        upgrade.Plan `json:"plan" yaml:"plan"`
}

func NewUpgradeProposal(title, description string, plan upgrade.Plan) Content {
	return &UpgradeProposal{title, description, plan}
}

// Implements Proposal Interface
var _ Content = UpgradeProposal{}

func init() {
	RegisterProposalType(ProposalTypeSoftwareUpgrade)
}

func (sup UpgradeProposal) GetTitle() string { return sup.Title }

func (sup UpgradeProposal) GetDescription() string { return sup.Description }

func (sup UpgradeProposal) ProposalRoute() string { return RouterKey }

func (sup UpgradeProposal) ProposalType() string { return ProposalTypeSoftwareUpgrade }

func (sup UpgradeProposal) ValidateBasic() error {
	if err := sup.Plan.ValidateBasic(); err != nil {
		return err
	}
	return ValidateAbstract(sup)
}

func (sup UpgradeProposal) String() string {
	return fmt.Sprintf(`Software Upgrade Proposal:
  Title:       %s
  Description: %s
`, sup.Title, sup.Description)
}
