package ops

import (
	"fmt"

	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	proposalDepositFlag     = "proposal-deposit"
	voteTimeoutFlag         = "vote-timeout"
	slashTimeoutFlag        = "slash-timeout"
	minValidatorNumFlag     = "min-validator-num"
	maxValidatorNumFlag     = "max-validator-num"
	minStakeInPoolFlag      = "min-stake-in-Pool"
	advanceNoticePeriodFlag = "advance-notice-period"
	migrationTimeFlag       = "migration-time"
	allParamsFlag           = "all-params"
	proposalIdFlag          = "proposal-id"
	checkVotesFlag          = "check-votes"
)

func flagSetGovParams() *pflag.FlagSet {
	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	fs.Bool(proposalDepositFlag, false, "proposal deposit")
	fs.Bool(voteTimeoutFlag, false, "vote timeout")
	fs.Bool(slashTimeoutFlag, false, "slash timeout")
	fs.Bool(minValidatorNumFlag, false, "minimal validator num")
	fs.Bool(maxValidatorNumFlag, false, "maximal validator num")
	fs.Bool(minStakeInPoolFlag, false, "minimal stakes to become a validator")
	fs.Bool(advanceNoticePeriodFlag, false, "advance notice period")
	fs.Bool(migrationTimeFlag, false, "migration time")
	return fs
}

func GovCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gov",
		Short: "Mainchain governance",
	}
	cmd.AddCommand(
		getParamsCommand(),
		getProposalCommand(),
		createParamProposalCommand(),
		voteParamProposalCommand(),
		confirmParamProposalCommand(),
	)
	return cmd
}

func createParamProposalCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-param-proposal",
		Short: "create parameter change proposal",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	cmd.Flags().AddFlagSet(flagSetGovParams())
	cmd.Flags().Int(proposalIdFlag, -1, "proposal id")
	return cmd
}

func voteParamProposalCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote-param-proposal",
		Short: "vote parameter change proposal",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	return cmd
}

func confirmParamProposalCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "confirm-param-proposal",
		Short: "confirm parameter change proposal",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	return cmd
}

func getProposalCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-param-proposal",
		Short: "get param change proposal infomation",
		RunE: func(cmd *cobra.Command, args []string) error {
			ethClient, err := common.NewEthClientFromConfig()
			if err != nil {
				return err
			}
			id := viper.GetInt(proposalIdFlag)
			if id == -1 {
				nextId, err2 := ethClient.DPoS.NextParamProposalId(&bind.CallOpts{})
				if err2 != nil {
					return err2
				}
				if nextId.Uint64() > 0 {
					id = int(nextId.Uint64() - 1)
				} else {
					return fmt.Errorf("no gov proposal")
				}
			}
			proposal := &mainchain.ParamProposal{}
			proposal, err = ethClient.GetParamProposal(int64(id))
			if err != nil {
				return err
			}
			if proposal.Status == 0 {
				return fmt.Errorf("proposal %d uninitialized", id)
			}
			fmt.Println("proposal Id:", id)
			fmt.Printf("proposer: %x", proposal.Proposer)
			fmt.Println("deposit:", proposal.Deposit)
			fmt.Println("vote deadline:", proposal.VoteDeadline)
			fmt.Println("record:", proposalParamName(proposal.Record.Uint64()))
			fmt.Println("new value:", proposal.NewValue)
			fmt.Println("status:", proposalStatusName(proposal.Status))

			checkVotes, err := cmd.Flags().GetBool(checkVotesFlag)
			if err != nil {
				return err
			}
			if checkVotes {
				fmt.Println("\n---- vote stats based on current validator set ----")
				yesVoters, noVoters, yesStakes, totalStakes, quorumStakes, err2 := ethClient.GetParamProposalVotes(int64(id))
				if err2 != nil {
					return err2
				}
				fmt.Println("yes voters:")
				for _, addr := range yesVoters {
					fmt.Printf("  %x", addr)
				}
				fmt.Println("no voters:")
				for _, addr := range noVoters {
					fmt.Printf("  %x", addr)
				}
				fmt.Println("yes stakes: ", yesStakes)
				fmt.Println("total stakes: ", totalStakes)
				fmt.Println("quorum stakes: ", quorumStakes)
			}

			return nil
		},
	}
	cmd.Flags().Int(proposalIdFlag, -1, "proposal id")
	cmd.Flags().Bool(checkVotesFlag, false, "check votes based on current validator set")
	return cmd
}

func getParamsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-params",
		Short: "Get mainchain parameters",
		RunE: func(cmd *cobra.Command, args []string) error {
			ethClient, err := common.NewEthClientFromConfig()
			if err != nil {
				return err
			}
			count := 0

			allParams, err := cmd.Flags().GetBool(allParamsFlag)
			if err != nil {
				return err
			}

			check, err := cmd.Flags().GetBool(proposalDepositFlag)
			if err != nil {
				return err
			}
			if check || allParams {
				count++
				value, err := ethClient.GetProposalDeposit()
				if err != nil {
					return nil
				}
				fmt.Println("proposal deposit:", value)
			}

			check, err = cmd.Flags().GetBool(voteTimeoutFlag)
			if err != nil {
				return err
			}
			if check || allParams {
				count++
				value, err := ethClient.GetGovernVoteTimeout()
				if err != nil {
					return nil
				}
				fmt.Println("vote timeout:", value)
			}

			check, err = cmd.Flags().GetBool(slashTimeoutFlag)
			if err != nil {
				return err
			}
			if check || allParams {
				count++
				value, err := ethClient.GetSlashTimeout()
				if err != nil {
					return nil
				}
				fmt.Println("slash timeout:", value)
			}

			check, err = cmd.Flags().GetBool(minValidatorNumFlag)
			if err != nil {
				return err
			}
			if check || allParams {
				count++
				value, err := ethClient.GetMinValidatorNum()
				if err != nil {
					return nil
				}
				fmt.Println("minimal validator num:", value)
			}

			check, err = cmd.Flags().GetBool(maxValidatorNumFlag)
			if err != nil {
				return err
			}
			if check || allParams {
				count++
				value, err := ethClient.GetMaxValidatorNum()
				if err != nil {
					return nil
				}
				fmt.Println("maximal validator num:", value)
			}

			check, err = cmd.Flags().GetBool(minStakeInPoolFlag)
			if err != nil {
				return err
			}
			if check || allParams {
				count++
				value, err := ethClient.GetMinStakeInPool()
				if err != nil {
					return nil
				}
				fmt.Println("minimal stakes to become a validator:", value)
			}

			check, err = cmd.Flags().GetBool(advanceNoticePeriodFlag)
			if err != nil {
				return err
			}
			if check || allParams {
				count++
				value, err := ethClient.GetAdvanceNoticePeriod()
				if err != nil {
					return nil
				}
				fmt.Println("advance notice period:", value)
			}

			check, err = cmd.Flags().GetBool(migrationTimeFlag)
			if err != nil {
				return err
			}
			if check || allParams {
				count++
				value, err := ethClient.GetMigrationTime()
				if err != nil {
					return nil
				}
				fmt.Println("migration time:", value)
			}

			if count == 0 {
				return fmt.Errorf("no param flag provided")
			}

			return nil
		},
	}
	cmd.Flags().AddFlagSet(flagSetGovParams())
	cmd.Flags().Bool(allParamsFlag, false, "query all parameters")

	return cmd
}

func paramFlagToRecordId() int64 {
	return -1
}

func proposalParamName(param uint64) string {
	switch param {
	case mainchain.ProposalDeposit:
		return "ProposalDeposit"
	case mainchain.GovernVoteTimeout:
		return "GovernVoteTimeout"
	case mainchain.SlashTimeout:
		return "SlashTimeout"
	case mainchain.MinValidatorNum:
		return "MinValidatorNum"
	case mainchain.MaxValidatorNum:
		return "MaxValidatorNum"
	case mainchain.MinStakeInPool:
		return "MinStakeInPool"
	case mainchain.AdvanceNoticePeriod:
		return "AdvanceNoticePeriod"
	case mainchain.MigrationTime:
		return "MigrationTime"
	}
	return "InvalidParam"
}

func proposalVoteName(vote uint8) string {
	switch vote {
	case mainchain.Unvoted:
		return "Unvoted"
	case mainchain.VoteYes:
		return "Yes"
	case mainchain.VoteNo:
		return "No"
	case mainchain.VoteAbstain:
		return "Abstain"
	}
	return "InvalidVote"
}

func proposalStatusName(status uint8) string {
	switch status {
	case mainchain.ProposalStatusUninitiated:
		return "Uninitiated"
	case mainchain.ProposalStatusVoting:
		return "Voting"
	case mainchain.ProposalStatusClosed:
		return "Closed"
	}
	return "InvalidStatus"
}
