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
	rateFlag    = "rate"
	addLockTime = "add-lock-time"
)

func UpdateCommissionRate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-commission-rate",
		Short: "Update commission rate",
	}
	cmd.AddCommand(
		announceIncreaseCommissionRate(),
		confirmIncreaseCommissionRate(),
		decreaseCommissionRate(),
		extendCommissionRateLockTime(),
	)
	return cmd
}

func announceIncreaseCommissionRate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "announce-increase-rate",
		Short: "Announce increase commission rate",
		RunE: func(cmd *cobra.Command, args []string) error {
			ethClient, err := initEthClient()
			if err != nil {
				return err
			}
			newRate := new(big.Int)
			newRate.SetString(viper.GetString(rateFlag), 10)
			rateLockPeriod := new(big.Int)
			rateLockPeriod.SetString(viper.GetString(addLockTime), 10)
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
			log.Infof("Announce increase commission rate to %s (in unit of 0.01%%), rateLockEndTime: %s", newRate, newLockEndTime)
			_, err = ethClient.Transactor.TransactWaitMined(
				"AnnounceIncreaseCommissionRate",
				func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*ethtypes.Transaction, error) {
					return ethClient.DPoS.AnnounceIncreaseCommissionRate(opts, newRate, newLockEndTime)
				},
			)
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().String(rateFlag, "", "Commission rate in unit of 0.01% (e.g., 120 is 1.2%)")
	cmd.Flags().String(addLockTime, "", "(optional) additional rate lock period")
	cmd.MarkFlagRequired(rateFlag)
	return cmd
}

func confirmIncreaseCommissionRate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "confirm-increase-rate",
		Short: "Confirm increase commission rate",
		RunE: func(cmd *cobra.Command, args []string) error {
			ethClient, err := initEthClient()
			if err != nil {
				return err
			}
			log.Infof("Confirm increase commission rate")
			_, err = ethClient.Transactor.TransactWaitMined(
				"ConfirmIncreaseCommissionRate",
				func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*ethtypes.Transaction, error) {
					return ethClient.DPoS.ConfirmIncreaseCommissionRate(opts)
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

func decreaseCommissionRate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "decrease-rate",
		Short: "Decrease commission rate",
		RunE: func(cmd *cobra.Command, args []string) error {
			ethClient, err := initEthClient()
			if err != nil {
				return err
			}
			newRate := new(big.Int)
			newRate.SetString(viper.GetString(rateFlag), 10)
			rateLockPeriod := new(big.Int)
			rateLockPeriod.SetString(viper.GetString(addLockTime), 10)
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
			log.Infof("Decrease commission rate to %s (in unit of 0.01%%), rateLockEndTime: %s", newRate, newLockEndTime)
			_, err = ethClient.Transactor.TransactWaitMined(
				"NonIncreaseCommissionRate",
				func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*ethtypes.Transaction, error) {
					return ethClient.DPoS.NonIncreaseCommissionRate(opts, newRate, newLockEndTime)
				},
			)
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().String(rateFlag, "", "Commission rate in unit of 0.01% (e.g., 120 is 1.2%)")
	cmd.Flags().String(addLockTime, "", "(optional) additional rate lock period")
	cmd.MarkFlagRequired(rateFlag)
	return cmd
}

func extendCommissionRateLockTime() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extend-lock",
		Short: "Extend commission rate lock end time",
		RunE: func(cmd *cobra.Command, args []string) error {
			ethClient, err := initEthClient()
			if err != nil {
				return err
			}
			rateLockPeriod := new(big.Int)
			rateLockPeriod.SetString(viper.GetString(addLockTime), 10)
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
			_, err = ethClient.Transactor.TransactWaitMined(
				"NonIncreaseCommissionRate",
				func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*ethtypes.Transaction, error) {
					return ethClient.DPoS.NonIncreaseCommissionRate(opts, info.CommissionRate, newLockEndTime)
				},
			)
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().String(addLockTime, "", "additional rate lock period")
	cmd.MarkFlagRequired(addLockTime)
	return cmd
}
