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
			return queryGuardian(ctx, req, keeper)
		case QueryRequest:
			return queryRequest(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("Unknown guardianmanager query endpoint")
		}
	}
}

func queryGuardian(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryGuardianParams
	err := ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("Failed to parse params: %s", err))
	}

	guaridan := keeper.GetGuardian(ctx, params.EthAddress)
	res, err := codec.MarshalJSONIndent(keeper.cdc, guaridan)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("Could not marshal result to JSON", err.Error()))

	}

	return res, nil
}

func queryRequest(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryRequestParams
	err := ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("Failed to parse params: %s", err))
	}

	request, found := keeper.GetRequest(ctx, params.ChannelId)
	if !found {
		return nil, sdk.ErrInternal("Could not find corresponding request")
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, request)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("Could not marshal result to JSON", err.Error()))

	}

	return res, nil
}
