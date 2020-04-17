package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/celer-network/sgn/x/sync/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group sync queries under a subcommand
	syncQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the sync module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	syncQueryCmd.AddCommand(flags.GetCommands(
		GetCmdQueryChange(queryRoute, cdc),
		GetCmdQueryChanges(queryRoute, cdc),
		GetCmdQueryParam(queryRoute, cdc),
		GetCmdQueryParams(queryRoute, cdc),

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
			cliCtx := context.NewCLIContext().WithCodec(cdc)

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
	res, err := gcutils.QueryChangeByID(changeID, cliCtx, queryRoute)
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
			fmt.Sprintf(`Query for a all paginated changes that match optional filters:

Example:
$ %s query sync changes --depositor cosmos1skjwj5whet0lpe65qaq4rpq03hjxlwd9nf39lk
$ %s query sync changes --voter cosmos1skjwj5whet0lpe65qaq4rpq03hjxlwd9nf39lk
$ %s query sync changes --status (DepositPeriod|VotingPeriod|Passed|Rejected)
$ %s query sync changes --page=2 --limit=100
`,
				version.ClientName, version.ClientName, version.ClientName, version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			bechDepositorAddr := viper.GetString(flagDepositor)
			bechVoterAddr := viper.GetString(flagVoter)
			strChangeStatus := viper.GetString(flagStatus)
			page := viper.GetInt(flags.FlagPage)
			limit := viper.GetInt(flags.FlagLimit)

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			matchingChanges, err := QueryChanges(cliCtx, queryRoute, page, limit, bechDepositorAddr, bechVoterAddr, strChangeStatus)
			if err != nil {
				return err
			}

			if len(matchingChanges) == 0 {
				return fmt.Errorf("no matching changes found")
			}

			return cliCtx.PrintOutput(matchingChanges) // nolint:errcheck
		},
	}

	cmd.Flags().Int(flags.FlagPage, 1, "pagination page of changes to to query for")
	cmd.Flags().Int(flags.FlagLimit, 100, "pagination limit of changes to query for")
	cmd.Flags().String(flagDepositor, "", "(optional) filter by changes deposited on by depositor")
	cmd.Flags().String(flagVoter, "", "(optional) filter by changes voted on by voted")
	cmd.Flags().String(flagStatus, "", "(optional) filter changes by change status, status: deposit_period/voting_period/passed/rejected")

	return cmd
}

func QueryChanges(cliCtx context.CLIContext, queryRoute string, page, limit int, bechDepositorAddr, bechVoterAddr, strChangeStatus string) (matchingChanges types.Changes, err error) {
	cdc := cliCtx.Codec
	var depositorAddr sdk.AccAddress
	var voterAddr sdk.AccAddress
	var changeStatus types.ChangeStatus

	params := types.NewQueryChangesParams(page, limit, changeStatus, voterAddr, depositorAddr)

	if len(bechDepositorAddr) != 0 {
		depositorAddr, err2 := sdk.AccAddressFromBech32(bechDepositorAddr)
		if err2 != nil {
			err = err2
			return
		}
		params.Depositor = depositorAddr
	}

	if len(bechVoterAddr) != 0 {
		voterAddr, err2 := sdk.AccAddressFromBech32(bechVoterAddr)
		if err2 != nil {
			err = err2
			return
		}
		params.Voter = voterAddr
	}

	if len(strChangeStatus) != 0 {
		changeStatus, err2 := types.ChangeStatusFromString(gcutils.NormalizeChangeStatus(strChangeStatus))
		if err2 != nil {
			err = err2
			return
		}
		params.ChangeStatus = changeStatus
	}

	bz, err := cdc.MarshalJSON(params)
	if err != nil {
		return
	}

	res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/changes", queryRoute), bz)
	if err != nil {
		return
	}

	err = cdc.UnmarshalJSON(res, &matchingChanges)
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
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			tp, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/params/tallying", queryRoute), nil)
			if err != nil {
				return err
			}
			dp, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/params/deposit", queryRoute), nil)
			if err != nil {
				return err
			}
			vp, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/params/voting", queryRoute), nil)
			if err != nil {
				return err
			}

			var tallyParams types.TallyParams
			cdc.MustUnmarshalJSON(tp, &tallyParams)
			var depositParams types.DepositParams
			cdc.MustUnmarshalJSON(dp, &depositParams)
			var votingParams types.VotingParams
			cdc.MustUnmarshalJSON(vp, &votingParams)

			return cliCtx.PrintOutput(types.NewParams(votingParams, tallyParams, depositParams))
		},
	}
}

// GetCmdQueryChange implements the query change command.
func GetCmdQueryParam(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "param [param-type]",
		Args:  cobra.ExactArgs(1),
		Short: "Query the parameters (voting|tallying|deposit) of the sync process",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the all the parameters for the sync process.

Example:
$ %s query sync param voting
$ %s query sync param tallying
$ %s query sync param deposit
`,
				version.ClientName, version.ClientName, version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// Query store
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/params/%s", queryRoute, args[0]), nil)
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
			case "deposit":
				var param types.DepositParams
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
