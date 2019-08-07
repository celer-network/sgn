package cli

import (
	"fmt"

	"github.com/celer-network/sgn/x/bridge/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	bridgeQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the bridge module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	bridgeQueryCmd.AddCommand(client.GetCommands(
		GetCmdEthAddress(storeKey, cdc),
	)...)
	return bridgeQueryCmd
}

// GetCmdEthAddress queries a list of all names
func GetCmdEthAddress(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "eth-address",
		Short: "query eth address associated with the sidechain address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			data, err := cdc.MarshalJSON(types.NewQueryEthAddressParams(addr))
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryEthAddress)
			res, _, err := cliCtx.QueryWithData(route, data)
			if err != nil {
				fmt.Printf("query error", err)
				return nil
			}

			var out types.EthAddress
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}
