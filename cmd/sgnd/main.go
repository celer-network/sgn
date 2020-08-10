package main

import (
	"github.com/celer-network/sgn/cmd"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

func main() {
	cobra.EnableCommandSorting = false
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(sdk.Bech32PrefixAccAddr, sdk.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(sdk.Bech32PrefixValAddr, sdk.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(sdk.Bech32PrefixConsAddr, sdk.Bech32PrefixConsPub)
	config.Seal()

	// prepare and add flags
	executor := cmd.GetSgndExecutor()
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}
