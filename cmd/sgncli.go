package cmd

import (
	"os"
	"path/filepath"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/app"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/gateway"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankcmd "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/cli"
)

func GetSgncliExecutor() cli.Executor {
	cdc := app.MakeCodec()
	rootCmd := &cobra.Command{
		Use:   "sgncli",
		Short: "SGN node command line interface",
	}

	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		return initConfig(rootCmd)
	}

	// Construct Root Command
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		client.ConfigCmd(app.DefaultCLIHome),
		queryCmd(cdc),
		txCmd(cdc),
		flags.LineBreak,
		gateway.ServeCommand(cdc),
		flags.LineBreak,
		keys.Commands(),
		flags.LineBreak,
		version.Cmd,
	)

	rootCmd.PersistentFlags().String(
		common.FlagConfig, filepath.Join(app.DefaultCLIHome, "config", "sgn.toml"), "Path to SGN-specific configs")

	return cli.PrepareMainCmd(rootCmd, "SGN", app.DefaultCLIHome)

}

func queryCmd(cdc *amino.Codec) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Querying subcommands",
	}

	queryCmd.AddCommand(
		authcmd.GetAccountCmd(cdc),
		flags.LineBreak,
		rpc.ValidatorCommand(cdc),
		rpc.BlockCommand(),
		authcmd.QueryTxsByEventsCmd(cdc),
		authcmd.QueryTxCmd(cdc),
		flags.LineBreak,
	)

	// add modules' query commands
	app.ModuleBasics.AddQueryCommands(queryCmd, cdc)

	var cmdsToRemove []*cobra.Command

	for _, cmd := range queryCmd.Commands() {
		if cmd.Use == auth.ModuleName || cmd.Use == supply.ModuleName || cmd.Use == staking.ModuleName {
			cmdsToRemove = append(cmdsToRemove, cmd)
		}
	}

	queryCmd.RemoveCommand(cmdsToRemove...)

	return queryCmd
}

func txCmd(cdc *amino.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "tx",
		Short: "Transactions subcommands",
	}

	txCmd.AddCommand(
		bankcmd.SendTxCmd(cdc),
		flags.LineBreak,
		authcmd.GetSignCommand(cdc),
		authcmd.GetMultiSignCommand(cdc),
		flags.LineBreak,
		authcmd.GetBroadcastCommand(cdc),
		authcmd.GetEncodeCommand(cdc),
		authcmd.GetDecodeCommand(cdc),
		flags.LineBreak,
	)

	// add modules' tx commands
	app.ModuleBasics.AddTxCommands(txCmd, cdc)

	// remove auth and bank commands as they're mounted under the root tx command
	var cmdsToRemove []*cobra.Command

	for _, cmd := range txCmd.Commands() {
		if cmd.Use == auth.ModuleName || cmd.Use == bank.ModuleName || cmd.Use == staking.ModuleName {
			cmdsToRemove = append(cmdsToRemove, cmd)
		}
	}

	txCmd.RemoveCommand(cmdsToRemove...)

	return txCmd
}

func initConfig(cmd *cobra.Command) error {
	home, err := cmd.PersistentFlags().GetString(cli.HomeFlag)
	if err != nil {
		return err
	}
	cfgFile := filepath.Join(home, "config", "config.toml")
	_, err = os.Stat(cfgFile)
	if err == nil {
		viper.SetConfigFile(cfgFile)
		readErr := viper.ReadInConfig()
		if readErr != nil {
			return readErr
		}
	}
	sgnCfgFile := viper.GetString(common.FlagConfig)
	_, err = os.Stat(sgnCfgFile)
	if err != nil {
		return err
	}
	viper.SetConfigFile(sgnCfgFile)
	err = viper.MergeInConfig()
	if err != nil {
		return err
	}

	err = viper.BindPFlag(cli.EncodingFlag, cmd.PersistentFlags().Lookup(cli.EncodingFlag))
	if err != nil {
		return err
	}

	log.SetLevelByName(viper.GetString(common.FlagLogLevel))
	if viper.GetBool(common.FlagLogColor) {
		log.EnableColor()
	}

	return viper.BindPFlag(cli.OutputFlag, cmd.PersistentFlags().Lookup(cli.OutputFlag))
}
