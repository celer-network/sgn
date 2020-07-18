package ops

import (
	"math/big"

	"github.com/celer-network/goutils/log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func withdrawFromUnbondedCandidate() error {
	ethClient, err := initEthClient()
	amount := new(big.Int)
	amount.SetString(viper.GetString(amountFlag), 10)
	candidate := ethcommon.HexToAddress(viper.GetString(candidateFlag))

	log.Infof(
		"Withdrawing amount %s delegated to an unbonded candidate %s",
		amount,
		candidate.Hex(),
	)
	receipt, err := ethClient.Transactor.TransactWaitMined(
		"WithdrawFromUnbondedCandidate",
		func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*ethtypes.Transaction, error) {
			return ethClient.DPoS.WithdrawFromUnbondedCandidate(opts, candidate, amount)
		},
	)
	if err != nil {
		return err
	}
	log.Infof("Transaction %x succeeded", receipt.TxHash.Hex())
	return nil
}

func WithdrawFromUnbondedCandidateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-from-unbonded-candidate",
		Short: "Withdraw delegated stake to an unbonded candidate",
		RunE: func(cmd *cobra.Command, args []string) error {
			return withdrawFromUnbondedCandidate()
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			err := cmd.MarkFlagRequired(candidateFlag)
			if err != nil {
				return err
			}
			err = cmd.MarkFlagRequired(amountFlag)
			if err != nil {
				return err
			}
			err = viper.BindPFlag(candidateFlag, cmd.Flags().Lookup(candidateFlag))
			if err != nil {
				return err
			}
			err = viper.BindPFlag(amountFlag, cmd.Flags().Lookup(amountFlag))
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().String(candidateFlag, "", "Candidate ETH address")
	cmd.Flags().String(amountFlag, "", "Withdraw amount")
	return cmd
}
