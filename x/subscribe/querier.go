package subscribe

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
		case QuerySubscrption:
			return querySubscription(ctx, req, keeper)
		case QueryRequest:
			return queryRequest(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("Unknown subscribe query endpoint")
		}
	}
}

func querySubscription(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QuerySubscrptionParams
	err := ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	subscription, found := keeper.GetSubscription(ctx, params.EthAddress)
	if !found {
		return nil, sdk.ErrInternal("cannot find subscription")
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, subscription)
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
