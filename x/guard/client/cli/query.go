package cli

import (
	"fmt"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/x/guard/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagEpochId = "epochId"
)

func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	guardQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the guard module",
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	guardQueryCmd.AddCommand(common.GetCommands(
		GetCmdSubscription(storeKey, cdc),
		GetCmdRequest(storeKey, cdc),
		GetCmdEpoch(storeKey, cdc),
		GetCmdQueryParams(storeKey, cdc),
	)...)
	return guardQueryCmd
}

// GetCmdSubscription queries subscription info
func GetCmdSubscription(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "subscription [ethAddress]",
		Short: "query subscription info associated with the eth address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := common.NewQueryCLIContext(cdc)
			subscription, err := QuerySubscription(cliCtx, queryRoute, args[0])
			if err != nil {
				log.Errorln("query error", err)
				return err
			}

			return cliCtx.PrintOutput(subscription)
		},
	}
}

// Query subscription info
func QuerySubscription(cliCtx context.CLIContext, queryRoute, ethAddress string) (subscription types.Subscription, err error) {
	data, err := cliCtx.Codec.MarshalJSON(types.NewQuerySubscriptionParams(ethAddress))
	if err != nil {
		return
	}

	route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QuerySubscription)
	res, err := common.RobustQueryWithData(cliCtx, route, data)
	if err != nil {
		return
	}

	err = cliCtx.Codec.UnmarshalJSON(res, &subscription)
	return
}

// GetCmdRequest queries request info
func GetCmdRequest(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "request [channelId] [simplexReceiver]",
		Short: "query request info associated with the channelId",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := common.NewQueryCLIContext(cdc)
			request, err := QueryRequest(cliCtx, queryRoute, mainchain.Hex2Bytes(args[0]), args[1])
			if err != nil {
				log.Errorln("query error", err)
				return err
			}

			return cliCtx.PrintOutput(request)
		},
	}
}

// Query request info
func QueryRequest(cliCtx context.CLIContext, queryRoute string, channelId []byte, simplexReceiver string) (request types.Request, err error) {
	data, err := cliCtx.Codec.MarshalJSON(types.NewQueryRequestParams(channelId, simplexReceiver))
	if err != nil {
		return
	}

	route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryRequest)
	res, err := common.RobustQueryWithData(cliCtx, route, data)
	if err != nil {
		return
	}

	err = cliCtx.Codec.UnmarshalJSON(res, &request)
	return
}

// GetCmdQueryParams implements the params query command.
func GetCmdQueryParams(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current guard parameters information",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := common.NewQueryCLIContext(cdc)
			params, err := QueryParams(cliCtx, queryRoute)
			if err != nil {
				log.Errorln("query error", err)
				return err
			}

			return cliCtx.PrintOutput(params)
		},
	}
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

			cliCtx := common.NewQueryCLIContext(cdc)
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

// Query params info
func QueryParams(cliCtx context.CLIContext, queryRoute string) (params types.Params, err error) {
	route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryParameters)
	res, err := common.RobustQuery(cliCtx, route)
	if err != nil {
		return
	}

	err = cliCtx.Codec.UnmarshalJSON(res, &params)
	return
}
