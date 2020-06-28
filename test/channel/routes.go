package channel

import (
	"bytes"
	"encoding/json"
	"math/big"
	"net/http"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/proto/chain"
	tc "github.com/celer-network/sgn/test/common"
	"github.com/celer-network/sgn/x/subscribe"
	"github.com/celer-network/sgn/x/sync"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/golang/protobuf/proto"
)

func (rs *RestServer) registerRoutes() {
	rs.Mux.HandleFunc(
		"/requestGuard",
		postInitGuardHandlerFn(rs),
	).Methods(http.MethodPost)

	rs.Mux.HandleFunc(
		"/intendSettle",
		postIntendSettleHandlerFn(rs),
	).Methods(http.MethodPost)

	rs.Mux.HandleFunc(
		"/channelInfo",
		getChannelInfoHandlerFn(rs),
	).Methods(http.MethodGet)
}

type (
	RequestGuardRequest struct {
		SeqNum uint64 `json:"seqNum"`
	}

	IntendSettleRequest struct {
		SeqNum uint64 `json:"seqNum"`
	}
)

func postInitGuardHandlerFn(rs *RestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestGuardRequest
		if !rest.ReadRESTReq(w, r, rs.cdc, &req) {
			return
		}

		signedSimplexStateProto, err := tc.PrepareSignedSimplexState(req.SeqNum, rs.channelID[:], rs.peer2.Address.Bytes(), rs.peer1, rs.peer2)
		if err != nil {
			log.Errorln("could not get SignedSimplexState:", err)
			return
		}

		signedSimplexStateBytes, err := proto.Marshal(signedSimplexStateProto)
		if err != nil {
			log.Errorln("could not marshal SignedSimplexState:", err)
			return
		}

		ownerSig, err := rs.peer1.Signer.SignEthMessage(signedSimplexStateBytes)
		if err != nil {
			return
		}

		if rs.gateway == "" {
			_, peerAddrs, peerFromIndex, err := subscribe.GetOnChainChannelSeqAndPeerIndex(
				tc.LedgerContract, rs.channelID, rs.peer2.Address)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, "Fail to get request onchain channel info")
				return
			}
			request := subscribe.NewRequest(
				rs.channelID.Bytes(),
				req.SeqNum,
				peerAddrs,
				peerFromIndex,
				signedSimplexStateBytes,
				ownerSig)
			requestData := rs.transactor.CliCtx.Codec.MustMarshalBinaryBare(request)
			msg := sync.NewMsgSubmitChange(sync.InitGuardRequest, requestData, rs.transactor.Key.GetAddress())
			rs.transactor.AddTxMsg(msg)
		} else {
			reqBody, err := json.Marshal(map[string]string{
				"ownerSig":                mainchain.Bytes2Hex(ownerSig),
				"signedSimplexStateBytes": mainchain.Bytes2Hex(signedSimplexStateBytes),
			})
			if err != nil {
				return
			}
			_, err = http.Post(rs.gateway+"/subscribe/initGuard",
				"application/json", bytes.NewBuffer(reqBody))
			if err != nil {
				return
			}
		}

		w.Header().Set("Content-Type", "text/plain")
		if _, err := w.Write([]byte("success")); err != nil {
			log.Errorln("could not write response:", err)
		}
	}
}

func postIntendSettleHandlerFn(rs *RestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req IntendSettleRequest
		if !rest.ReadRESTReq(w, r, rs.cdc, &req) {
			return
		}

		signedSimplexStateProto, err := tc.PrepareSignedSimplexState(req.SeqNum, rs.channelID[:], rs.peer2.Address.Bytes(), rs.peer1, rs.peer2)
		if err != nil {
			log.Errorln("could not get SignedSimplexState:", err)
			return
		}

		signedSimplexStateArrayBytes, err := proto.Marshal(&chain.SignedSimplexStateArray{
			SignedSimplexStates: []*chain.SignedSimplexState{signedSimplexStateProto},
		})
		if err != nil {
			log.Errorln("could not marshal SignedSimplexStateArray:", err)
			return
		}

		_, err = tc.LedgerContract.IntendSettle(rs.peer2.Auth, signedSimplexStateArrayBytes)
		if err != nil {
			log.Errorln("could not intendSettle:", err)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		if _, err := w.Write([]byte("success")); err != nil {
			log.Errorln("could not write response:", err)
		}
	}
}

func getChannelInfoHandlerFn(rs *RestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		addresses, seqNums, err := tc.LedgerContract.GetStateSeqNumMap(&bind.CallOpts{}, rs.channelID)
		if err != nil {
			log.Errorln("Query StateSeqNumMap err", err)
			return
		}

		result, err := rs.cdc.MarshalJSON(struct {
			Addresses [2]mainchain.Addr
			SeqNums   [2]*big.Int
		}{addresses, seqNums})

		if err != nil {
			log.Errorln("MarshalJSON err", err)
			return
		}

		if _, err := w.Write(result); err != nil {
			log.Errorln("could not write response:", err)
		}
	}
}
