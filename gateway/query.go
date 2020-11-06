package gateway

import (
	"math/big"
	"net/http"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/proto/sgn"
	"github.com/celer-network/sgn/x/guard"
	"github.com/celer-network/sgn/x/validator"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
)

func (rs *RestServer) registerQueryRoutes() {
	rs.Mux.HandleFunc(
		"/guard/params",
		guardParamsHandlerFn(rs),
	).Methods(http.MethodGet, http.MethodOptions)

	rs.Mux.HandleFunc(
		"/guard/subscription/{ethAddr}",
		subscriptionHandlerFn(rs),
	).Methods(http.MethodGet, http.MethodOptions)

	rs.Mux.HandleFunc(
		"/guard/request/{channelId}/{simplexReceiver}",
		guardRequestHandlerFn(rs),
	).Methods(http.MethodGet, http.MethodOptions)

	rs.Mux.HandleFunc(
		"/validator/candidate/{ethAddr}",
		candidateHandlerFn(rs),
	).Methods(http.MethodGet, http.MethodOptions)

	rs.Mux.HandleFunc(
		"/validator/candidate-delegators/{ethAddr}",
		candidateDelegatorsHandlerFn(rs),
	).Methods(http.MethodGet, http.MethodOptions)

	rs.Mux.HandleFunc(
		"/validator/reward/{ethAddr}",
		rewardHandlerFn(rs),
	).Methods(http.MethodGet, http.MethodOptions)

	rs.Mux.HandleFunc(
		"/validator/rewardRequest/{ethAddr}",
		rewardRequestHandlerFn(rs),
	).Methods(http.MethodGet, http.MethodOptions)
}

// http request handler to query guard params
func guardParamsHandlerFn(rs *RestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		transactor := rs.transactorPool.GetTransactor()
		params, err := guard.CLIQueryParams(transactor.CliCtx, guard.RouterKey)

		postProcessResponse(w, transactor.CliCtx, params, err)
	}
}

// http request handler to query subscription
func subscriptionHandlerFn(rs *RestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ethAddr := vars["ethAddr"]
		transactor := rs.transactorPool.GetTransactor()
		subscription, err := guard.CLIQuerySubscription(transactor.CliCtx, guard.RouterKey, ethAddr)
		postProcessResponse(w, transactor.CliCtx, subscription, err)
	}
}

// http request handler to query guard request
func guardRequestHandlerFn(rs *RestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		channelId := mainchain.Hex2Bytes(vars["channelId"])
		simplexReceiver := vars["simplexReceiver"]
		transactor := rs.transactorPool.GetTransactor()
		request, err := guard.CLIQueryRequest(transactor.CliCtx, guard.RouterKey, channelId, simplexReceiver)
		postProcessResponse(w, transactor.CliCtx, request, err)
	}
}

// http request handler to query candidate
func candidateHandlerFn(rs *RestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == http.MethodOptions {
			return
		}

		vars := mux.Vars(r)
		ethAddr := vars["ethAddr"]
		transactor := rs.transactorPool.GetTransactor()
		candidate, err := validator.CLIQueryCandidate(transactor.CliCtx, validator.RouterKey, ethAddr)
		postProcessResponse(w, transactor.CliCtx, candidate, err)
	}
}

// http request handler to query candidate delegators
func candidateDelegatorsHandlerFn(rs *RestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == http.MethodOptions {
			return
		}

		vars := mux.Vars(r)
		ethAddr := vars["ethAddr"]
		transactor := rs.transactorPool.GetTransactor()
		delegators, err := validator.CLIQueryCandidateDelegators(transactor.CliCtx, validator.RouterKey, ethAddr)
		postProcessResponse(w, transactor.CliCtx, delegators, err)
	}
}

// http request handler to query reward
func rewardHandlerFn(rs *RestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == http.MethodOptions {
			return
		}

		vars := mux.Vars(r)
		ethAddr := vars["ethAddr"]
		transactor := rs.transactorPool.GetTransactor()
		reward, err := validator.CLIQueryReward(transactor.CliCtx, validator.RouterKey, ethAddr)
		postProcessResponse(w, transactor.CliCtx, reward, err)
	}
}

type RewardRequest struct {
	MiningReward            string `json:"mining_reward"`
	ServiceReward           string `json:"service_reward"`
	RewardRequestProtoBytes string `json:"reward_request_proto_bytes"` // proto msg for reward snapshot from latest intendWithdraw
}

// http request handler to query reward request
func rewardRequestHandlerFn(rs *RestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == http.MethodOptions {
			return
		}

		vars := mux.Vars(r)
		ethAddr := vars["ethAddr"]
		transactor := rs.transactorPool.GetTransactor()
		reward, err := validator.CLIQueryReward(transactor.CliCtx, validator.RouterKey, ethAddr)
		rewardRequest := RewardRequest{
			RewardRequestProtoBytes: mainchain.Bytes2Hex(reward.GetRewardRequest()),
		}

		if len(reward.RewardProtoBytes) != 0 {
			var pbReward sgn.Reward
			err = proto.Unmarshal(reward.RewardProtoBytes, &pbReward)
			if err != nil {
				log.Errorln("proto umarshal err", err, reward.RewardProtoBytes)
			} else {
				rewardRequest.MiningReward = new(big.Int).SetBytes(pbReward.CumulativeMiningReward).String()
				rewardRequest.ServiceReward = new(big.Int).SetBytes(pbReward.CumulativeServiceReward).String()
			}
		}

		postProcessResponse(w, transactor.CliCtx, rewardRequest, err)
	}
}

func postProcessResponse(w http.ResponseWriter, cliCtx context.CLIContext, resp interface{}, err error) {
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	rest.PostProcessResponse(w, cliCtx, resp)
}
