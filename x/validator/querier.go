package validator

import (
	"fmt"

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
		case QueryPusher:
			return queryPusher(ctx, req, keeper)
		case QueryDelegator:
			return queryDelegator(ctx, req, keeper)
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

func queryPusher(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	pusher := keeper.GetPusher(ctx)
	res, err := codec.MarshalJSONIndent(keeper.cdc, pusher)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("Could not marshal result to JSON", err.Error()))
	}

	return res, nil
}

func queryDelegator(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryDelegatorParams
	err := ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("Failed to parse params: %s", err))
	}

	delegator := keeper.GetDelegator(ctx, params.CandidateAddress, params.DelegatorAddress)
	res, err := codec.MarshalJSONIndent(keeper.cdc, delegator)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("Could not marshal result to JSON", err.Error()))

	}

	return res, nil
}

func queryCandidate(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryCandidateParams
	err := ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("Failed to parse params: %s", err))
	}

	candidate := keeper.GetCandidate(ctx, params.CandidateAddress)
	res, err := codec.MarshalJSONIndent(keeper.cdc, candidate)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("Could not marshal result to JSON", err.Error()))

	}

	return res, nil
}
