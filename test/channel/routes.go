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
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
)

func (rs *RestServer) registerRoutes() {
	rs.Mux.HandleFunc(
		"/requestGuard",
		postRequestGuardHandlerFn(rs),
	).Methods(http.MethodPost)

	rs.Mux.HandleFunc(
		"/intendSettle",
		postIntendSettleHandlerFn(rs),
	).Methods(http.MethodPost)

	rs.Mux.HandleFunc(
		"/channelInfo/{channelId}",
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

func postRequestGuardHandlerFn(rs *RestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestGuardRequest
		if !rest.ReadRESTReq(w, r, rs.cdc, &req) {
			return
		}

		signedSimplexStateProto, err := tc.PrepareSignedSimplexState(req.SeqNum, rs.channelID[:], rs.peer1.Address.Bytes(), rs.peer1, rs.peer2)
		if err != nil {
			log.Errorln("could not get SignedSimplexState:", err)
			return
		}

		signedSimplexStateBytes, err := proto.Marshal(signedSimplexStateProto)
		if err != nil {
			log.Errorln("could not marshal SignedSimplexState:", err)
			return
		}

		if rs.gateway == "" {
			msgRequestGuard := subscribe.NewMsgRequestGuard(rs.peer1.Address.Hex(), signedSimplexStateBytes, rs.transactor.Key.GetAddress())
			rs.transactor.AddTxMsg(msgRequestGuard)
		} else {
			reqBody, err := json.Marshal(map[string]string{
				"ethAddr":                 rs.peer1.Address.Hex(),
				"signedSimplexStateBytes": mainchain.Bytes2Hex(signedSimplexStateBytes),
			})
			if err != nil {
				return
			}
			_, err = http.Post(rs.gateway+"/subscribe/request",
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

		signedSimplexStateProto, err := tc.PrepareSignedSimplexState(req.SeqNum, rs.channelID[:], rs.peer1.Address.Bytes(), rs.peer1, rs.peer2)
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

		_, err = rs.peer2.Ledger.IntendSettle(rs.peer2.Auth, signedSimplexStateArrayBytes)
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
		vars := mux.Vars(r)
		channelId := vars["channelId"]
		addresses, seqNums, err := rs.peer1.Ledger.GetStateSeqNumMap(&bind.CallOpts{}, mainchain.Hex2Cid(channelId))
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
