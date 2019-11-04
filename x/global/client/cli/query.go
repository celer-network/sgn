package cli

import (
	"encoding/binary"
	"fmt"

	"github.com/celer-network/sgn/x/global/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagEpochId = "epochId"
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
		GetCmdEpoch(storeKey, cdc),
		GetCmdQueryParams(storeKey, cdc),
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
			block, err := QueryLatestBlock(cliCtx, queryRoute)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(block)
		},
	}
}

// Query latest block
func QueryLatestBlock(cliCtx context.CLIContext, queryRoute string) (block types.Block, err error) {
	route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryLatestBlock)
	res, _, err := cliCtx.Query(route)
	if err != nil {
		fmt.Println("query error", err)
		return
	}

	cliCtx.Codec.MustUnmarshalJSON(res, &block)
	return
}

// Query secure block number
func QuerySecureBlockNum(cliCtx context.CLIContext, queryRoute string) (secureBlockNum uint64, err error) {
	route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QuerySecureBlockNum)
	res, _, err := cliCtx.Query(route)
	if err != nil {
		fmt.Println("query error", err)
		return
	}

	secureBlockNum = binary.BigEndian.Uint64(res)
	return
}

// GetCmdEpoch queries request info
func GetCmdEpoch(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "epoch",
		Short: "query epoch info by epochId",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			epochId := viper.GetInt64(flagEpochId)
			data, err := cdc.MarshalJSON(types.NewQueryEpochParams(epochId))
			if err != nil {
				return err
			}

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryEpoch)
			bz, _, err := cliCtx.QueryWithData(route, data)
			if err != nil {
				return err
			}

			var epoch types.Epoch
			cdc.MustUnmarshalJSON(bz, &epoch)
			return cliCtx.PrintOutput(epoch)
		},
	}

	cmd.Flags().Int64(flagEpochId, 0, "Epoch id")

	return cmd
}

// GetCmdQueryParams implements the params query command.
func GetCmdQueryParams(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current global parameters information",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryParameters)
			bz, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var params types.Params
			cdc.MustUnmarshalJSON(bz, &params)
			return cliCtx.PrintOutput(params)
		},
	}
}
