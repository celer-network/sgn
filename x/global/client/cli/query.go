package cli

import (
	"fmt"

	"github.com/celer-network/sgn/x/global/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	globalQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the global module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	globalQueryCmd.AddCommand(client.GetCommands(
		GetCmdLatestBlock(storeKey, cdc),
	)...)
	return globalQueryCmd
}

// GetCmdLatestBlock queries the lastest block
func GetCmdLatestBlock(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "latest-block",
		Short: "query last block info",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryLatestBlock)
			res, _, err := cliCtx.Query(route)
			if err != nil {
				fmt.Printf("query error", err)
				return nil
			}

			var out types.Block
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}
