package ops

import (
	"encoding/json"
	"fmt"

	"github.com/celer-network/sgn/common"
	"github.com/spf13/cobra"
)

func getAllValidators() error {
	ethClient, err := common.NewEthClientFromConfig()
	if err != nil {
		return err
	}
	validators, totalStakes, _, err := ethClient.GetValidators()
	bytes, err := json.MarshalIndent(&validators, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Println(string(bytes))
	fmt.Println("total stakes:", totalStakes)
	return err
}

func GetAllValidatorsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-all-validators",
		Short: "Get all validators info",
		RunE: func(cmd *cobra.Command, args []string) error {
			return getAllValidators()
		},
	}
	return cmd
}
