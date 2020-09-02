package ops

import (
	"encoding/json"
	"fmt"

	"github.com/celer-network/sgn/common"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func getCandidateInfo() error {
	ethClient, err := common.NewEthClientFromConfig()
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
	}
	cmd.Flags().String(candidateFlag, "", "Candidate ETH address")
	cmd.MarkFlagRequired(candidateFlag)
	return cmd
}
