package ops

import (
	"github.com/celer-network/goutils/log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func confirmWithdraw() error {
	ethClient, err := initEthClient()
	candidate := ethcommon.HexToAddress(viper.GetString(candidateFlag))

	log.Infof(
		"Confirming withdrawal intents for the stake delegated to candidate %s",
		candidate.Hex(),
	)
	receipt, err := ethClient.Transactor.TransactWaitMined(
		"ConfirmWithdraw",
		func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*ethtypes.Transaction, error) {
			return ethClient.DPoS.ConfirmWithdraw(opts, candidate)
		},
	)
	if err != nil {
		return err
	}
	log.Infof("Transaction %x succeeded", receipt.TxHash)
	return nil
}

func ConfirmWithdrawCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "confirm-withdraw",
		Short: "Confirm withdrawal intents for the stake delegated to a candidate",
		RunE: func(cmd *cobra.Command, args []string) error {
			return confirmWithdraw()
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			err := cmd.MarkFlagRequired(candidateFlag)
			if err != nil {
				return err
			}
			err = viper.BindPFlag(candidateFlag, cmd.Flags().Lookup(candidateFlag))
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().String(candidateFlag, "", "Candidate ETH address")
	return cmd
}
