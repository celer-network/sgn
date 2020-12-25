package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/x/global/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group global queries under a subcommand
	globalQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the global module",
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	globalQueryCmd.AddCommand(common.GetCommands(
		GetCmdQueryEthBlkNum(queryRoute, cdc),
	)...)

	return globalQueryCmd
}

// GetCmdQueryEthBlkNum implements the query block number command.
func GetCmdQueryEthBlkNum(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "ethBlkNum",
		Args:  cobra.ExactArgs(0),
		Short: "Query latest mainchain block number globaled to sidechain",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := common.NewQueryCLIContext(cdc)

			res, err := common.RobustQuery(cliCtx, fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryEthBlkNum))
			if err != nil {
				return err
			}

			var ethBlkNum uint64
			cdc.MustUnmarshalBinaryBare(res, &ethBlkNum)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(ethBlkNum) // nolint:errcheck
		},
	}
}
