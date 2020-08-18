package ops

import (
	"github.com/celer-network/goutils/log"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func WithdrawCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw",
		Short: "Withdraw delegated stake",
	}
	cmd.AddCommand(
		intendWithdrawCommand(),
		confirmWithdrawCommand(),
		withdrawFromUnbondedCandidateCommand(),
	)
	return cmd
}

func intendWithdrawCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "intend",
		Short: "Send a withdrawal intent for the stake delegated to a candidate",
		RunE: func(cmd *cobra.Command, args []string) error {
			ethClient, err := initEthClient()
			amount := calcRawAmount(viper.GetString(amountFlag))
			candidate := ethcommon.HexToAddress(viper.GetString(candidateFlag))

			log.Infof(
				"Sending intent to withdraw amount %s from candidate %s",
				amount,
				candidate.Hex(),
			)
			_, err = ethClient.Transactor.TransactWaitMined(
				"IntendWithdraw",
				func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*ethtypes.Transaction, error) {
					return ethClient.DPoS.IntendWithdraw(opts, candidate, amount)
				},
			)
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().String(candidateFlag, "", "Candidate ETH address")
	cmd.Flags().String(amountFlag, "", "Withdraw amount")
	cmd.MarkFlagRequired(candidateFlag)
	cmd.MarkFlagRequired(amountFlag)
	return cmd
}

func confirmWithdrawCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "confirm",
		Short: "Confirm withdrawal intents for the stake delegated to a candidate",
		RunE: func(cmd *cobra.Command, args []string) error {
			ethClient, err := initEthClient()
			candidate := ethcommon.HexToAddress(viper.GetString(candidateFlag))

			log.Infof(
				"Confirming withdrawal intents for the stake delegated to candidate %s",
				candidate.Hex(),
			)
			_, err = ethClient.Transactor.TransactWaitMined(
				"ConfirmWithdraw",
				func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*ethtypes.Transaction, error) {
					return ethClient.DPoS.ConfirmWithdraw(opts, candidate)
				},
			)
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().String(candidateFlag, "", "Candidate ETH address")
	cmd.MarkFlagRequired(candidateFlag)
	return cmd
}

func withdrawFromUnbondedCandidateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unbonded-candidate",
		Short: "Withdraw delegated stake from an unbonded candidate",
		RunE: func(cmd *cobra.Command, args []string) error {
			ethClient, err := initEthClient()
			amount := calcRawAmount(viper.GetString(amountFlag))
			candidate := ethcommon.HexToAddress(viper.GetString(candidateFlag))

			log.Infof(
				"Withdrawing amount %s delegated from an unbonded candidate %s",
				amount,
				candidate.Hex(),
			)
			_, err = ethClient.Transactor.TransactWaitMined(
				"WithdrawFromUnbondedCandidate",
				func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*ethtypes.Transaction, error) {
					return ethClient.DPoS.WithdrawFromUnbondedCandidate(opts, candidate, amount)
				},
			)
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().String(candidateFlag, "", "Candidate ETH address")
	cmd.Flags().String(amountFlag, "", "Withdraw amount")
	cmd.MarkFlagRequired(candidateFlag)
	cmd.MarkFlagRequired(amountFlag)
	return cmd
}
