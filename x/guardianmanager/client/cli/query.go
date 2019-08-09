package cli

import (
	"fmt"

	"github.com/celer-network/sgn/x/guardianmanager/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	guardianmanagerQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the guardianmanager module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	guardianmanagerQueryCmd.AddCommand(client.GetCommands(
		GetCmdGuardian(storeKey, cdc),
	)...)
	return guardianmanagerQueryCmd
}

// GetCmdGuardian queries the guardian associated with the eth address
func GetCmdGuardian(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "guardian",
		Short: "query guardian associated with the eth address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			data, err := cdc.MarshalJSON(types.NewQueryGuardianParams(args[0]))
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryGuardian)
			res, _, err := cliCtx.QueryWithData(route, data)
			if err != nil {
				fmt.Printf("query error", err)
				return nil
			}

			var out types.Guardian
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}
