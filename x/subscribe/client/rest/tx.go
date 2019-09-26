package rest

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/celer-network/sgn/x/subscribe/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/subscribe/subscribe",
		postSubscribeHandlerFn(cliCtx),
	).Methods("POST")
}

type (
	// SubscribeRequest defines the properties of a subscribe request's body.
	SubscribeRequest struct {
		BaseReq rest.BaseReq `json:"baseReq" yaml:"baseReq"`
		EthAddr string       `json:"ethAddr" yaml:"ethAddr"`
	}

	// RequestGuardRequest defines the properties of a subscribe request's body.
	RequestGuardRequest struct {
		BaseReq                 rest.BaseReq `json:"baseReq" yaml:"baseReq"`
		EthAddr                 string       `json:"ethAddr" yaml:"ethAddr"`
		SignedSimplexStateBytes string       `json:"signedSimplexStateBytes" yaml:"signedSimplexStateBytes"`
	}
)

func postSubscribeHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SubscribeRequest

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgSubscribe(req.EthAddr, fromAddr)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func postRequestGuardHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestGuardRequest

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		signedSimplexStateBytes := ethcommon.Hex2Bytes(req.SignedSimplexStateBytes)
		msg := types.NewMsgRequestGuard(req.EthAddr, signedSimplexStateBytes, fromAddr)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
