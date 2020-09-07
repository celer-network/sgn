package validator

import (
	"errors"
	"fmt"

	"github.com/celer-network/sgn/x/validator/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case QuerySyncer:
			return querySyncer(ctx, req, keeper)
		case QueryDelegator:
			return queryDelegator(ctx, req, keeper)
		case QueryCandidate:
			return queryCandidate(ctx, req, keeper)
		case QueryCandidateDelegators:
			return queryCandidateDelegators(ctx, req, keeper)
		case QueryReward:
			return queryReward(ctx, req, keeper)
		case QueryRewardEpoch:
			return queryRewardEpoch(ctx, req, keeper)
		case QueryParameters:
			return queryParameters(ctx, keeper)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Unknown validator query endpoint")
		}
	}
}

func querySyncer(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	syncer := keeper.GetSyncer(ctx)
	res, err := codec.MarshalJSONIndent(keeper.cdc, syncer)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())

	}

	return res, nil
}

func queryDelegator(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params QueryDelegatorParams
	err := ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	delegator, found := keeper.GetDelegator(ctx, params.CandidateAddress, params.DelegatorAddress)
	if !found {
		return nil, fmt.Errorf("Cannot find delegator %s for candidate %s", params.DelegatorAddress, params.CandidateAddress)
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, delegator)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())

	}

	return res, nil
}

func queryCandidate(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params QueryCandidateParams
	err := ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	candidate, found := keeper.GetCandidate(ctx, params.CandidateAddress)
	if !found {
		return nil, errors.New("Cannot find candidate " + params.CandidateAddress)
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, candidate)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())

	}

	return res, nil
}

func queryCandidateDelegators(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params QueryCandidateParams
	err := ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	delegators := keeper.GetAllDelegators(ctx, params.CandidateAddress)

	res, err := codec.MarshalJSONIndent(keeper.cdc, delegators)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())

	}

	return res, nil
}

func queryReward(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params QueryRewardParams
	err := ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	reward, found := keeper.GetReward(ctx, params.EthAddress)
	if !found {
		return nil, errors.New("Reward does not exist for " + params.EthAddress)
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, reward)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())

	}

	return res, nil
}

func queryRewardEpoch(ctx sdk.Context, _ abci.RequestQuery, keeper Keeper) ([]byte, error) {
	epoch := keeper.GetRewardEpoch(ctx)
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
