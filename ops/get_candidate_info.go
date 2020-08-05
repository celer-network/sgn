package ops

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func getCandidateInfo() error {
	ethClient, err := initEthClient()
	if err != nil {
		return err
	}
	candidate := ethcommon.HexToAddress(viper.GetString(candidateFlag))
	info, err := ethClient.DPoS.GetCandidateInfo(&bind.CallOpts{}, candidate)
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

func GetCandidateInfoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-candidate-info",
		Short: "Get candidate info",
		RunE: func(cmd *cobra.Command, args []string) error {
			return getCandidateInfo()
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			err := cmd.MarkFlagRequired(candidateFlag)
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
	cmd.Flags().String(candidateFlag, "", "Candidate ETH address")
	return cmd
}
