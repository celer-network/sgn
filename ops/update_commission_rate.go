package ops

import (
	"context"
	"fmt"
	"math/big"

	"github.com/celer-network/goutils/log"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	rateFlag                 = "rate"
	additionalLockPeriodFlag = "additional-lock-period"
)

func UpdateCommissionRate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-commission-rate",
		Short: "Update commission rate",
	}
	cmd.AddCommand(
		announceIncreaseCommissionRate(),
		confirmIncreaseCommissionRate(),
		reduceCommissionRate(),
		extendCommissionRateLockTime(),
	)
	return cmd
}

func announceIncreaseCommissionRate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "announce-increase",
		Short: "Announce increase commission rate",
		RunE: func(cmd *cobra.Command, args []string) error {
			ethClient, err := initEthClient()
			if err != nil {
				return err
			}
			newRate := new(big.Int)
			newRate.SetString(viper.GetString(rateFlag), 10)
			rateLockPeriod := new(big.Int)
			rateLockPeriod.SetString(viper.GetString(additionalLockPeriodFlag), 10)
			info, err := ethClient.DPoS.GetCandidateInfo(&bind.CallOpts{}, ethClient.Address)
			if err != nil {
				return err
			}
			if !info.Initialized {
				return fmt.Errorf("candidate not initialized")
			}
			if newRate.Cmp(info.CommissionRate) <= 0 {
				return fmt.Errorf("commission rate not increasing")
			}
			newLockEndTime := new(big.Int).Add(info.RateLockEndTime, rateLockPeriod)
			header, err := ethClient.Client.HeaderByNumber(context.Background(), nil)
			if err != nil {
				return err
			}
			if header.Number.Cmp(info.RateLockEndTime) > 0 {
				newLockEndTime = new(big.Int).Add(header.Number, rateLockPeriod)
			}
			log.Infof("Announce increase commission rate: %s, rateLockEndTime: %s", newRate, newLockEndTime)
			receipt, err := ethClient.Transactor.TransactWaitMined(
				"AnnounceIncreaseCommissionRate",
				func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*ethtypes.Transaction, error) {
					return ethClient.DPoS.AnnounceIncreaseCommissionRate(opts, newRate, newLockEndTime)
				},
			)
			if err != nil {
				return err
			}
			log.Infof("Transaction %x succeeded", receipt.TxHash)
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			err := cmd.MarkFlagRequired(rateFlag)
			if err != nil {
				return err
			}
			err = viper.BindPFlag(rateFlag, cmd.Flags().Lookup(rateFlag))
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().String(rateFlag, "", "Commission rate in unit of 0.01% (e.g., 120 is 1.2%)")
	cmd.Flags().String(additionalLockPeriodFlag, "", "(optional) additional rate lock period")
	return cmd
}

func confirmIncreaseCommissionRate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "confirm-increase",
		Short: "Confirm increase commission rate",
		RunE: func(cmd *cobra.Command, args []string) error {
			ethClient, err := initEthClient()
			if err != nil {
				return err
			}
			log.Infof("Confirm increase commission rate")
			receipt, err := ethClient.Transactor.TransactWaitMined(
				"ConfirmIncreaseCommissionRate",
				func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*ethtypes.Transaction, error) {
					return ethClient.DPoS.ConfirmIncreaseCommissionRate(opts)
				},
			)
			if err != nil {
				return err
			}
			log.Infof("Transaction %x succeeded", receipt.TxHash)
			return nil
		},
	}
	return cmd
}

func reduceCommissionRate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reduce",
		Short: "Reduce commission rate",
		RunE: func(cmd *cobra.Command, args []string) error {
			ethClient, err := initEthClient()
			if err != nil {
				return err
			}
			newRate := new(big.Int)
			newRate.SetString(viper.GetString(rateFlag), 10)
			rateLockPeriod := new(big.Int)
			rateLockPeriod.SetString(viper.GetString(additionalLockPeriodFlag), 10)
			info, err := ethClient.DPoS.GetCandidateInfo(&bind.CallOpts{}, ethClient.Address)
			if err != nil {
				return err
			}
			if !info.Initialized {
				return fmt.Errorf("candidate not initialized")
			}
			if newRate.Cmp(info.CommissionRate) >= 0 {
				return fmt.Errorf("commission rate not decreasing")
			}
			newLockEndTime := new(big.Int).Add(info.RateLockEndTime, rateLockPeriod)
			header, err := ethClient.Client.HeaderByNumber(context.Background(), nil)
			if err != nil {
				return err
			}
			if header.Number.Cmp(info.RateLockEndTime) > 0 {
				newLockEndTime = new(big.Int).Add(header.Number, rateLockPeriod)
			}
			log.Infof("Reduce commission rate: %s, rateLockEndTime: %s", newRate, newLockEndTime)
			receipt, err := ethClient.Transactor.TransactWaitMined(
				"NonIncreaseCommissionRate",
				func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*ethtypes.Transaction, error) {
					return ethClient.DPoS.NonIncreaseCommissionRate(opts, newRate, newLockEndTime)
				},
			)
			if err != nil {
				return err
			}
			log.Infof("Transaction %x succeeded", receipt.TxHash)
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			err := cmd.MarkFlagRequired(rateFlag)
			if err != nil {
				return err
			}
			err = viper.BindPFlag(rateFlag, cmd.Flags().Lookup(rateFlag))
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().String(rateFlag, "", "Commission rate in unit of 0.01% (e.g., 120 is 1.2%)")
	cmd.Flags().String(additionalLockPeriodFlag, "", "(optional) additional rate lock period")
	return cmd
}

func extendCommissionRateLockTime() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extend-lock-period",
		Short: "Extend commission rate lock end time",
		RunE: func(cmd *cobra.Command, args []string) error {
			ethClient, err := initEthClient()
			if err != nil {
				return err
			}
			rateLockPeriod := new(big.Int)
			rateLockPeriod.SetString(viper.GetString(additionalLockPeriodFlag), 10)
			info, err := ethClient.DPoS.GetCandidateInfo(&bind.CallOpts{}, ethClient.Address)
			if err != nil {
				return err
			}
			if !info.Initialized {
				return fmt.Errorf("candidate not initialized")
			}
			newLockEndTime := new(big.Int).Add(info.RateLockEndTime, rateLockPeriod)
			header, err := ethClient.Client.HeaderByNumber(context.Background(), nil)
			if err != nil {
				return err
			}
			if header.Number.Cmp(info.RateLockEndTime) > 0 {
				newLockEndTime = new(big.Int).Add(header.Number, rateLockPeriod)
			}
			log.Infof("Extend commission rateLockEndTime: %s", newLockEndTime)
			receipt, err := ethClient.Transactor.TransactWaitMined(
				"NonIncreaseCommissionRate",
				func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*ethtypes.Transaction, error) {
					return ethClient.DPoS.NonIncreaseCommissionRate(opts, info.CommissionRate, newLockEndTime)
				},
			)
			if err != nil {
				return err
			}
			log.Infof("Transaction %x succeeded", receipt.TxHash)
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			err := cmd.MarkFlagRequired(additionalLockPeriodFlag)
			if err != nil {
				return err
			}
			err = viper.BindPFlag(additionalLockPeriodFlag, cmd.Flags().Lookup(additionalLockPeriodFlag))
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().String(additionalLockPeriodFlag, "", "additional rate lock period")
	return cmd
}
