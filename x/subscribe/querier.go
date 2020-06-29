package subscribe

import (
	"errors"
	"fmt"

	"github.com/celer-network/sgn/x/subscribe/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case QuerySubscription:
			return querySubscription(ctx, req, keeper)
		case QueryRequest:
			return queryRequest(ctx, req, keeper)
		case QueryEpoch:
			return queryEpoch(ctx, req, keeper)
		case QueryParameters:
			return queryParameters(ctx, keeper)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Unknown subscribe query endpoint")
		}
	}
}

func querySubscription(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params QuerySubscriptionParams
	err := ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to parse params: %s", err))
	}

	subscription, found := keeper.GetSubscription(ctx, params.EthAddress)
	if !found {
		return nil, errors.New(fmt.Sprintf("cannot find subscription for %s", params.EthAddress))
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, subscription)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())

	}

	return res, nil
}

func queryRequest(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params QueryRequestParams
	err := ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	request, found := keeper.GetRequest(ctx, params.ChannelId, params.Receiver)
	if !found {
		return nil, errors.New("Could not find corresponding request")
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, request)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())

	}

	return res, nil
}

func queryEpoch(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params QueryEpochParams
	err := ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	var epoch Epoch

	// Get the latest epoch if user does not specify the epoch id
	if params.EpochId <= 0 {
		epoch = keeper.GetLatestEpoch(ctx)
	} else {
		epoch, _ = keeper.GetEpoch(ctx, sdk.NewInt(params.EpochId))
		if epoch.Id.IsZero() {
			return nil, errors.New("Could not find corresponding epoch")
		}
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, epoch)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryParameters(ctx sdk.Context, k Keeper) ([]byte, error) {
	params := k.GetParams(ctx)

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}
