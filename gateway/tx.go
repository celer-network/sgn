package gateway

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/celer-network/goutils/eth"
	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/transactor"
	"github.com/celer-network/sgn/x/guard"
	"github.com/celer-network/sgn/x/sync"
	"github.com/celer-network/sgn/x/validator"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func (rs *RestServer) registerTxRoutes() {
	rs.Mux.HandleFunc(
		"/guard/requestGuard",
		postRequestGuardHandlerFn(rs),
	).Methods(http.MethodPost, http.MethodOptions)

	rs.Mux.HandleFunc(
		"/validator/withdrawReward",
		postWithdrawRewardHandlerFn(rs),
	).Methods(http.MethodPost, http.MethodOptions)

	/*
		rs.Mux.HandleFunc(
			"/guard/subscribe",
			postSubscribeHandlerFn(rs),
		).Methods(http.MethodPost, http.MethodOptions)

		rs.Mux.HandleFunc(
			"/validator/updateSidechainAddr",
			postUpdateSidechainAddrHandlerFn(rs),
		).Methods(http.MethodPost, http.MethodOptions)

		rs.Mux.HandleFunc(
			"/validator/syncDelegator",
			postSyncDelegatorHandlerFn(rs),
		).Methods(http.MethodPost, http.MethodOptions)
	*/
}

