package osp

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/proto/chain"
	tf "github.com/celer-network/sgn/testing"
	"github.com/celer-network/sgn/x/subscribe"
	"github.com/cosmos/cosmos-sdk/types/rest"
	protobuf "github.com/golang/protobuf/proto"
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

		signedSimplexStateProto, err := tf.PrepareSignedSimplexState(req.SeqNum, rs.channelID[:], rs.user.Address.Bytes(), rs.user.PrivateKey, rs.osp.PrivateKey)
		if err != nil {
			log.Errorln("could not get SignedSimplexState:", err)
			return
		}

		signedSimplexStateBytes, err := protobuf.Marshal(signedSimplexStateProto)
		if err != nil {
			log.Errorln("could not marshal SignedSimplexState:", err)
			return
		}

		if rs.gateway == "" {
			msgRequestGuard := subscribe.NewMsgRequestGuard(rs.user.Address.Hex(), signedSimplexStateBytes, rs.transactor.Key.GetAddress())
			rs.transactor.AddTxMsg(msgRequestGuard)
		} else {
			reqBody, err := json.Marshal(map[string]string{
				"ethAddr":                 rs.user.Address.Hex(),
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

		signedSimplexStateProto, err := tf.PrepareSignedSimplexState(req.SeqNum, rs.channelID[:], rs.user.Address.Bytes(), rs.user.PrivateKey, rs.osp.PrivateKey)
		if err != nil {
			log.Errorln("could not get SignedSimplexState:", err)
			return
		}

		signedSimplexStateArrayBytes, err := protobuf.Marshal(&chain.SignedSimplexStateArray{
			SignedSimplexStates: []*chain.SignedSimplexState{signedSimplexStateProto},
		})
		if err != nil {
			log.Errorln("could not marshal SignedSimplexStateArray:", err)
			return
		}

		_, err = rs.osp.Ledger.IntendSettle(rs.osp.Auth, signedSimplexStateArrayBytes)
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
