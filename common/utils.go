package common

import (
	"math/big"
	"time"

	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/proto/chain"
	"github.com/celer-network/sgn/proto/entity"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/golang/protobuf/proto"
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

func NewCommission(ethClient *mainchain.EthClient, commissionRate *big.Int) (staking.Commission, error) {
	commissionBase, err := ethClient.DPoS.COMMISSIONRATEBASE(&bind.CallOpts{})
	if err != nil {
		return staking.Commission{}, err
	}

	prec := int64(len(commissionBase.String()) - 1)
	return staking.Commission{
		CommissionRates: staking.CommissionRates{
			Rate:          sdk.NewDecFromBigIntWithPrec(commissionRate, prec),
			MaxRate:       sdk.NewDec(1),
			MaxChangeRate: sdk.NewDec(1),
		},
	}, nil
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
