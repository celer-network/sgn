package main

import (
	"github.com/celer-network/sgn/cmd"
	"github.com/celer-network/sgn/common"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

func main() {
	cobra.EnableCommandSorting = false
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(common.Bech32PrefixAccAddr, common.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(common.Bech32PrefixValAddr, common.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(common.Bech32PrefixConsAddr, common.Bech32PrefixConsPub)
	config.Seal()

	// prepare and add flags
	executor := cmd.GetSgndExecutor()
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}
