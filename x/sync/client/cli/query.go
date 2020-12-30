package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/x/sync/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/version"
)

const (
	flagStatus = "status"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group sync queries under a subcommand
	syncQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the sync module",
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	syncQueryCmd.AddCommand(common.GetCommands(
		GetCmdQueryChange(queryRoute, cdc),
		GetCmdQueryChanges(queryRoute, cdc),
		GetCmdQueryParam(queryRoute, cdc),
		GetCmdQueryParams(queryRoute, cdc),
	)...)

	return syncQueryCmd
}

// GetCmdQueryChange implements the query change command.
func GetCmdQueryChange(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "change [change-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query details of a single change",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details for a change. You can find the
change-id by running "%s query sync changes".

Example:
$ %s query sync change 1
`,
				version.ClientName, version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := common.NewQueryCLIContext(cdc)

			// validate that the change id is a uint
			changeID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("change-id %s not a valid uint, please input a valid change-id", args[0])
			}

			// Query the change
			change, err := QueryChange(cliCtx, queryRoute, changeID)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(change) // nolint:errcheck
		},
	}
}

func QueryChange(cliCtx context.CLIContext, queryRoute string, changeID uint64) (change types.Change, err error) {
	cdc := cliCtx.Codec
	params := types.NewQueryChangeParams(changeID)
	bz, err := cdc.MarshalJSON(params)
	if err != nil {
		return
	}

	res, err := common.RobustQueryWithData(cliCtx, fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryChange), bz)
	if err != nil {
		return
	}

	cdc.MustUnmarshalJSON(res, &change)
	return
}

// GetCmdQueryChanges implements a query changes command.
func GetCmdQueryChanges(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "changes",
		Short: "Query changes with optional filters",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for a all  changes:
Example:
$ query sync changes
`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := common.NewQueryCLIContext(cdc)
			changes, err := QueryChanges(cliCtx, queryRoute)
			if err != nil {
				return err
			}

			if len(changes) == 0 {
				return fmt.Errorf("no matching changes found")
			}

			return cliCtx.PrintOutput(changes) // nolint:errcheck
		},
	}

	return cmd
}

func QueryChanges(cliCtx context.CLIContext, queryRoute string) (changes types.Changes, err error) {
	cdc := cliCtx.Codec

	res, err := common.RobustQuery(cliCtx, fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryChanges))
	if err != nil {
		return
	}

	err = cdc.UnmarshalJSON(res, &changes)
	return
}

// GetCmdQueryChange implements the query change command.
func GetCmdQueryParams(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Query the parameters of the sync process",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the all the parameters for the sync process.

Example:
$ %s query sync params
`,
				version.ClientName,
			),
		),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := common.NewQueryCLIContext(cdc)
			tp, err := common.RobustQueryWithData(cliCtx, fmt.Sprintf("custom/%s/params/tallying", queryRoute), nil)
			if err != nil {
				return err
			}
			vp, err := common.RobustQueryWithData(cliCtx, fmt.Sprintf("custom/%s/params/voting", queryRoute), nil)
			if err != nil {
				return err
			}

			var tallyParams types.TallyParams
			cdc.MustUnmarshalJSON(tp, &tallyParams)
			var votingParams types.VotingParams
			cdc.MustUnmarshalJSON(vp, &votingParams)

			return cliCtx.PrintOutput(types.NewParams(votingParams, tallyParams))
		},
	}
}

// GetCmdQueryChange implements the query change command.
func GetCmdQueryParam(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "param [param-type]",
		Args:  cobra.ExactArgs(1),
		Short: "Query the parameters (voting|tallying) of the sync process",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the all the parameters for the sync process.

Example:
$ %s query sync param voting
$ %s query sync param tallying
`,
				version.ClientName, version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := common.NewQueryCLIContext(cdc)

			// Query store
			res, err := common.RobustQueryWithData(cliCtx, fmt.Sprintf("custom/%s/params/%s", queryRoute, args[0]), nil)
			if err != nil {
				return err
			}
			var out fmt.Stringer
			switch args[0] {
			case "voting":
				var param types.VotingParams
				cdc.MustUnmarshalJSON(res, &param)
				out = param
			case "tallying":
				var param types.TallyParams
				cdc.MustUnmarshalJSON(res, &param)
				out = param
			default:
				return fmt.Errorf("argument must be one of (voting|tallying|deposit), was %s", args[0])
			}

			return cliCtx.PrintOutput(out)
		},
	}
}

// DONTCOVER
