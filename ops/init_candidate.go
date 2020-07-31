package ops

import (
	"context"
	"math/big"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	minSelfStakeFlag   = "min-self-stake"
	commissionRateFlag = "commission-rate"
	rateLockPeriodFlag = "rate-lock-period"
)

func initCandidate() error {
	ethClient, err := initEthClient()
	if err != nil {
		return err
	}
	minSelfStake := new(big.Int)
	minSelfStake.SetString(viper.GetString(minSelfStakeFlag), 10)
	commissionRate := new(big.Int)
	commissionRate.SetString(viper.GetString(commissionRateFlag), 10)
	rateLockPeriod := new(big.Int)
	rateLockPeriod.SetString(viper.GetString(rateLockPeriodFlag), 10)

	dposContract := ethClient.DPoS
	info, err := dposContract.GetCandidateInfo(&bind.CallOpts{}, ethClient.Address)
	if err != nil {
		return err
	}
	if !info.Initialized {
		header, headerErr := ethClient.Client.HeaderByNumber(context.Background(), nil)
		if headerErr != nil {
			return headerErr
		}
		rateLockEndTime := new(big.Int).Add(header.Number, rateLockPeriod)
		log.Infof(
			"Sending initialize candidate transaction with minSelfStake: %s, commissionRate: %s, rateLockEndTime: %s",
			minSelfStake,
			commissionRate,
			rateLockEndTime,
		)
		receipt, initCandidateErr := ethClient.Transactor.TransactWaitMined(
			"InitializeCandidate",
			func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*ethtypes.Transaction, error) {
				return dposContract.InitializeCandidate(opts, minSelfStake, commissionRate, rateLockEndTime)
			},
		)
		if initCandidateErr != nil {
			return initCandidateErr
		}
		log.Infof("Initialize candidate transaction %x succeeded", receipt.TxHash)
	}

	operatorAddress, err := sdk.AccAddressFromBech32(viper.GetString(common.FlagSgnOperator))
	if err != nil {
		return err
	}
	log.Infof("Calling updateSidechainAddr for %s", operatorAddress)
	receipt, err := ethClient.Transactor.TransactWaitMined(
		"UpdateSidechainAddr",
		func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*ethtypes.Transaction, error) {
			return ethClient.SGN.UpdateSidechainAddr(opts, operatorAddress.Bytes())
		},
	)
	if err != nil {
		return err
	}
	log.Infof("Update sidechain address transaction %x succeeded", receipt.TxHash)
	return nil
}

func InitCandidateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init-candidate",
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
			err = cmd.MarkFlagRequired(rateLockPeriodFlag)
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
			err = viper.BindPFlag(rateLockPeriodFlag, cmd.Flags().Lookup(rateLockPeriodFlag))
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().String(minSelfStakeFlag, "", "Minimum self stake")
	cmd.Flags().String(commissionRateFlag, "", "Commission rate")
	cmd.Flags().String(rateLockPeriodFlag, "", "Rate lock period")
	return cmd
}
