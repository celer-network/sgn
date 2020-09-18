package ops

import (
	"math/big"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/viper"
)

const (
	candidateFlag = "candidate"
	amountFlag    = "amount"
)

func calcRawAmount(amount string) *big.Int {
	rawAmount := new(big.Int)
	rawAmount.SetString(amount, 10)
	return new(big.Int).Mul(rawAmount, big.NewInt(common.TokenDec))
}

func approveCelrToDPoS(ethClient *mainchain.EthClient, amount *big.Int) error {
	dPoSAddress := ethClient.DPoSAddress
	celrContract, err := mainchain.NewERC20(
		mainchain.Hex2Addr(viper.GetString(common.FlagEthCelrAddress)),
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
	return nil
}
