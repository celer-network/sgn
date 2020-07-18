package main

import (
	"github.com/celer-network/sgn/app"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/ops"
	"github.com/celer-network/sgn/test/channel"
	tc "github.com/celer-network/sgn/test/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
)

func main() {
	cobra.EnableCommandSorting = false

	rootCmd := &cobra.Command{
		Use:   "sgnops",
		Short: "sgn ops utility",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			viper.SetConfigFile(viper.GetString(common.FlagConfig))
			return viper.ReadInConfig()
		},
	}

	// Construct Root Command
	rootCmd.AddCommand(
		tc.DeployCommand(),
		tc.AccountsCommand(),
		channel.ServeCommand(),
		ops.InitCandidateCommand(),
		ops.DelegateCommand(),
		ops.IntendWithdrawCommand(),
		ops.ConfirmWithdrawCommand(),
		ops.ConfirmUnbondedCandidateCommand(),
		ops.WithdrawFromUnbondedCandidateCommand(),
	)

	rootCmd.PersistentFlags().String(common.FlagConfig, "./config.json", "config path")

	executor := cli.PrepareMainCmd(rootCmd, "SGN", app.DefaultCLIHome)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}
