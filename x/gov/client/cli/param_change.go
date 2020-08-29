package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/transactor"
	govutils "github.com/celer-network/sgn/x/gov/client/utils"
	govtypes "github.com/celer-network/sgn/x/gov/types"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GetCmdSubmitProposal implements a command handler for submitting a parameter
// change proposal transaction.
func GetCmdSubmitParamChangeProposal(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "param-change [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a parameter change proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a parameter proposal along with an initial deposit.
The proposal details must be supplied via a JSON file. For values that contains
objects, only non-empty fields will be updated.

IMPORTANT: Currently parameter changes are evaluated but not validated, so it is
very important that any "value" change is valid (ie. correct type and within bounds)
for its respective parameter.

Proper vetting of a parameter change proposal should prevent this from happening
(no deposits should occur during the governance process), but it should be noted
regardless.

Example:
$ %s tx gov submit-proposal param-change <path/to/proposal.json>

Where proposal.json contains:

{
  "title": "Guard Param Change",
  "description": "Update guard request cost",
  "changes": [
    {
      "subspace": "guard",
      "key": "RequestCost",
      "value": "5000000000000"
    }
  ],
  "deposit": "10"
}
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			txr, err := transactor.NewCliTransactor(cdc, viper.GetString(flags.FlagHome))
			if err != nil {
				log.Error(err)
				return err
			}

			proposal, err := govutils.ParseParamChangeProposalJSON(cdc, args[0])
			if err != nil {
				log.Error(err)
				return err
			}

			content := govtypes.NewParameterProposal(proposal.Title, proposal.Description, proposal.Changes.ToParamChanges())

			msg := govtypes.NewMsgSubmitProposal(content, proposal.Deposit, txr.Key.GetAddress())
			if err := msg.ValidateBasic(); err != nil {
				log.Error(err)
				return err
			}

			txr.AddTxMsg(msg)
			time.Sleep(5 * time.Second)

			return nil
		},
	}

	return cmd
}
