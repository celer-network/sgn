package ops

import (
	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/x/validator"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func WithdrawCommand(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw",
		Short: "Withdraw delegated stake",
	}
	cmd.AddCommand(
		intendWithdrawCommand(),
		confirmWithdrawCommand(),
		withdrawFromUnbondedCandidateCommand(),
		withdrawReward(cdc),
	)
	return cmd
}

func intendWithdrawCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "intend",
		Short: "Send a withdrawal intent for the stake delegated to a candidate",
		RunE: func(cmd *cobra.Command, args []string) error {
			ethClient, err := common.NewEthClientFromConfig()
			if err != nil {
				return err
			}
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
	cmd.Flags().String(amountFlag, "", "Withdraw amount (integer in unit of CELR)")
	cmd.MarkFlagRequired(candidateFlag)
	cmd.MarkFlagRequired(amountFlag)
	return cmd
}

func confirmWithdrawCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "confirm",
		Short: "Confirm withdrawal intents for the stake delegated to a candidate",
		RunE: func(cmd *cobra.Command, args []string) error {
			ethClient, err := common.NewEthClientFromConfig()
			if err != nil {
				return err
			}
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
			ethClient, err := common.NewEthClientFromConfig()
			if err != nil {
				return err
			}
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
	cmd.Flags().String(amountFlag, "", "Withdraw amount (integer in unit of CELR)")
	cmd.MarkFlagRequired(candidateFlag)
	cmd.MarkFlagRequired(amountFlag)
	return cmd
}

func withdrawReward(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reward [eth-address]",
		Short: "Withdraw reward on mainchain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := common.NewQueryCLIContext(cdc)
			reward, err := validator.CLIQueryReward(cliCtx, validator.RouterKey, args[0])
			if err != nil {
				log.Errorln("query reward error", err)
				return err
			}

			ethClient, err := common.NewEthClientFromConfig()
			if err != nil {
				return err
			}

			log.Infof("Withdrawing reward for %s", args[0])

			_, err = ethClient.Transactor.TransactWaitMined(
				"WithdrawFromUnbondedCandidate",
				func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*ethtypes.Transaction, error) {
					return ethClient.SGN.RedeemReward(opts, reward.GetRewardRequest())
				},
			)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}
