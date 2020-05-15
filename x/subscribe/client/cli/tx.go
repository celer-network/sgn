package cli

import (
	"github.com/celer-network/sgn/x/subscribe/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	subscribeTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "subscribe transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	subscribeTxCmd.AddCommand(flags.PostCommands()...)

	return subscribeTxCmd
}
