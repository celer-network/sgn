package common

import (
	"bufio"
	"os"
	"time"

	"github.com/celer-network/sgn/proto/chain"
	"github.com/celer-network/sgn/proto/entity"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/input"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/protobuf/proto"
	"github.com/spf13/viper"
)

const (
	retryTimeout = 500 * time.Millisecond
)

func RobustQuery(cliCtx context.CLIContext, route string) ([]byte, error) {
	res, _, err := cliCtx.Query(route)
	if err != nil {
		time.Sleep(retryTimeout)
		res, _, err = cliCtx.Query(route)
		return res, err
	}

	return res, err
}

func RobustQueryWithData(cliCtx context.CLIContext, route string, data []byte) ([]byte, error) {
	res, _, err := cliCtx.QueryWithData(route, data)
	if err != nil {
		time.Sleep(retryTimeout)
		res, _, err = cliCtx.QueryWithData(route, data)
		return res, err
	}

	return res, err
}

func UnmarshalSignedSimplexStateBytes(input []byte) (*chain.SignedSimplexState, *entity.SimplexPaymentChannel, error) {
	var signedSimplexState chain.SignedSimplexState
	err := proto.Unmarshal(input, &signedSimplexState)
	if err != nil {
		return nil, nil, err
	}
	var simplexChannel entity.SimplexPaymentChannel
	err = proto.Unmarshal(signedSimplexState.SimplexState, &simplexChannel)
	if err != nil {
		return nil, nil, err
	}
	return &signedSimplexState, &simplexChannel, nil
}

func ParseTransactorAddrs(ts []string) ([]sdk.AccAddress, error) {
	var transactors []sdk.AccAddress
	for _, t := range ts {
		transactor, err := sdk.AccAddressFromBech32(t)
		if err != nil {
			return transactors, err
		}

		transactors = append(transactors, transactor)
	}

	return transactors, nil
}

func SetupUserPassword() error {
	buf := bufio.NewReader(os.Stdin)

	if viper.Get(FlagEthPassphrase) == nil {
		pass, err := input.GetString("Enter eth keystore passphrase:", buf)
		if err != nil {
			return err
		}

		viper.Set(FlagEthPassphrase, pass)
	}

	if viper.Get(FlagSgnPassphrase) == nil {
		pass, err := input.GetString("Enter sidechain validator passphrase:", buf)
		if err != nil {
			return err
		}

		viper.Set(FlagSgnPassphrase, pass)
	}

	return nil
}
