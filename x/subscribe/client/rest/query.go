package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/celer-network/sgn/x/subscribe/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/subscribe/subscription/{ethAddr}",
		subscriptionHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		"/subscribe/request/{channelId}",
		requestHandlerFn(cliCtx),
	).Methods("GET")
}

// http request handler to query subscription
func subscriptionHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		vars := mux.Vars(r)
		ethAddr := vars["ethAddr"]
		params := types.NewQuerySubscriptionParams(ethAddr)
		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.ModuleName, types.QuerySubscription)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// http request handler to query request
func requestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		vars := mux.Vars(r)
		channelId := ethcommon.Hex2Bytes(vars["channelId"])
		params := types.NewQueryRequestParams(channelId)
		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.ModuleName, types.QuerySubscription)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
