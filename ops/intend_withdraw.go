package ops

import (
	"github.com/celer-network/goutils/log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func intendWithdraw() error {
	ethClient, err := initEthClient()
	amount := calcRawAmount(viper.GetString(amountFlag))
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
	}
	cmd.Flags().String(candidateFlag, "", "Candidate ETH address")
	cmd.Flags().String(amountFlag, "", "Withdraw amount")
	cmd.MarkFlagRequired(candidateFlag)
	cmd.MarkFlagRequired(amountFlag)
	return cmd
}
