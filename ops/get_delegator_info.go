package ops

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	delegatorFlag = "delegator"
)

func getDelegatorInfo() error {
	ethClient, err := initEthClient()
	if err != nil {
		return err
	}
	candidate := ethcommon.HexToAddress(viper.GetString(candidateFlag))
	delegator := ethcommon.HexToAddress(viper.GetString(delegatorFlag))
	info, err := ethClient.DPoS.GetDelegatorInfo(&bind.CallOpts{}, candidate, delegator)
	if err != nil {
		return err
	}
	bytes, err := json.MarshalIndent(&info, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Println(string(bytes))
	if err != nil {
		return err
	}
	return nil
}

func GetDelegatorInfoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-delegator-info",
		Short: "Get delegator info",
		RunE: func(cmd *cobra.Command, args []string) error {
			return getDelegatorInfo()
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			err := cmd.MarkFlagRequired(candidateFlag)
			if err != nil {
				return err
			}
			err = cmd.MarkFlagRequired(delegatorFlag)
			if err != nil {
				return err
			}
			err = viper.BindPFlag(candidateFlag, cmd.Flags().Lookup(candidateFlag))
			if err != nil {
				return err
			}
			err = viper.BindPFlag(delegatorFlag, cmd.Flags().Lookup(delegatorFlag))
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().String(candidateFlag, "", "Candidate ETH address")
	cmd.Flags().String(delegatorFlag, "", "Delegator ETH address")
	return cmd
}
