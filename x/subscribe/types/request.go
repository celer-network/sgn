package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Request struct {
	ChannelId               []byte           `json:"channelId"`
	SeqNum                  uint64           `json:"seqNum"`
	PeerAddresses           []string         `json:"peerAddresses"`
	PeerFromIndex           uint8            `json:"peerFromIndex"`
	DisputeTimeout          uint64           `json:"disputeTimeout"`
	SignedSimplexStateBytes []byte           `json:"signedSimplexStateBytes"`
	PeerToSig               []byte           `json:"peerToSig"`
	RequestGuards           []sdk.AccAddress `json:"requestGuards"`
	TriggerTxHash           string           `json:"triggerTxHash"`
	TriggerTxBlkNum         uint64           `json:"triggerTxBlkNum"`
	GuardTxHash             string           `json:"guardTxHash"`
	GuardTxBlkNum           uint64           `json:"guardTxBlkNum"`
	GuardSender             string           `json:"guardSender"`
}

func NewRequest(
	channelId []byte,
	seqNum uint64,
	peerAddresses []string,
	peerFromIndex uint8,
	signedSimplex []byte,
	peerToSig []byte) Request {
	return Request{
		ChannelId:               channelId,
		SeqNum:                  seqNum,
		PeerAddresses:           peerAddresses,
		PeerFromIndex:           peerFromIndex,
		SignedSimplexStateBytes: signedSimplex,
		PeerToSig:               peerToSig,
	}
}

func (r Request) GetPeerFromAddress() string {
	return r.PeerAddresses[r.PeerFromIndex]
}

func (r Request) GetPeerToAddress() string {
	return r.PeerAddresses[1-r.PeerFromIndex]
}

// implement fmt.Stringer
func (r Request) String() string {
	return strings.TrimSpace(fmt.Sprintf(`SeqNum: %d, PeerAddresses: %s, PeerFromIndex: %d, DisputeTimeout: %d, TriggerTxHash: %s, TriggerTxBlkNum: %d, GuardTxHash: %s, GuardTxBlkNum: %d, GuardSender: %s`,
		r.SeqNum, r.PeerAddresses, r.PeerFromIndex, r.DisputeTimeout,
		r.TriggerTxHash, r.TriggerTxBlkNum, r.GuardTxHash, r.GuardTxBlkNum, r.GuardSender))
}
