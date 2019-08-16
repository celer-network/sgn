package cli

import (
	"fmt"

	"github.com/celer-network/sgn/x/subscribe/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
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
	)...)
	return subscribeQueryCmd
}

// GetCmdSubscription queries subscription info
func GetCmdSubscription(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "subscription",
		Short: "query subscription info associated with the eth address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			subscription, err := QuerySubscrption(cdc, cliCtx, queryRoute, args[0])
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(subscription)
		},
	}
}

func QuerySubscrption(cdc *codec.Codec, cliCtx context.CLIContext, queryRoute, ethAddress string) (subscription types.Subscription, err error) {
	data, err := cdc.MarshalJSON(types.NewQuerySubscrptionParams(ethAddress))
	if err != nil {
		return
	}

	route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QuerySubscrption)
	res, _, err := cliCtx.QueryWithData(route, data)
	if err != nil {
		fmt.Printf("query error", err)
		return
	}

	cdc.MustUnmarshalJSON(res, &subscription)
	return
}
