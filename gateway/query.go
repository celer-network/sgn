package gateway

import (
	"net/http"

	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/x/global"
	"github.com/celer-network/sgn/x/subscribe"
	"github.com/celer-network/sgn/x/validator"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
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

	rs.Mux.HandleFunc(
		"/validator/candidate/{ethAddr}",
		candidateHandlerFn(rs),
	).Methods("GET")

	rs.Mux.HandleFunc(
		"/validator/reward/{ethAddr}",
		rewardHandlerFn(rs),
	).Methods("GET")

	rs.Mux.HandleFunc(
		"/validator/rewardRequest/{ethAddr}",
		rewardRequestHandlerFn(rs),
	).Methods("GET")
}

// http request handler to query latest block
func latestBlockHandlerFn(rs *RestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		transactor := rs.transactorPool.GetTransactor()
		block, err := global.CLIQueryLatestBlock(transactor.CliCtx, global.RouterKey)
		postProcessResponse(w, transactor.CliCtx, block, err)
	}
}

// http request handler to query subscribe params
func subscribeParamsHandlerFn(rs *RestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		transactor := rs.transactorPool.GetTransactor()
		params, err := subscribe.CLIQueryParams(transactor.CliCtx, subscribe.RouterKey)

		postProcessResponse(w, transactor.CliCtx, params, err)
	}
}

// http request handler to query subscription
func subscriptionHandlerFn(rs *RestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ethAddr := vars["ethAddr"]
		transactor := rs.transactorPool.GetTransactor()
		subscription, err := subscribe.CLIQuerySubscription(transactor.CliCtx, subscribe.RouterKey, ethAddr)
		postProcessResponse(w, transactor.CliCtx, subscription, err)
	}
}

// http request handler to query guard request
func guardRequestHandlerFn(rs *RestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		channelId := mainchain.Hex2Bytes(vars["channelId"])
		transactor := rs.transactorPool.GetTransactor()
		request, err := subscribe.CLIQueryRequest(transactor.CliCtx, subscribe.RouterKey, channelId)
		postProcessResponse(w, transactor.CliCtx, request, err)
	}
}

// http request handler to query candidate
func candidateHandlerFn(rs *RestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ethAddr := vars["ethAddr"]
		transactor := rs.transactorPool.GetTransactor()
		candidate, err := validator.CLIQueryCandidate(transactor.CliCtx, validator.RouterKey, ethAddr)
		postProcessResponse(w, transactor.CliCtx, candidate, err)
	}
}

// http request handler to query reward
func rewardHandlerFn(rs *RestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ethAddr := vars["ethAddr"]
		transactor := rs.transactorPool.GetTransactor()
		reward, err := validator.CLIQueryReward(transactor.CliCtx, validator.RouterKey, ethAddr)
		postProcessResponse(w, transactor.CliCtx, reward, err)
	}
}

// http request handler to query reward request
func rewardRequestHandlerFn(rs *RestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ethAddr := vars["ethAddr"]
		transactor := rs.transactorPool.GetTransactor()
		reward, err := validator.CLIQueryReward(transactor.CliCtx, validator.RouterKey, ethAddr)
		postProcessResponse(w, transactor.CliCtx, reward.GetRewardRequest(), err)
	}
}

func postProcessResponse(w http.ResponseWriter, cliCtx context.CLIContext, resp interface{}, err error) {
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	rest.PostProcessResponse(w, cliCtx, resp)
}
