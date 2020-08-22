package ops

import (
	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func delegate() error {
	ethClient, err := initEthClient()
	if err != nil {
		log.Error(err)
		return err
	}
	amount := calcRawAmount(viper.GetString(amountFlag))
	candidate := ethcommon.HexToAddress(viper.GetString(candidateFlag))

	dPoSAddress := ethClient.DPoSAddress
	celrContract, err := mainchain.NewERC20(
		ethcommon.HexToAddress(viper.GetString(common.FlagEthCelrAddress)),
		ethClient.Client,
	)
	if err != nil {
		return err
	}
	allowance, err :=
		celrContract.Allowance(&bind.CallOpts{}, ethClient.Transactor.Address(), dPoSAddress)
	if err != nil {
		return err
	}
	if allowance.Cmp(amount) < 0 {
		log.Info("Approving CELR to DPoS contract")
		_, approveErr := ethClient.Transactor.TransactWaitMined(
			"Approve",
			func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*ethtypes.Transaction, error) {
				return celrContract.Approve(opts, dPoSAddress, amount)
			},
		)
		if approveErr != nil {
			return approveErr
		}
	}

	log.Infof(
		"Sending delegate transaction with amount: %s, candidate: %s",
		amount,
		candidate.Hex(),
	)
	_, err = ethClient.Transactor.TransactWaitMined(
		"Delegate",
		func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*ethtypes.Transaction, error) {
			return ethClient.DPoS.Delegate(opts, candidate, amount)
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func DelegateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegate",
		Short: "Delegate stake to a validator candidate",
		RunE: func(cmd *cobra.Command, args []string) error {
			return delegate()
		},
	}
	cmd.Flags().String(amountFlag, "", "Stake amount")
	cmd.Flags().String(candidateFlag, "", "Candidate ETH address")
	cmd.MarkFlagRequired(amountFlag)
	cmd.MarkFlagRequired(candidateFlag)
	return cmd
}
