package cli

import (
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/x/guard/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	guardTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "guard transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	guardTxCmd.AddCommand(common.PostCommands()...)

	return guardTxCmd
}
