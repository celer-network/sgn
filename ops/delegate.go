package ops

import (
	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func delegate() error {
	ethClient, err := common.NewEthClientFromConfig()
	if err != nil {
		return err
	}
	amount := calcRawAmount(viper.GetString(amountFlag))
	candidate := mainchain.Hex2Addr(viper.GetString(candidateFlag))

	err = approveCelrToDPoS(ethClient, amount)
	if err != nil {
		return err
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
	cmd.Flags().String(amountFlag, "", "Stake amount (integer in unit of CELR)")
	cmd.Flags().String(candidateFlag, "", "Candidate ETH address")
	cmd.MarkFlagRequired(amountFlag)
	cmd.MarkFlagRequired(candidateFlag)
	return cmd
}
