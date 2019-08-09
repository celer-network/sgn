package guardianmanager

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
		case QueryGuardian:
			return queryEthAddress(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown guardianmanager query endpoint")
		}
	}
}

func queryEthAddress(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryGuardianParams
	err := ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	guaridan := keeper.GetGuardian(ctx, params.EthAddress)
	res, err := codec.MarshalJSONIndent(keeper.cdc, guaridan)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}
