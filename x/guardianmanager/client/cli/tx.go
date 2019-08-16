package cli

import (
	"github.com/celer-network/sgn/x/guardianmanager/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	guardianmanagerTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "guardianmanager transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	guardianmanagerTxCmd.AddCommand(client.PostCommands(
		GetCmdSetEthAddress(cdc),
	)...)

	return guardianmanagerTxCmd
}

// GetCmdSetEthAddress is the CLI command for sending a SetEthAddress transaction
func GetCmdSetEthAddress(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "set-eth-address [eth-addr]",
		Short: "set the eth address associated with the from address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			msg := types.NewMsgDeposit(args[0], cliCtx.GetFromAddress())
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdRequestGuard is the CLI command for sending a SetEthAddress transaction
func GetCmdRequestGuard(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "set-eth-address [eth-addr]",
		Short: "set the eth address associated with the from address",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			signedSimplexStateBytes := ethcommon.Hex2Bytes(args[1])
			msg := types.NewMsgRequestGuard(args[0], signedSimplexStateBytes, cliCtx.GetFromAddress())
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
