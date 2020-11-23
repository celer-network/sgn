package cli

import (
	"fmt"
	"sort"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/types"

	v03 "github.com/celer-network/sgn/x/genutil/legacy/v0_3"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	extypes "github.com/cosmos/cosmos-sdk/x/genutil"
	tmtime "github.com/tendermint/tendermint/types/time"
)

const (
	flagChainID = "chain-id"
)

var migrationMap = extypes.MigrationMap{
	"v0.3": v03.Migrate,
}

// GetMigrationVersions get all migration version in a sorted slice.
func GetMigrationVersions() []string {
	versions := make([]string, len(migrationMap))

	var i int
	for version := range migrationMap {
		versions[i] = version
		i++
	}

	sort.Strings(versions)
	return versions
}

// MigrateGenesisCmd returns a command to execute genesis state migration.
func MigrateGenesisCmd(_ *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate [target-version] [genesis-file]",
		Short: "Migrate genesis to a specified target version",
		Long: `Migrate the source genesis into the target version and print to STDOUT.

Example:
$ sgnd migrate v0.3 /path/to/genesis.json --chain-id=sgnchain-3
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			target := args[0]
			importGenesis := args[1]

			genDoc, err := types.GenesisDocFromFile(importGenesis)
			if err != nil {
				return errors.Wrapf(err, "failed to read genesis document from file %s", importGenesis)
			}

			var initialState extypes.AppMap
			if err = cdc.UnmarshalJSON(genDoc.AppState, &initialState); err != nil {
				return errors.Wrap(err, "failed to JSON unmarshal initial genesis state")
			}

			migrationFunc := migrationMap[target]
			if migrationFunc == nil {
				return fmt.Errorf("unknown migration function for version: %s. Supported versions: %s", target, GetMigrationVersions())
			}

			newGenState := migrationFunc(initialState)

			genDoc.AppState, err = cdc.MarshalJSON(newGenState)
			if err != nil {
				return errors.Wrap(err, "failed to JSON marshal migrated genesis state")
			}

			genDoc.GenesisTime = tmtime.Now()

			chainID := cmd.Flag(flagChainID).Value.String()
			if chainID != "" {
				genDoc.ChainID = chainID
			}

			bz, err := cdc.MarshalJSONIndent(genDoc, "", "  ")
			if err != nil {
				return errors.Wrap(err, "failed to marshal genesis doc")
			}

			sortedBz, err := sdk.SortJSON(bz)
			if err != nil {
				return errors.Wrap(err, "failed to sort JSON genesis doc")
			}

			cmd.Println(string(sortedBz))
			return nil
		},
	}

	cmd.Flags().String(flagChainID, "", "override chain_id with this flag")

	return cmd
}
