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

func intendWithdraw() error {
	ethClient, err := initEthClient()
	amount := new(big.Int)
	amount.SetString(viper.GetString(amountFlag), 10)
	candidate := ethcommon.HexToAddress(viper.GetString(candidateFlag))

	log.Infof(
		"Sending intent to withdraw amount %s from candidate %s",
		amount,
		candidate.Hex(),
	)
	receipt, err := ethClient.Transactor.TransactWaitMined(
		"IntendWithdraw",
		func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*ethtypes.Transaction, error) {
			return ethClient.DPoS.IntendWithdraw(opts, candidate, amount)
		},
	)
	if err != nil {
		return err
	}
	log.Infof("Transaction %x succeeded", receipt.TxHash)
	return nil
}

func IntendWithdrawCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "intend-withdraw",
		Short: "Send a withdraw intent for the stake delegated to a candidate",
		RunE: func(cmd *cobra.Command, args []string) error {
			return intendWithdraw()
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
