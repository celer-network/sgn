package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/celer-network/sgn/app"
	"github.com/celer-network/sgn/common"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	tlog "github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

func GetSgndExecutor() cli.Executor {
	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()
	rootCmd := &cobra.Command{
		Use:               "sgnd",
		Short:             "SGN App Daemon (server)",
		PersistentPreRunE: persistentPreRunEFn(ctx),
	}
	// CLI commands to initialize the chain
	rootCmd.AddCommand(
		genutilcli.InitCmd(ctx, cdc, app.ModuleBasics, app.DefaultNodeHome),
		genutilcli.CollectGenTxsCmd(ctx, cdc, auth.GenesisAccountIterator{}, app.DefaultNodeHome),
		genutilcli.MigrateGenesisCmd(ctx, cdc),
		genutilcli.GenTxCmd(ctx, cdc, app.ModuleBasics, staking.AppModuleBasic{}, auth.GenesisAccountIterator{}, app.DefaultNodeHome, app.DefaultCLIHome),
		genutilcli.ValidateGenesisCmd(ctx, cdc, app.ModuleBasics),
		// addGenesisAccountCmd allows users to add accounts to the genesis file
		addGenesisAccountCmd(ctx, cdc, app.DefaultNodeHome, app.DefaultCLIHome),
	)

	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)
	rootCmd.PersistentFlags().String(common.FlagCLIHome, app.DefaultCLIHome, "Directory for cli config and data")
	rootCmd.PersistentFlags().String(
		common.FlagConfig, filepath.Join(app.DefaultCLIHome, "config", "sgn.toml"), "Path to SGN-specific configs")

	// prepare and add flags
	return cli.PrepareBaseCmd(rootCmd, "SGN", app.DefaultNodeHome)
}

func persistentPreRunEFn(context *server.Context) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		err := server.PersistentPreRunEFn(context)(cmd, args)
		if err != nil {
			return err
		}
		sgnConfigPath := viper.GetString(common.FlagConfig)
		_, err = os.Stat(sgnConfigPath)
		if err != nil {
			return err
		}
		viper.SetConfigFile(sgnConfigPath)
		return viper.MergeInConfig()
	}
}

func newApp(logger tlog.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	var cache sdk.MultiStorePersistentCache
	if viper.GetBool(server.FlagInterBlockCache) {
		cache = store.NewCommitKVStoreCacheManager()
	}
	skipUpgradeHeights := make(map[int64]bool)
	for _, h := range viper.GetIntSlice(server.FlagUnsafeSkipUpgrades) {
		skipUpgradeHeights[int64(h)] = true
	}
	pruningOpts, err := server.GetPruningOptionsFromFlags()
	if err != nil {
		panic(err)
	}

	app := app.NewSgnApp(
		logger,
		db,
		skipUpgradeHeights,
		baseapp.SetHaltHeight(viper.GetUint64(server.FlagHaltHeight)),
		baseapp.SetHaltTime(viper.GetUint64(server.FlagHaltTime)),
		baseapp.SetInterBlockCache(cache),
		baseapp.SetMinGasPrices(viper.GetString(server.FlagMinGasPrices)),
		baseapp.SetPruning(pruningOpts),
	)
	if err := app.LoadLatestVersion(); err != nil {
		panic(err)
	}
	return app
}

func exportAppStateAndTMValidators(
	logger tlog.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailWhiteList []string,
) (json.RawMessage, []tmtypes.GenesisValidator, error) {

	if height != -1 {
		sgnApp := app.NewSgnApp(logger, db, map[int64]bool{})
		err := sgnApp.LoadHeight(height)
		if err != nil {
			return nil, nil, err
		}
		return sgnApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
	}

	sgnApp := app.NewSgnApp(logger, db, map[int64]bool{})

	return sgnApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
}

func addGenesisAccountCmd(
	ctx *server.Context, cdc *codec.Codec, defaultNodeHome, defaultClientHome string,
) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "add-genesis-account [address_or_key_name] [coin][,[coin]]",
		Short: "Add a genesis account to genesis.json",
		Long: `Add a genesis account to genesis.json. The provided account must specify
the account address or key name and a list of initial coins. If a key name is given,
the address will be looked up in the local Keybase. The list of initial tokens must
contain valid denominations. Accounts may optionally be supplied with vesting parameters.
`,
		Args: cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				// attempt to lookup address from Keybase if no address was provided
				kb, err2 := keys.NewKeyBaseFromDir(viper.GetString(common.FlagCLIHome))
				if err2 != nil {
					return err2
				}

				info, err2 := kb.Get(args[0])
				if err2 != nil {
					return fmt.Errorf("failed to get address from Keybase: %w", err)
				}

				addr = info.GetAddress()
			}

			coins, err := sdk.ParseCoins(args[1])
			if err != nil {
				return fmt.Errorf("failed to parse coins: %w", err)
			}

			// create concrete account type based on input parameters
			var genAccount authexported.GenesisAccount

			baseAccount := auth.NewBaseAccount(addr, coins.Sort(), nil, 0, 0)

			genAccount = baseAccount

			if err = genAccount.Validate(); err != nil {
				return fmt.Errorf("failed to validate new genesis account: %w", err)
			}

			genFile := config.GenesisFile()
			appState, genDoc, err := genutil.GenesisStateFromGenFile(cdc, genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}

			authGenState := auth.GetGenesisStateFromAppState(cdc, appState)

			if authGenState.Accounts.Contains(addr) {
				return fmt.Errorf("cannot add account at existing address %s", addr)
			}

			// Add the new account to the set of genesis accounts and sanitize the
			// accounts afterwards.
			authGenState.Accounts = append(authGenState.Accounts, genAccount)
			authGenState.Accounts = auth.SanitizeGenesisAccounts(authGenState.Accounts)

			authGenStateBz, err := cdc.MarshalJSON(authGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal auth genesis state: %w", err)
			}

			appState[auth.ModuleName] = authGenStateBz

			appStateJSON, err := cdc.MarshalJSON(appState)
			if err != nil {
				return fmt.Errorf("failed to marshal application genesis state: %w", err)
			}

			genDoc.AppState = appStateJSON
			return genutil.ExportGenesisFile(genDoc, genFile)
		},
	}

	cmd.Flags().String(cli.HomeFlag, defaultNodeHome, "node's home directory")
	cmd.Flags().String(common.FlagCLIHome, defaultClientHome, "client's home directory")

	return cmd
}
