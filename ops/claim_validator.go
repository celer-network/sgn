package ops

import (
	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func claimValidator() error {
	ethClient, err := initEthClient()
	if err != nil {
		return err
	}

	operatorAddress, err := sdk.AccAddressFromBech32(viper.GetString(common.FlagSgnOperator))
	if err != nil {
		return err
	}
	log.Infof("Calling claimValidator for %s", operatorAddress)
	_, err = ethClient.Transactor.TransactWaitMined(
		"ClaimValidator",
		func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*ethtypes.Transaction, error) {
			return ethClient.DPoS.ClaimValidator(opts)
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func ClaimValidatorCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim-validator",
		Short: "Claim the validator status",
		RunE: func(cmd *cobra.Command, args []string) error {
			return claimValidator()
		},
	}
	return cmd
}
