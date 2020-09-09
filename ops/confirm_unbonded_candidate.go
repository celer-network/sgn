package ops

import (
	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func confirmUnbondedCandidate() error {
	ethClient, err := common.NewEthClientFromConfig()
	if err != nil {
		return err
	}
	candidate := ethcommon.HexToAddress(viper.GetString(candidateFlag))

	log.Infof(
		"Confirming unbonded candidate %s",
		candidate.Hex(),
	)
	_, err = ethClient.Transactor.TransactWaitMined(
		"ConfirmUnbondedCandidate",
		func(transactor bind.ContractTransactor, opts *bind.TransactOpts) (*ethtypes.Transaction, error) {
			return ethClient.DPoS.ConfirmUnbondedCandidate(opts, candidate)
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func ConfirmUnbondedCandidateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "confirm-unbonded-candidate",
		Short: "Confirm an unbonded candidate",
		RunE: func(cmd *cobra.Command, args []string) error {
			return confirmUnbondedCandidate()
		},
	}
	cmd.Flags().String(candidateFlag, "", "Candidate ETH address")
	cmd.MarkFlagRequired(candidateFlag)
	return cmd
}
