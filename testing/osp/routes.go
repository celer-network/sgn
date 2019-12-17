package osp

import (
	"net/http"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/proto/chain"
	tf "github.com/celer-network/sgn/testing"
	"github.com/celer-network/sgn/x/subscribe"
	"github.com/cosmos/cosmos-sdk/types/rest"
	protobuf "github.com/golang/protobuf/proto"
)

func (rs *RestServer) registerRoutes() {
	rs.Mux.HandleFunc(
		"/intendSettle",
		postIntendSettleHandlerFn(rs),
	).Methods("POST")
}

type (
	IntendSettleRequest struct {
		SeqNum uint64 `json:"seqNum"`
	}

	RequestGuardRequest struct {
		SeqNum uint64 `json:"seqNum"`
	}
)

func postIntendSettleHandlerFn(rs *RestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req IntendSettleRequest
		if !rest.ReadRESTReq(w, r, rs.transactor.CliCtx.Codec, &req) {
			return
		}

		signedSimplexStateProto, err := tf.PrepareSignedSimplexState(req.SeqNum, rs.channelID[:], rs.user.Address.Bytes(), rs.osp.PrivateKey, rs.user.PrivateKey)
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

		_, err = tf.DefaultTestEthClient.Ledger.IntendSettle(rs.osp.Auth, signedSimplexStateArrayBytes)
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

func postRequestGuardHandlerFn(rs *RestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestGuardRequest
		if !rest.ReadRESTReq(w, r, rs.transactor.CliCtx.Codec, &req) {
			return
		}

		signedSimplexStateProto, err := tf.PrepareSignedSimplexState(req.SeqNum, rs.channelID[:], rs.user.Address.Bytes(), rs.osp.PrivateKey, rs.user.PrivateKey)
		if err != nil {
			log.Errorln("could not get SignedSimplexState:", err)
			return
		}

		signedSimplexStateBytes, err := protobuf.Marshal(signedSimplexStateProto)
		if err != nil {
			log.Errorln("could not marshal SignedSimplexState:", err)
			return
		}
		msgRequestGuard := subscribe.NewMsgRequestGuard(rs.user.Address.Hex(), signedSimplexStateBytes, rs.transactor.Key.GetAddress())
		rs.transactor.AddTxMsg(msgRequestGuard)

		w.Header().Set("Content-Type", "text/plain")
		if _, err := w.Write([]byte("success")); err != nil {
			log.Errorln("could not write response:", err)
		}
	}
}
