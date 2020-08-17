package ops

import (
	"github.com/celer-network/goutils/log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func updateMinSelfStake() error {
	ethClient, err := initEthClient()
	amount := calcRawAmount(viper.GetString(amountFlag))

	log.Infof("Update minimal self-delegated stake to %s", amount)
	receipt, err := ethClient.Transactor.TransactWaitMined(
		"UpdateMinSelfStake",
		func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*ethtypes.Transaction, error) {
			return ethClient.DPoS.UpdateMinSelfStake(opts, amount)
		},
	)
	if err != nil {
		return err
	}
	log.Infof("UpdateMinSelfStake transaction %x succeeded", receipt.TxHash)
	return nil
}

func UpdateMinSelfStakeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-min-self-stake",
		Short: "Update minimal self-delegated stake",
		RunE: func(cmd *cobra.Command, args []string) error {
			return updateMinSelfStake()
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			err := cmd.MarkFlagRequired(amountFlag)
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
	cmd.Flags().String(amountFlag, "", "Minimal self-delegated stake")
	return cmd
}
