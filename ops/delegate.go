package ops

import (
	"math/big"

	"github.com/celer-network/goutils/eth"
	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	stakeFlag     = "stake"
	candidateFlag = "candidate"
)

func delegate() error {
	ethClient, err := initEthClient()
	if err != nil {
		log.Error(err)
		return err
	}
	stake := new(big.Int)
	stake.SetString(viper.GetString(stakeFlag), 10)
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
	if allowance.Cmp(stake) < 0 {
		log.Info("Approving CELR to DPoS contract")
		_, approveErr := ethClient.Transactor.TransactWaitMined(
			"Approve",
			func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*ethtypes.Transaction, error) {
				return celrContract.Approve(opts, dPoSAddress, stake)
			},
		)
		if approveErr != nil {
			return approveErr
		}
	}

	log.Infof(
		"Sending delegate transaction with stake: %s, candidate: %s ",
		stake,
		candidate.Hex(),
	)
	tx, err := ethClient.Transactor.Transact(
		&eth.TransactionStateHandler{
			OnMined: func(receipt *ethtypes.Receipt) {
				if receipt.Status == ethtypes.ReceiptStatusSuccessful {
					log.Infof("Delegate transaction %x succeeded", receipt.TxHash)
				} else {
					log.Errorf("Delegate transaction %x failed", receipt.TxHash)
				}
			},
			OnError: func(tx *ethtypes.Transaction, err error) {
				log.Errorf("Delegate transaction %x err: %s", tx.Hash(), err)
			},
		},
		func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*ethtypes.Transaction, error) {
			return ethClient.DPoS.Delegate(opts, candidate, stake)
		},
	)
	if err != nil {
		return err
	}
	log.Infof("Delegate transaction: %x", tx.Hash())
	return nil
}

func DelegateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegate",
		Short: "Delegate stake to a validator candidate",
		RunE: func(cmd *cobra.Command, args []string) error {
			return delegate()
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			err := cmd.MarkFlagRequired(stakeFlag)
			if err != nil {
				return err
			}
			err = cmd.MarkFlagRequired(candidateFlag)
			if err != nil {
				return err
			}
			err = viper.BindPFlag(stakeFlag, cmd.Flags().Lookup(stakeFlag))
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
	cmd.Flags().String(stakeFlag, "", "stake")
	cmd.Flags().String(candidateFlag, "", "candidate ETH address")
	return cmd
}
