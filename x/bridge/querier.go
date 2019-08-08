package bridge

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryEthAddress:
			return queryEthAddress(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown bridge query endpoint")
		}
	}
}

// nolint: unparam
func queryEthAddress(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryEthAddressParams
	err := ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	ethAddress := keeper.GetEthAddress(ctx, params.Address)
	res, err := codec.MarshalJSONIndent(keeper.cdc, ethAddress)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}
