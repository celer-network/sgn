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

var migrationOrder = []string{"v0.2", "v0.3"}

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

func findMigrationIndex(version string) int {
	for i, v := range migrationOrder {
		if v == version {
			return i
		}
	}

	return -1
}

// MigrateGenesisCmd returns a command to execute genesis state migration.
func MigrateGenesisCmd(_ *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate [target-version] [genesis-file]",
		Short: "Migrate genesis to a specified target version",
		Long: fmt.Sprintf(`Migrate the source genesis into the target version and print to STDOUT.

Example:
$ sgnd migrate v0.2 v0.3 /path/to/genesis.json --chain-id=sgnchain-3

Supported versions: %s
`, migrationOrder),
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			source := args[0]
			target := args[1]
			importGenesis := args[2]

			sourceIndex := findMigrationIndex(source)
			if sourceIndex == -1 {
				return fmt.Errorf("invalid source version")
			}

			targetIndex := findMigrationIndex(target)
			if targetIndex == -1 {
				return fmt.Errorf("invalid target version")
			}

			if source >= target {
				return fmt.Errorf("target must be newer than source")
			}

			genDoc, err := types.GenesisDocFromFile(importGenesis)
			if err != nil {
				return errors.Wrapf(err, "failed to read genesis document from file %s", importGenesis)
			}

			var appState extypes.AppMap
			if err = cdc.UnmarshalJSON(genDoc.AppState, &appState); err != nil {
				return errors.Wrap(err, "failed to JSON unmarshal initial genesis state")
			}

			for current := sourceIndex + 1; current <= targetIndex; current++ {
				migrationFunc := migrationMap[migrationOrder[current]]
				appState = migrationFunc(appState)
			}

			genDoc.AppState, err = cdc.MarshalJSON(appState)
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
