package types

import (
	"fmt"
	"strings"

	"github.com/celer-network/sgn/mainchain"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Request struct {
	ChannelId               []byte           `json:"channelId"`
	SeqNum                  uint64           `json:"seqNum"`
	SimplexSender           string           `json:"simplexSender"`
	SimplexReceiver         string           `json:"simplexReceiver"`
	SignedSimplexStateBytes []byte           `json:"signedSimplexStateBytes"`
	DisputeTimeout          uint64           `json:"disputeTimeout"`
	RequestGuards           []sdk.AccAddress `json:"requestGuards"`
	TriggerTxHash           string           `json:"triggerTxHash"`
	TriggerTxBlkNum         uint64           `json:"triggerTxBlkNum"`
	GuardTxHash             string           `json:"guardTxHash"`
	GuardTxBlkNum           uint64           `json:"guardTxBlkNum"`
	GuardSender             string           `json:"guardSender"`
}

func NewRequest(
	channelId []byte, seqNum uint64,
	senderAddr, receiverAddr mainchain.Addr,
	signedSimplexStateBytes []byte,
	disputeTimeout uint64) Request {
	return Request{
		ChannelId:               channelId,
		SeqNum:                  seqNum,
		SimplexSender:           mainchain.Addr2Hex(senderAddr),
		SimplexReceiver:         mainchain.Addr2Hex(receiverAddr),
		SignedSimplexStateBytes: signedSimplexStateBytes,
		DisputeTimeout:          disputeTimeout,
	}
}

// implement fmt.Stringer
func (r Request) String() string {
	return strings.TrimSpace(fmt.Sprintf(`SeqNum: %d, SimplexSender: %s, SimplexReceiver: %s, DisputeTimeout: %d, TriggerTxHash: %s, TriggerTxBlkNum: %d, GuardTxHash: %s, GuardTxBlkNum: %d, GuardSender: %s`,
		r.SeqNum, r.SimplexSender, r.SimplexReceiver, r.DisputeTimeout,
		r.TriggerTxHash, r.TriggerTxBlkNum, r.GuardTxHash, r.GuardTxBlkNum, r.GuardSender))
}

type InitRequest struct {
	SignedSimplexStateBytes []byte `json:"signedSimplexStateBytes"`
	SimplexReceiverSig      []byte `json:"simplexReceiverSig"`
	DisputeTimeout          uint64 `json:"disputeTimeout"`
}

func NewInitRequest(signedSimplexStateBytes, simplexReceiverSig []byte, disputeTimeout uint64) *InitRequest {
	return &InitRequest{
		SignedSimplexStateBytes: signedSimplexStateBytes,
		SimplexReceiverSig:      simplexReceiverSig,
		DisputeTimeout:          disputeTimeout,
	}
}
