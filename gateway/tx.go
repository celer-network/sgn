package gateway

import (
	"net/http"

	"github.com/celer-network/goutils/eth"
	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/transactor"
	"github.com/celer-network/sgn/x/subscribe"
	"github.com/celer-network/sgn/x/sync"
	"github.com/celer-network/sgn/x/validator"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func (rs *RestServer) registerTxRoutes() {
	rs.Mux.HandleFunc(
		"/subscribe/subscribe",
		postSubscribeHandlerFn(rs),
	).Methods(http.MethodPost, http.MethodOptions)

	rs.Mux.HandleFunc(
		"/subscribe/initGuardRequest",
		postInitGuardRequestHandlerFn(rs),
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

func postInitGuardRequestHandlerFn(rs *RestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestGuardRequest
		transactor := rs.transactorPool.GetTransactor()
		if !rest.ReadRESTReq(w, r, transactor.CliCtx.Codec, &req) {
			return
		}
		ownerSig := mainchain.Hex2Bytes(req.OwnerSig)
		signedSimplexStateBytes := mainchain.Hex2Bytes(req.SignedSimplexStateBytes)
		_, simplexChannel, err := common.UnmarshalSignedSimplexStateBytes(signedSimplexStateBytes)
		if err != nil {
			log.Errorln("Failed UnmarshalSignedSimplexStateBytes:", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Fail UnmarshalSignedSimplexStateBytes")
			return
		}

		_, err = subscribe.CLIQueryRequest(
			transactor.CliCtx, subscribe.RouterKey, simplexChannel.ChannelId, mainchain.Bytes2Hex(simplexChannel.PeerFrom))
		if err == nil {
			log.Errorf("Request for channel %x owner %x already initiated", simplexChannel.ChannelId, simplexChannel.PeerFrom)
			rest.WriteErrorResponse(w, http.StatusBadRequest, "request for channel and owner already initiated")
			return
		}

		ownerAddr, err := eth.RecoverSigner(signedSimplexStateBytes, ownerSig)
		if err != nil {
			log.Errorln("recover signer err:", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, "recover signer err")
			return
		}

		if mainchain.Bytes2Addr(simplexChannel.PeerFrom) != ownerAddr {
			log.Errorf("Owner signer %x does not match peerFrom: %x", ownerAddr, simplexChannel.PeerFrom)
			rest.WriteErrorResponse(w, http.StatusBadRequest, "owner signer not match")
			return
		}

		_, peerAddrs, peerFromIndex, err := subscribe.GetOnChainChannelSeqAndPeerIndex(
			rs.ledgerContract, mainchain.Bytes2Cid(simplexChannel.ChannelId), mainchain.Bytes2Addr(simplexChannel.PeerFrom))
		if err != nil {
			log.Errorln("Failed to get onchain channel info:", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, "get onchain channel info")
			return
		}

		// TODO: verify request before send
		request := subscribe.NewRequest(
			simplexChannel.ChannelId,
			simplexChannel.SeqNum,
			peerAddrs,
			peerFromIndex,
			signedSimplexStateBytes,
			ownerSig)
		requestData := transactor.CliCtx.Codec.MustMarshalBinaryBare(request)
		msg := sync.NewMsgSubmitChange(sync.InitGuardRequest, requestData, transactor.Key.GetAddress())
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

		sidechainAddr, err := rs.sgnContract.SidechainAddrMap(&bind.CallOpts{}, mainchain.Hex2Addr(req.EthAddr))
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

		di, err := rs.dposContract.GetDelegatorInfo(&bind.CallOpts{}, mainchain.Hex2Addr(req.CandidateAddress), mainchain.Hex2Addr(req.DelegatorAddress))
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
