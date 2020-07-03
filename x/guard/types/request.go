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
	GuardPending            bool             `json:"guardPending"`
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
	out := fmt.Sprintf(`SeqNum: %d, SimplexSender: %s, SimplexReceiver: %s, DisputeTimeout: %d`,
		r.SeqNum, r.SimplexSender, r.SimplexReceiver, r.DisputeTimeout)
	if r.TriggerTxHash != "" {
		out += fmt.Sprintf(`, TriggerTxHash: %s, TriggerTxBlkNum: %d`, r.TriggerTxHash, r.TriggerTxBlkNum)
	}
	if r.GuardTxHash != "" {
		out += fmt.Sprintf(`, GuardTxHash: %s, GuardTxBlkNum: %d, GuardSender: %s`, r.GuardTxHash, r.GuardTxBlkNum, r.GuardSender)
	}
	return strings.TrimSpace(out)
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

type GuardTrigger struct {
	ChannelId       []byte `json:"channelId"`
	SimplexReceiver string `json:"simplexReceiver"`
	TriggerTxHash   string `json:"triggerTxHash"`
	TriggerTxBlkNum uint64 `json:"triggerTxBlkNum"`
}

func (gt GuardTrigger) String() string {
	return strings.TrimSpace(fmt.Sprintf(`ChannelId: %x, SimplexReceiver: %s, TriggerTxHash: %s, TriggerTxBlkNum: %d`,
		gt.ChannelId, gt.SimplexReceiver, gt.TriggerTxHash, gt.TriggerTxBlkNum))
}

func NewGuardTrigger(
	channelId mainchain.CidType,
	simplexReceiver mainchain.Addr,
	triggerTxHash mainchain.HashType,
	triggerTxBlkNum uint64) *GuardTrigger {
	return &GuardTrigger{
		ChannelId:       channelId.Bytes(),
		SimplexReceiver: mainchain.Addr2Hex(simplexReceiver),
		TriggerTxHash:   triggerTxHash.Hex(),
		TriggerTxBlkNum: triggerTxBlkNum,
	}
}

type GuardProof struct {
	ChannelId       []byte `json:"channelId"`
	SimplexReceiver string `json:"simplexReceiver"`
	GuardTxHash     string `json:"guardTxHash"`
	GuardTxBlkNum   uint64 `json:"guardTxBlkNum"`
	GuardSender     string `json:"guardSender"`
}

func (gp GuardProof) String() string {
	return strings.TrimSpace(fmt.Sprintf(`ChannelId: %x, SimplexReceiver: %s, GuardTxHash: %s, GuardTxBlkNum: %d, GuardSender: %s`,
		gp.ChannelId, gp.SimplexReceiver, gp.GuardTxHash, gp.GuardTxBlkNum, gp.GuardSender))
}

func NewGuardProof(
	channelId mainchain.CidType,
	simplexReceiver mainchain.Addr,
	guardTxHash mainchain.HashType,
	guardTxBlkNum uint64,
	guardSender mainchain.Addr) *GuardProof {
	return &GuardProof{
		ChannelId:       channelId.Bytes(),
		SimplexReceiver: mainchain.Addr2Hex(simplexReceiver),
		GuardTxHash:     guardTxHash.Hex(),
		GuardTxBlkNum:   guardTxBlkNum,
		GuardSender:     mainchain.Addr2Hex(guardSender),
	}
}
