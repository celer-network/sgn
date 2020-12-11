package cli

import (
	"fmt"
	"strconv"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
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
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	slashQueryCmd.AddCommand(common.GetCommands(
		GetCmdPenalty(storeKey, cdc),
		GetCmdPenalties(storeKey, cdc),
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

			cliCtx := common.NewQueryCLIContext(cdc)
			penalty, err := QueryPenalty(cliCtx, queryRoute, nonce)
			if err != nil {
				log.Errorln("query error", err)
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
	res, err := common.RobustQueryWithData(cliCtx, route, data)
	if err != nil {
		return
	}

	err = cliCtx.Codec.UnmarshalJSON(res, &penalty)
	return
}

// GetCmdPenalties queries penalty info
func GetCmdPenalties(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "penalties",
		Short: "query penalties info",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := common.NewQueryCLIContext(cdc)
			penalties, err := QueryPenalties(cliCtx, queryRoute)
			if err != nil {
				log.Errorln("query error", err)
				return err
			}

			return cliCtx.PrintOutput(penalties)
		},
	}
}

// Query penalties info
func QueryPenalties(cliCtx context.CLIContext, queryRoute string) (penalties []types.Penalty, err error) {
	route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryPenalties)
	res, err := common.RobustQuery(cliCtx, route)
	if err != nil {
		return
	}

	err = cliCtx.Codec.UnmarshalJSON(res, &penalties)
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

			cliCtx := common.NewQueryCLIContext(cdc)
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
			cliCtx := common.NewQueryCLIContext(cdc)

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryParameters)
			res, err := common.RobustQuery(cliCtx, route)
			if err != nil {
				return err
			}

			var params types.Params
			cdc.MustUnmarshalJSON(res, &params)
			return cliCtx.PrintOutput(params)
		},
	}
}