type (
	SubscribeRequest struct {
		EthAddr string `json:"ethAddr" yaml:"ethAddr"`
		Amount  string `json:"amount" yaml:"amount"`
	}

	GuardRequest struct {
		SignedSimplexStateBytes string `json:"signedSimplexStateBytes" yaml:"signedSimplexStateBytes"`
		SimplexReceiverSig      string `json:"simplexReceiverSig" yaml:"simplexReceiverSig"`
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

func postRequestGuardHandlerFn(rs *RestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req GuardRequest
		transactor := rs.transactorPool.GetTransactor()
		if !rest.ReadRESTReq(w, r, transactor.CliCtx.Codec, &req) {
			return
		}

		simplexReceiverSig := mainchain.Hex2Bytes(req.SimplexReceiverSig)
		signedSimplexStateBytes := mainchain.Hex2Bytes(req.SignedSimplexStateBytes)
		signedSimplexState, simplexChannel, err := common.UnmarshalSignedSimplexStateBytes(signedSimplexStateBytes)
		if err != nil {
			errmsg := fmt.Sprintf("UnmarshalSignedSimplexStateBytes err: %s", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, errmsg)
			return
		}

		simplexReceiver, err := eth.RecoverSigner(signedSimplexStateBytes, simplexReceiverSig)
		if err != nil {
			errmsg := fmt.Sprintf("RecoverSigner err: %s", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, errmsg)
			return
		}

		// verify signature
		simplexSender := mainchain.Bytes2Addr(simplexChannel.PeerFrom)
		err = guard.VerifySimplexStateSigs(signedSimplexState, simplexSender, simplexReceiver)
		if err != nil {
			errmsg := fmt.Sprintf("Invalid signature: %s", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, errmsg)
			return
		}

		// check existing record
		storedRequest, err := guard.CLIQueryRequest(
			transactor.CliCtx, guard.RouterKey, simplexChannel.ChannelId, mainchain.Addr2Hex(simplexReceiver))
		if err == nil {
			if simplexChannel.SeqNum <= storedRequest.SeqNum {
				errmsg := fmt.Sprintf("Invalid sequence number. request: %d, stored: %d", simplexChannel.SeqNum, storedRequest.SeqNum)
				rest.WriteErrorResponse(w, http.StatusBadRequest, errmsg)
				return
			}
			if mainchain.Hex2Addr(storedRequest.SimplexSender) != mainchain.Bytes2Addr(simplexChannel.PeerFrom) {
				errmsg := fmt.Sprintf("Sender not match stored request: %s", storedRequest.SimplexSender)
				rest.WriteErrorResponse(w, http.StatusBadRequest, errmsg)
				return
			}
			if mainchain.Hex2Addr(storedRequest.SimplexReceiver) != simplexReceiver {
				errmsg := fmt.Sprintf("Receiver not match stored request: %s", storedRequest.SimplexReceiver)
				rest.WriteErrorResponse(w, http.StatusBadRequest, errmsg)
				return
			}
			if storedRequest.Status != guard.ChanStatus_Idle {
				errmsg := fmt.Sprintf("Guard state is not idle: %d", storedRequest.Status)
				rest.WriteErrorResponse(w, http.StatusBadRequest, errmsg)
				return
			}
			msg := guard.NewMsgRequestGuard(signedSimplexStateBytes, simplexReceiverSig, transactor.Key.GetAddress())
			writeGenerateStdTxResponse(w, transactor, msg)
			return
		} else if !strings.Contains(err.Error(), common.ErrRecordNotFound.Error()) {
			errmsg := fmt.Sprintf("Failed to get request: %s", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, errmsg)
			return
		}

		// Initialize first guard request
		// verify peer addr
		cid := mainchain.Bytes2Cid(simplexChannel.ChannelId)
		addrs, seqNums, err := rs.ledgerContract.GetStateSeqNumMap(&bind.CallOpts{}, cid)
		if err != nil {
			errmsg := fmt.Sprintf("Failed to get mainchain channel info: %s", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, errmsg)
			return
		}
		seqIndex := 0
		var match bool
		if simplexSender == addrs[0] {
			match = (simplexReceiver == addrs[1])
		} else if simplexSender == addrs[1] {
			match = (simplexReceiver == addrs[0])
			seqIndex = 1
		}
		if !match {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "channel peers not match")
			return
		}

		// verify seq
		if simplexChannel.SeqNum <= seqNums[seqIndex].Uint64() {
			errmsg := fmt.Sprintf("invalid sequence number, request: %d, mainchain: %d",
				simplexChannel.SeqNum, seqNums[seqIndex].Uint64())
			rest.WriteErrorResponse(w, http.StatusBadRequest, errmsg)
			return
		}

		disputeTimeout, err := rs.ledgerContract.GetDisputeTimeout(&bind.CallOpts{}, cid)
		if err != nil {
			errmsg := fmt.Sprintf("Failed to get dispute timeout: %s", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, errmsg)
			return
		}

		request := guard.NewInitRequest(signedSimplexStateBytes, simplexReceiverSig, disputeTimeout.Uint64())
		syncData := transactor.CliCtx.Codec.MustMarshalBinaryBare(request)
		msg := transactor.NewMsgSubmitChange(sync.InitGuardRequest, syncData, rs.ethClient)
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
		reward, err := validator.CLIQueryReward(transactor.CliCtx, validator.RouterKey, req.EthAddr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if !reward.HasNewReward() {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "reward request is already the latest")
			return
		}

		msg := validator.NewMsgWithdrawReward(req.EthAddr, transactor.CliCtx.GetFromAddress())
		writeGenerateStdTxResponse(w, transactor, msg)
	}
}

func postSubscribeHandlerFn(rs *RestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SubscribeRequest
		transactor := rs.transactorPool.GetTransactor()
		if !rest.ReadRESTReq(w, r, transactor.CliCtx.Codec, &req) {
			return
		}

		subscription := guard.NewSubscription(req.EthAddr)
		deposit, ok := sdk.NewIntFromString(req.Amount)
		if !ok {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid deposit amount")
			return
		}

		subscription.Deposit = deposit
		subscriptionData := transactor.CliCtx.Codec.MustMarshalBinaryBare(subscription)
		msg := transactor.NewMsgSubmitChange(sync.Subscribe, subscriptionData, rs.ethClient)
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
		msg := transactor.NewMsgSubmitChange(sync.UpdateSidechainAddr, candidateData, rs.ethClient)
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
		msg := transactor.NewMsgSubmitChange(sync.SyncDelegator, delegatorData, rs.ethClient)
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
