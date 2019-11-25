package cli

import (
	"fmt"
	"strconv"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/x/slash/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	slashQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the slash module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	slashQueryCmd.AddCommand(client.GetCommands(
		GetCmdPenalty(storeKey, cdc),
		GetCmdQueryParams(storeKey, cdc),
	)...)
	return slashQueryCmd
}

// GetCmdPenalty queries penalty info
func GetCmdPenalty(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "penalty [nonce]",
		Short: "query penalty info by nonce",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			nonce, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			penalty, err := QueryPenalty(cliCtx, queryRoute, nonce)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(penalty)
		},
	}
}

// Query penalty info
func QueryPenalty(cliCtx context.CLIContext, queryRoute string, nonce uint64) (penalty types.Penalty, err error) {
	data, err := cliCtx.Codec.MarshalJSON(types.NewQueryPenaltyParams(nonce))
	if err != nil {
		return
	}

	route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryPenalty)
	res, _, err := cliCtx.QueryWithData(route, data)
	if err != nil {
		log.Errorln("query error", err)
		return
	}

	cliCtx.Codec.MustUnmarshalJSON(res, &penalty)
	return
}

// GetCmdPenaltyRequest queries penalty request proto
func GetCmdPenaltyRequest(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "penalty-request [nonce]",
		Short: "query penalty request proto by nonce",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			nonce, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			penaltyRequest, err := QueryPenaltyRequest(cliCtx, queryRoute, nonce)
			if err != nil {
				return err
			}

			log.Info(string(penaltyRequest))
			return nil
		},
	}
}

// Query penalty info
func QueryPenaltyRequest(cliCtx context.CLIContext, queryRoute string, nonce uint64) (penaltyRequest []byte, err error) {
	penalty, err := QueryPenalty(cliCtx, queryRoute, nonce)
	if err != nil {
		return
	}

	penaltyRequest = penalty.GetPenaltyRequest()
	return
}

// GetCmdQueryParams implements the params query command.
func GetCmdQueryParams(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current slash parameters information",
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
