package cli

import (
	"fmt"

	"github.com/celer-network/sgn/x/validator/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingCli "github.com/cosmos/cosmos-sdk/x/staking/client/cli"
	"github.com/spf13/cobra"
)

const (
	flagSeq = "seq"
)

func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	validatorQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the validator module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	validatorQueryCmd.AddCommand(client.GetCommands(
		GetCmdPuller(storeKey, cdc),
		GetCmdDelegator(storeKey, cdc),
		GetCmdCandidate(storeKey, cdc),
		stakingCli.GetCmdQueryValidator(staking.StoreKey, cdc),
		stakingCli.GetCmdQueryValidators(staking.StoreKey, cdc),
	)...)
	return validatorQueryCmd
}

// GetCmdPuller queries puller info
func GetCmdPuller(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "puller",
		Short: "query puller info",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			puller, err := QueryPuller(cdc, cliCtx, queryRoute)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(puller)
		},
	}
}

// Query puller info
func QueryPuller(cdc *codec.Codec, cliCtx context.CLIContext, queryRoute string) (puller types.Puller, err error) {
	route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryPuller)
	res, _, err := cliCtx.Query(route)
	if err != nil {
		fmt.Printf("query puller error", err)
		return
	}

	cdc.MustUnmarshalJSON(res, &puller)
	return
}

// GetCmdPusher queries pusher info
func GetCmdPusher(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "pusher",
		Short: "query pusher info",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			pusher, err := QueryPusher(cdc, cliCtx, queryRoute)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(pusher)
		},
	}
}

// Query pusher info
func QueryPusher(cdc *codec.Codec, cliCtx context.CLIContext, queryRoute string) (pusher types.Pusher, err error) {
	route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryPusher)
	res, _, err := cliCtx.Query(route)
	if err != nil {
		fmt.Printf("query pusher error", err)
		return
	}

	cdc.MustUnmarshalJSON(res, &pusher)
	return
}

// GetCmdDelegator queries request info
func GetCmdDelegator(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "delegator [candidateAddress] [delegatorAddress]",
		Short: "query delegator info by candidateAddress and delegatorAddress",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := cdc.MarshalJSON(types.NewQueryDelegatorParams(args[0], args[1]))
			if err != nil {
				return err
			}

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryDelegator)
			bz, _, err := cliCtx.QueryWithData(route, data)
			if err != nil {
				return err
			}

			var delegator types.Delegator
			cdc.MustUnmarshalJSON(bz, &delegator)
			return cliCtx.PrintOutput(delegator)
		},
	}
}

// GetCmdCandidate queries request info
func GetCmdCandidate(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "candidate [candidateAddress]",
		Short: "query candidate info by candidateAddress and seq",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := cdc.MarshalJSON(types.NewQueryCandidateParams(args[0]))
			if err != nil {
				return err
			}

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryCandidate)
			bz, _, err := cliCtx.QueryWithData(route, data)
			if err != nil {
				return err
			}

			var candidate types.Candidate
			cdc.MustUnmarshalJSON(bz, &candidate)
			return cliCtx.PrintOutput(candidate)
		},
	}

	return cmd
}
