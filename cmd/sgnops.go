package cmd

import (
	"github.com/celer-network/sgn/app"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/ops"
	"github.com/celer-network/sgn/testing/channel"
	tc "github.com/celer-network/sgn/testing/common"
	"github.com/celer-network/sgn/transactor"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
)

func GetSgnopsExecutor() cli.Executor {
	cdc := app.MakeCodec()

	rootCmd := &cobra.Command{
		Use:   "sgnops",
		Short: "sgn ops utility",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			viper.SetConfigFile(viper.GetString(common.FlagConfig))
			err := viper.ReadInConfig()
			if err != nil {
				return err
			}

			return common.SetupUserPassword()
		},
	}

	rootCmd.AddCommand(
		transactor.AccountsCommand(),
		tc.DeployCommand(),
		channel.ServeCommand(),
		ops.InitCandidateCommand(),
		ops.DelegateCommand(),
		ops.WithdrawCommand(),
		ops.ClaimValidatorCommand(),
		ops.ConfirmUnbondedCandidateCommand(),
		ops.UpdateMinSelfStakeCommand(),
		ops.UpdateCommissionRateCommand(),
		ops.GetCandidateInfoCommand(),
		ops.GetDelegatorInfoCommand(),
		ops.GetSyncCmd(cdc),
	)
	rootCmd.PersistentFlags().String(common.FlagConfig, "./config.json", "config path")

	return cli.PrepareMainCmd(rootCmd, "SGN", app.DefaultCLIHome)
}
