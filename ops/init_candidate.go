package ops

import (
	"math/big"

	"github.com/celer-network/goutils/eth"
	"github.com/celer-network/goutils/log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	minSelfStakeFlag    = "minselfstake"
	commissionRateFlag  = "commissionrate"
	rateLockEndTimeFlag = "ratelockendtime"
)

func initCandidate() error {
	ethClient, err := initEthClient()
	minSelfStake := new(big.Int)
	minSelfStake.SetString(viper.GetString(minSelfStakeFlag), 10)
	commissionRate := new(big.Int)
	commissionRate.SetString(viper.GetString(commissionRateFlag), 10)
	rateLockEndTime := new(big.Int)
	rateLockEndTime.SetString(viper.GetString(rateLockEndTimeFlag), 10)

	log.Infof(
		"Sending initialize candidate transaction with minSelfStake: %s, commissionRate: %s, rateLockEndTime: %s",
		minSelfStake,
		commissionRate,
		rateLockEndTime,
	)
	tx, err := ethClient.Transactor.Transact(
		&eth.TransactionStateHandler{
			OnMined: func(receipt *ethtypes.Receipt) {
				if receipt.Status == ethtypes.ReceiptStatusSuccessful {
					log.Infof("Initialize candidate transaction %x succeeded", receipt.TxHash)
				} else {
					log.Errorf("Initialize candidate transaction %x failed", receipt.TxHash)
				}
			},
			OnError: func(tx *ethtypes.Transaction, err error) {
				log.Errorf("Initialize candidate transaction %x err: %s", tx.Hash(), err)
			},
		},
		func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*ethtypes.Transaction, error) {
			return ethClient.DPoS.InitializeCandidate(opts, minSelfStake, commissionRate, rateLockEndTime)
		},
	)
	if err != nil {
		return err
	}
	log.Infof("Initialize candidate transaction: %x", tx.Hash())
	return nil
}

func InitCandidateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "initcandidate",
		Short: "Initialize a validator candidate",
		RunE: func(cmd *cobra.Command, args []string) error {
			return initCandidate()
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			err := cmd.MarkFlagRequired(minSelfStakeFlag)
			if err != nil {
				return err
			}
			err = cmd.MarkFlagRequired(commissionRateFlag)
			if err != nil {
				return err
			}
			err = cmd.MarkFlagRequired(rateLockEndTimeFlag)
			if err != nil {
				return err
			}
			err = viper.BindPFlag(minSelfStakeFlag, cmd.Flags().Lookup(minSelfStakeFlag))
			if err != nil {
				return err
			}
			err = viper.BindPFlag(commissionRateFlag, cmd.Flags().Lookup(commissionRateFlag))
			if err != nil {
				return err
			}
			err = viper.BindPFlag(rateLockEndTimeFlag, cmd.Flags().Lookup(rateLockEndTimeFlag))
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().String(minSelfStakeFlag, "", "Minimum self stake")
	cmd.Flags().String(commissionRateFlag, "", "Commission rate")
	cmd.Flags().String(rateLockEndTimeFlag, "", "Rate lock end time")
	return cmd
}
