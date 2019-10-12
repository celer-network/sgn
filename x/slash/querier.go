package slash

import (
	"fmt"

	"github.com/celer-network/sgn/x/slash/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryPenalty:
			return queryPenalty(ctx, req, keeper)
		case QueryParameters:
			return queryParameters(ctx, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("Unknown slash query endpoint")
		}
	}
}

func queryPenalty(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryPenaltyParams
	err := ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	penalty, found := keeper.GetPenalty(ctx, params.Nonce)
	if !found {
		return nil, sdk.ErrInternal("Penalty does not exist")
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, penalty)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("Could not marshal result to JSON", err.Error()))
	}

	return res, nil
}

func queryParameters(ctx sdk.Context, k Keeper) ([]byte, sdk.Error) {
	params := k.GetParams(ctx)

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, params)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return res, nil
}
