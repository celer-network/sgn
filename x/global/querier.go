package global

import (
	"fmt"

	"github.com/celer-network/sgn/x/global/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryLatestBlock:
			return queryLatestBlock(ctx, req, keeper)
		case QueryEpoch:
			return queryEpoch(ctx, req, keeper)
		case QueryParameters:
			return queryParameters(ctx, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("Unknown global query endpoint")
		}
	}
}

func queryLatestBlock(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	latestBlock := keeper.GetLatestBlock(ctx)
	res, err := codec.MarshalJSONIndent(keeper.cdc, latestBlock)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("Could not marshal result to JSON", err.Error()))

	}

	return res, nil
}

func queryEpoch(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryEpochParams
	err := ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("Failed to parse params: %s", err))
	}

	var epoch Epoch

	if params.EpochId <= 0 {
		epoch = keeper.GetLatestEpoch(ctx)
	} else {
		epoch, _ = keeper.GetEpoch(ctx, sdk.NewInt(params.EpochId))
		if epoch.Id.IsZero() {
			return nil, sdk.ErrInternal("Could not find corresponding epoch")
		}
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, epoch)
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
