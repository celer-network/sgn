package ops

import (
	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func updateMinSelfStake() error {
	ethClient, err := common.NewEthClientFromConfig()
	if err != nil {
		return err
	}
	amount := calcRawAmount(viper.GetString(amountFlag))

	log.Infof("Update minimal self-delegated stake to %s", amount)
	_, err = ethClient.Transactor.TransactWaitMined(
		"UpdateMinSelfStake",
		func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*ethtypes.Transaction, error) {
			return ethClient.DPoS.UpdateMinSelfStake(opts, amount)
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func UpdateMinSelfStakeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-min-self-stake",
		Short: "Update minimal self-delegated stake",
		RunE: func(cmd *cobra.Command, args []string) error {
			return updateMinSelfStake()
		},
	}
	cmd.Flags().String(amountFlag, "", "Minimal self-delegated stake")
	cmd.MarkFlagRequired(amountFlag)
	return cmd
}
