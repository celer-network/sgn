package gateway

import (
	"net/http"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/transactor"
	"github.com/celer-network/sgn/x/subscribe"
	"github.com/celer-network/sgn/x/sync"
	"github.com/celer-network/sgn/x/validator"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

func (rs *RestServer) registerTxRoutes() {
	rs.Mux.HandleFunc(
		"/subscribe/subscribe",
		postSubscribeHandlerFn(rs),
	).Methods(http.MethodPost, http.MethodOptions)

	rs.Mux.HandleFunc(
		"/subscribe/request",
		postRequestGuardHandlerFn(rs),
	).Methods(http.MethodPost, http.MethodOptions)

	rs.Mux.HandleFunc(
		"/validator/updateSidechainAddr",
		postUpdateSidechainAddrHandlerFn(rs),
	).Methods(http.MethodPost, http.MethodOptions)

	rs.Mux.HandleFunc(
		"/validator/syncDelegator",
		postSyncDelegatorHandlerFn(rs),
	).Methods(http.MethodPost, http.MethodOptions)

	rs.Mux.HandleFunc(
		"/validator/withdrawReward",
		postWithdrawRewardHandlerFn(rs),
	).Methods(http.MethodPost, http.MethodOptions)
}

type (
	SubscribeRequest struct {
		EthAddr string `json:"ethAddr" yaml:"ethAddr"`
		Amount  string `json:"amount" yaml:"amount"`
	}

	RequestGuardRequest struct {
		EthAddr                 string `json:"ethAddr" yaml:"ethAddr"`
		SignedSimplexStateBytes string `json:"signedSimplexStateBytes" yaml:"signedSimplexStateBytes"`
	}

	UpdateSidechainAddrRequest struct {
		EthAddr string `json:"ethAddr" yaml:"ethAddr"`
	}

	SyncDelegatorRequest struct {
		CandidateAddress string `json:"candidateAddress"`
		DelegatorAddress string `json:"delegatorAddress"`
	}

	WithdrawRewardRequest struct {
		EthAddr string `json:"ethAddr" yaml:"ethAddr"`
	}
)

func postSubscribeHandlerFn(rs *RestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SubscribeRequest
		transactor := rs.transactorPool.GetTransactor()
		if !rest.ReadRESTReq(w, r, transactor.CliCtx.Codec, &req) {
			return
		}

		subscription := subscribe.NewSubscription(req.EthAddr)
		deposit, ok := sdk.NewIntFromString(req.Amount)
		if !ok {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid deposit amount")
			return
		}

		subscription.Deposit = deposit
		subscriptionData := transactor.CliCtx.Codec.MustMarshalBinaryBare(subscription)
		msg := sync.NewMsgSubmitChange(sync.Subscribe, subscriptionData, transactor.Key.GetAddress())
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

		request, err := subscribe.GetRequest(transactor.CliCtx, tc.Client0, req.SignedSimplexStateBytes)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Fail to get request from SignedSimplexStateBytes")
			return
		}

		request.SeqNum = seqNum
		request.SignedSimplexStateBytes = signedSimplexStateBytes
		request.OwnerSig = requestSig
		requestData := transactor.CliCtx.Codec.MustMarshalBinaryBare(request)
		msgSubmitChange = sync.NewMsgSubmitChange(sync.Request, requestData, transactor.Key.GetAddress())

		// FIX
		// signedSimplexStateBytes := mainchain.Hex2Bytes(req.SignedSimplexStateBytes)
		// msg := subscribe.NewMsgRequestGuard(req.EthAddr, signedSimplexStateBytes, transactor.CliCtx.GetFromAddress())
		// writeGenerateStdTxResponse(w, transactor, msg)
	}
}

func postUpdateSidechainAddrHandlerFn(rs *RestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UpdateSidechainAddrRequest
		transactor := rs.transactorPool.GetTransactor()
		if !rest.ReadRESTReq(w, r, transactor.CliCtx.Codec, &req) {
			return
		}

		// FIX
		// msg := validator.NewMsgUpdateSidechainAddr(req.EthAddr, transactor.CliCtx.GetFromAddress())
		// writeGenerateStdTxResponse(w, transactor, msg)
	}
}

func postSyncDelegatorHandlerFn(rs *RestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SyncDelegatorRequest
		transactor := rs.transactorPool.GetTransactor()
		if !rest.ReadRESTReq(w, r, transactor.CliCtx.Codec, &req) {
			return
		}

		// FIX
		// msg := validator.NewMsgSyncDelegator(req.CandidateAddress, req.DelegatorAddress, transactor.CliCtx.GetFromAddress())
		// writeGenerateStdTxResponse(w, transactor, msg)
	}
}

func postWithdrawRewardHandlerFn(rs *RestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "content-type")
		if r.Method == http.MethodOptions {
			return
		}

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

	transactor.AddTxMsg(msg)

	w.Header().Set("Content-Type", "text/plain")
	if _, err := w.Write([]byte("success")); err != nil {
		log.Errorln("could not write response:", err)
	}
}
