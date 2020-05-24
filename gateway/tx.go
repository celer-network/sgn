package gateway

import (
	"net/http"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/proto/chain"
	"github.com/celer-network/sgn/transactor"
	"github.com/celer-network/sgn/x/subscribe"
	"github.com/celer-network/sgn/x/sync"
	"github.com/celer-network/sgn/x/validator"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/golang/protobuf/proto"
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
		OwnerSig                string `json:"ownerSig" yaml:"ownerSig"`
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

		signedSimplexStateBytes := mainchain.Hex2Bytes(req.SignedSimplexStateBytes)
		var signedSimplexState chain.SignedSimplexState
		err := proto.Unmarshal(signedSimplexStateBytes, &signedSimplexState)
		if err != nil {
			log.Errorln("Failed to unmarshal signedSimplexStateBytes:", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Fail to unmarshal signedSimplexStateBytes")
			return
		}

		request, err := subscribe.GetRequest(transactor.CliCtx, rs.ethClient, &signedSimplexState)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Fail to get request from SignedSimplexStateBytes")
			return
		}

		request.SignedSimplexStateBytes = signedSimplexStateBytes
		request.OwnerSig = mainchain.Hex2Bytes(req.OwnerSig)
		requestData := transactor.CliCtx.Codec.MustMarshalBinaryBare(request)
		msg := sync.NewMsgSubmitChange(sync.Request, requestData, transactor.Key.GetAddress())
		writeGenerateStdTxResponse(w, transactor, msg)
	}
}

func postUpdateSidechainAddrHandlerFn(rs *RestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UpdateSidechainAddrRequest
		transactor := rs.transactorPool.GetTransactor()
		if !rest.ReadRESTReq(w, r, transactor.CliCtx.Codec, &req) {
			return
		}

		sidechainAddr, err := rs.ethClient.SGN.SidechainAddrMap(&bind.CallOpts{}, mainchain.Hex2Addr(req.EthAddr))
		if err != nil {
			log.Errorln("Query sidechain address error:", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Fail to query sidechain address")
			return
		}

		candidate := validator.NewCandidate(req.EthAddr, sdk.AccAddress(sidechainAddr))
		candidateData := transactor.CliCtx.Codec.MustMarshalBinaryBare(candidate)
		msg := sync.NewMsgSubmitChange(sync.UpdateSidechainAddr, candidateData, transactor.Key.GetAddress())
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

		di, err := rs.ethClient.DPoS.GetDelegatorInfo(&bind.CallOpts{}, mainchain.Hex2Addr(req.CandidateAddress), mainchain.Hex2Addr(req.DelegatorAddress))
		if err != nil {
			log.Errorf("Failed to query delegator info: %s", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Fail to query delegator info")
			return
		}

		delegator := validator.NewDelegator(req.CandidateAddress, req.DelegatorAddress)
		delegator.DelegatedStake = sdk.NewIntFromBigInt(di.DelegatedStake)
		delegatorData := transactor.CliCtx.Codec.MustMarshalBinaryBare(delegator)
		msg := sync.NewMsgSubmitChange(sync.SyncDelegator, delegatorData, transactor.Key.GetAddress())
		writeGenerateStdTxResponse(w, transactor, msg)
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
