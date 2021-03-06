package channel

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"math/big"
	"net/http"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/proto/chain"
	tc "github.com/celer-network/sgn/testing/common"
	"github.com/celer-network/sgn/x/guard"
	"github.com/celer-network/sgn/x/sync"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/golang/protobuf/proto"
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
		"/channelInfo",
		getChannelInfoHandlerFn(rs),
	).Methods(http.MethodGet)
}

type (
	GuardRequest struct {
		SeqNum uint64 `json:"seq_num"`
	}

	IntendSettleRequest struct {
		SeqNum uint64 `json:"seq_num"`
	}
)

func postRequestGuardHandlerFn(rs *RestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req GuardRequest
		if !rest.ReadRESTReq(w, r, rs.cdc, &req) {
			return
		}

		signedSimplexStateProto, err := tc.PrepareSignedSimplexState(req.SeqNum, rs.channelID[:], rs.peer1.Address.Bytes(), rs.peer0, rs.peer1)
		if err != nil {
			log.Errorln("could not get SignedSimplexState:", err)
			return
		}

		signedSimplexStateBytes, err := proto.Marshal(signedSimplexStateProto)
		if err != nil {
			log.Errorln("could not marshal SignedSimplexState:", err)
			return
		}

		simplexReceiverSig, err := rs.peer0.Signer.SignEthMessage(signedSimplexStateBytes)
		if err != nil {
			return
		}

		if rs.gateway == "" {
			disputeTimeout, err := tc.LedgerContract.GetDisputeTimeout(&bind.CallOpts{}, rs.channelID)
			if err != nil {
				log.Errorln("Failed to get dispute timeout:", err)
				return
			}
			request := guard.NewInitRequest(signedSimplexStateBytes, simplexReceiverSig, disputeTimeout.Uint64())
			syncData := rs.transactor.CliCtx.Codec.MustMarshalBinaryBare(request)
			msg := sync.NewMsgSubmitChange(sync.InitGuardRequest, syncData, tc.EthClient, rs.transactor.Key.GetAddress())
			rs.transactor.AddTxMsg(msg)
		} else {
			reqBody, err := json.Marshal(map[string]string{
				"signed_simplex_state_bytes": mainchain.Bytes2Hex(signedSimplexStateBytes),
				"simplex_receiver_sig":       mainchain.Bytes2Hex(simplexReceiverSig),
			})
			if err != nil {
				return
			}
			resp, err := http.Post(rs.gateway+"/guard/requestGuard",
				"application/json", bytes.NewBuffer(reqBody))
			if err != nil {
				return
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return
			}
			if resp.StatusCode != http.StatusOK {
				log.Error(string(body))
				return
			}
		}

		w.Header().Set("Content-Type", "text/plain")
		if _, err := w.Write([]byte("success\n")); err != nil {
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

		signedSimplexStateProto, err := tc.PrepareSignedSimplexState(req.SeqNum, rs.channelID[:], rs.peer1.Address.Bytes(), rs.peer0, rs.peer1)
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

		_, err = tc.LedgerContract.IntendSettle(rs.peer1.Auth, signedSimplexStateArrayBytes)
		if err != nil {
			log.Errorln("could not intendSettle:", err)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		if _, err := w.Write([]byte("success\n")); err != nil {
			log.Errorln("could not write response:", err)
		}
	}
}

func getChannelInfoHandlerFn(rs *RestServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status, err := tc.LedgerContract.GetChannelStatus(&bind.CallOpts{}, rs.channelID)
		if err != nil {
			log.Errorln("Query ChannelStatus err", err)
			return
		}

		var statusStr string
		switch status {
		case 0:
			statusStr = "Uninitialized"
		case 1:
			statusStr = "Operable"
		case 2:
			statusStr = "Settling"
		case 3:
			statusStr = "Closed"
		case 4:
			statusStr = "Migrated"
		}

		addresses, seqNums, err := tc.LedgerContract.GetStateSeqNumMap(&bind.CallOpts{}, rs.channelID)
		if err != nil {
			log.Errorln("Query StateSeqNumMap err", err)
			return
		}

		result, err := rs.cdc.MarshalJSON(struct {
			Status    string
			Addresses [2]mainchain.Addr
			SeqNums   [2]*big.Int
		}{statusStr, addresses, seqNums})

		if err != nil {
			log.Errorln("MarshalJSON err", err)
			return
		}

		if _, err := w.Write(result); err != nil {
			log.Errorln("could not write response:", err)
		}
	}
}
