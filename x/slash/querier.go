package slash

import (
	"errors"
	"fmt"

	"github.com/celer-network/sgn/x/slash/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case QueryPenalty:
			return queryPenalty(ctx, req, keeper)
		case QueryPenalties:
			return queryPenalties(ctx, req, keeper)
		case QueryParameters:
			return queryParameters(ctx, keeper)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Unknown slash query endpoint")
		}
	}
}

func queryPenalty(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params QueryPenaltyParams
	err := ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to parse params: %s", err))
	}

	penalty, found := keeper.GetPenalty(ctx, params.Nonce)
	if !found {
		return nil, errors.New("Penalty does not exist")
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, penalty)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())

	}

	return res, nil
}

func queryPenalties(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	penalties := keeper.GetPenalties(ctx)
	res, err := codec.MarshalJSONIndent(keeper.cdc, penalties)
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
