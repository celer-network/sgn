package validator

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryPuller:
			return queryPuller(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("Unknown validator query endpoint")
		}
	}
}

func queryPuller(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	puller := keeper.GetPuller(ctx)
	res, err := codec.MarshalJSONIndent(keeper.cdc, puller)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("Could not marshal result to JSON", err.Error()))
	}

	return res, nil
}
