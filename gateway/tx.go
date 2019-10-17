package gateway

import (
	"log"
	"net/http"

	"github.com/celer-network/sgn/utils"
	"github.com/celer-network/sgn/x/subscribe/types"
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
}

type (
	// SubscribeRequest defines the properties of a subscribe request's body.
	SubscribeRequest struct {
		EthAddr string `json:"ethAddr" yaml:"ethAddr"`
	}

	// RequestGuardRequest defines the properties of a request guard request's body.
	RequestGuardRequest struct {
		EthAddr                 string `json:"ethAddr" yaml:"ethAddr"`
		SignedSimplexStateBytes string `json:"signedSimplexStateBytes" yaml:"signedSimplexStateBytes"`
	}
)

func postSubscribeHandlerFn(rs *RestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SubscribeRequest
		if !rest.ReadRESTReq(w, r, rs.transactor.CliCtx.Codec, &req) {
			return
		}

		msg := types.NewMsgSubscribe(req.EthAddr, rs.transactor.CliCtx.GetFromAddress())
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		writeGenerateStdTxResponse(w, rs.transactor, msg)
	}
}

func postRequestGuardHandlerFn(rs *RestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestGuardRequest
		if !rest.ReadRESTReq(w, r, rs.transactor.CliCtx.Codec, &req) {
			return
		}

		signedSimplexStateBytes := ethcommon.Hex2Bytes(req.SignedSimplexStateBytes)
		msg := types.NewMsgRequestGuard(req.EthAddr, signedSimplexStateBytes, rs.transactor.CliCtx.GetFromAddress())
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		writeGenerateStdTxResponse(w, rs.transactor, msg)
	}
}

func writeGenerateStdTxResponse(w http.ResponseWriter, transactor *utils.Transactor, msg sdk.Msg) {
	transactor.BroadcastTx(msg)

	w.Header().Set("Content-Type", "text/plain")
	if _, err := w.Write([]byte("success")); err != nil {
		log.Printf("could not write response: %v", err)
	}
}
