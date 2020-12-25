package cli

import (
	"github.com/spf13/cobra"

	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/x/global/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
)

// Change flags
const ()

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	globalTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "global transactions subcommands",
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	globalTxCmd.AddCommand(common.PostCommands()...)

	return globalTxCmd
}

// DONTCOVER
