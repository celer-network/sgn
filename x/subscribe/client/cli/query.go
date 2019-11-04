package cli

import (
	"fmt"

	"github.com/celer-network/sgn/x/subscribe/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	subscribeQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the subscribe module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	subscribeQueryCmd.AddCommand(client.GetCommands(
		GetCmdSubscription(storeKey, cdc),
		GetCmdRequest(storeKey, cdc),
		GetCmdQueryParams(storeKey, cdc),
	)...)
	return subscribeQueryCmd
}

// GetCmdSubscription queries subscription info
func GetCmdSubscription(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "subscription [ethAddress]",
		Short: "query subscription info associated with the eth address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			subscription, err := QuerySubscription(cliCtx, queryRoute, args[0])
			if err != nil {
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
	res, _, err := cliCtx.QueryWithData(route, data)
	if err != nil {
		fmt.Println("query error", err)
		return
	}

	cliCtx.Codec.MustUnmarshalJSON(res, &subscription)
	return
}

// GetCmdRequest queries request info
func GetCmdRequest(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "request [channelId]",
		Short: "query request info associated with the channelId",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			request, err := QueryRequest(cliCtx, queryRoute, ethcommon.Hex2Bytes(args[0]))
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(request)
		},
	}
}

// Query request info
func QueryRequest(cliCtx context.CLIContext, queryRoute string, channelId []byte) (request types.Request, err error) {
	data, err := cliCtx.Codec.MarshalJSON(types.NewQueryRequestParams(channelId))
	if err != nil {
		return
	}

	route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryRequest)
	res, _, err := cliCtx.QueryWithData(route, data)
	if err != nil {
		fmt.Println("query error", err)
		return
	}

	cliCtx.Codec.MustUnmarshalJSON(res, &request)
	return
}

// GetCmdQueryParams implements the params query command.
func GetCmdQueryParams(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current subscribe parameters information",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryParameters)
			bz, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var params types.Params
			cliCtx.Codec.MustUnmarshalJSON(bz, &params)
			return cliCtx.PrintOutput(params)
		},
	}
}
