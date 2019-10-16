package gateway

import (
	"fmt"
	"net/http"

	"github.com/celer-network/sgn/x/global"
	"github.com/celer-network/sgn/x/subscribe"
	"github.com/cosmos/cosmos-sdk/types/rest"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
)

func (rs *RestServer) registerQueryRoutes() {
	rs.Mux.HandleFunc(
		"/global/latestBlock",
		latestBlockHandlerFn(rs),
	).Methods("GET")

	rs.Mux.HandleFunc(
		"/subscribe/params",
		subscribeParamsHandlerFn(rs),
	).Methods("GET")

	rs.Mux.HandleFunc(
		"/subscribe/subscription/{ethAddr}",
		subscriptionHandlerFn(rs),
	).Methods("GET")

	rs.Mux.HandleFunc(
		"/subscribe/request/{channelId}",
		guardRequestHandlerFn(rs),
	).Methods("GET")
}

// http request handler to query latest block
func latestBlockHandlerFn(rs *RestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", global.ModuleName, global.QueryLatestBlock)
		res, _, err := rs.transactor.CliCtx.Query(route)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, rs.transactor.CliCtx, res)
	}
}

// http request handler to query subscribe params
func subscribeParamsHandlerFn(rs *RestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", subscribe.ModuleName, subscribe.QueryParameters)
		res, _, err := rs.transactor.CliCtx.Query(route)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, rs.transactor.CliCtx, res)
	}
}

// http request handler to query subscription
func subscriptionHandlerFn(rs *RestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ethAddr := vars["ethAddr"]
		subscription, err := subscribe.CLIQuerySubscription(rs.transactor.CliCtx.Codec, rs.transactor.CliCtx, subscribe.RouterKey, ethAddr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, rs.transactor.CliCtx, subscription)
	}
}

// http request handler to query guard request
func guardRequestHandlerFn(rs *RestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		channelId := ethcommon.Hex2Bytes(vars["channelId"])
		request, err := subscribe.CLIQueryRequest(rs.transactor.CliCtx.Codec, rs.transactor.CliCtx, subscribe.RouterKey, channelId)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, rs.transactor.CliCtx, request)
	}
}
