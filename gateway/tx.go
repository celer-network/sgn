package gateway

import (
	"log"
	"net/http"

	"github.com/celer-network/sgn/transactor"
	"github.com/celer-network/sgn/x/subscribe"
	"github.com/celer-network/sgn/x/validator"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

func (rs *RestServer) registerTxRoutes() {
	rs.Mux.HandleFunc(
		"/subscribe/subscribe",
		postSubscribeHandlerFn(rs),
	).Methods("POST")

	rs.Mux.HandleFunc(
		"/subscribe/request",
		postRequestGuardHandlerFn(rs),
	).Methods("POST")

	rs.Mux.HandleFunc(
		"/validator/initializeCandidate",
		postInitializeCandidateHandlerFn(rs),
	).Methods("POST")

	rs.Mux.HandleFunc(
		"/validator/syncDelegator",
		postSyncDelegatorHandlerFn(rs),
	).Methods("POST")

	rs.Mux.HandleFunc(
		"/validator/withdrawReward",
		postWithdrawRewardHandlerFn(rs),
	).Methods("POST")
}

type (
	ethAddr struct {
		EthAddr string `json:"ethAddr" yaml:"ethAddr"`
	}

	SubscribeRequest struct {
		ethAddr
	}

	RequestGuardRequest struct {
		ethAddr
		SignedSimplexStateBytes string `json:"signedSimplexStateBytes" yaml:"signedSimplexStateBytes"`
	}

	InitializeCandidateRequest struct {
		ethAddr
	}

	SyncDelegatorRequest struct {
		CandidateAddress string `json:"candidateAddress"`
		DelegatorAddress string `json:"delegatorAddress"`
	}

	WithdrawRewardRequest struct {
		ethAddr
	}
)

func postSubscribeHandlerFn(rs *RestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SubscribeRequest
		transactor := rs.transactorPool.GetTransactor()
		if !rest.ReadRESTReq(w, r, transactor.CliCtx.Codec, &req) {
			return
		}

		msg := subscribe.NewMsgSubscribe(req.EthAddr, transactor.CliCtx.GetFromAddress())
		writeGenerateStdTxResponse(w, transactor, msg)
	}
}

func postRequestGuardHandlerFn(rs *RestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestGuardRequest
		transactor := rs.transactorPool.GetTransactor()
		if !rest.ReadRESTReq(w, r, transactor.CliCtx.Codec, &req) {
			return
		}

		signedSimplexStateBytes := ethcommon.Hex2Bytes(req.SignedSimplexStateBytes)
		msg := subscribe.NewMsgRequestGuard(req.EthAddr, signedSimplexStateBytes, transactor.CliCtx.GetFromAddress())
		writeGenerateStdTxResponse(w, transactor, msg)
	}
}

func postInitializeCandidateHandlerFn(rs *RestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req InitializeCandidateRequest
		transactor := rs.transactorPool.GetTransactor()
		if !rest.ReadRESTReq(w, r, transactor.CliCtx.Codec, &req) {
			return
		}

		msg := validator.NewMsgInitializeCandidate(req.EthAddr, transactor.CliCtx.GetFromAddress())
		writeGenerateStdTxResponse(w, transactor, msg)
	}
}

func postSyncDelegatorHandlerFn(rs *RestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SyncDelegatorRequest
		transactor := rs.transactorPool.GetTransactor()
		if !rest.ReadRESTReq(w, r, transactor.CliCtx.Codec, &req) {
			return
		}

		msg := validator.NewMsgSyncDelegator(req.CandidateAddress, req.DelegatorAddress, transactor.CliCtx.GetFromAddress())
		writeGenerateStdTxResponse(w, transactor, msg)
	}
}

func postWithdrawRewardHandlerFn(rs *RestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req WithdrawRewardRequest
		transactor := rs.transactorPool.GetTransactor()
		if !rest.ReadRESTReq(w, r, transactor.CliCtx.Codec, &req) {
			return
		}

		msg := validator.NewMsgWithdrawReward(req.EthAddr, transactor.CliCtx.GetFromAddress())
		writeGenerateStdTxResponse(w, transactor, msg)
	}
}

func writeGenerateStdTxResponse(w http.ResponseWriter, transactor *transactor.Transactor, msg sdk.Msg) {
	if err := msg.ValidateBasic(); err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	transactor.BroadcastTx(msg)

	w.Header().Set("Content-Type", "text/plain")
	if _, err := w.Write([]byte("success")); err != nil {
		log.Printf("could not write response: %v", err)
	}
}
