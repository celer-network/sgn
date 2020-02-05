package main

import (
	"github.com/celer-network/sgn/app"
	"github.com/celer-network/sgn/common"
	tc "github.com/celer-network/sgn/test/common"
	"github.com/celer-network/sgn/test/osp"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/cli"
)

func main() {
	cobra.EnableCommandSorting = false

	rootCmd := &cobra.Command{
		Use:   "sgntest",
		Short: "sgn test utility",
	}

	// Construct Root Command
	rootCmd.AddCommand(
		tc.DeployCommand(),
		osp.ServeCommand(),
	)

	rootCmd.PersistentFlags().String(common.FlagConfig, "./config.json", "config path")

	executor := cli.PrepareMainCmd(rootCmd, "SGN", app.DefaultCLIHome)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}
