package cli

import (
	"github.com/celer-network/sgn/x/validator/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingCli "github.com/cosmos/cosmos-sdk/x/staking/client/cli"
	"github.com/spf13/cobra"
)

func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	validatorQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the validator module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	validatorQueryCmd.AddCommand(client.GetCommands(
		stakingCli.GetCmdQueryValidator(staking.StoreKey, cdc),
		stakingCli.GetCmdQueryValidators(staking.StoreKey, cdc),
	)...)
	return validatorQueryCmd
}
