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
	OwnerSig                []byte           `json:"ownerSig"`
	RequestGuards           []sdk.AccAddress `json:"requestGuards"`
	TriggerTxHash           string           `json:"triggerTxHash"`
	TriggerTxBlkNum         uint64           `json:"triggerTxBlkNum"`
	GuardTxHash             string           `json:"guardTxHash"`
	GuardTxBlkNum           uint64           `json:"guardTxBlkNum"`
	GuardSender             string           `json:"guardSender"`
}

func NewRequest(channelId []byte, seqNum uint64, peerAddresses []string, peerFromIndex uint8) Request {
	return Request{
		ChannelId:     channelId,
		SeqNum:        seqNum,
		PeerAddresses: peerAddresses,
		PeerFromIndex: peerFromIndex,
	}
}

func (r Request) GetOwnerAddress() string {
	if r.PeerFromIndex == 0 {
		return r.PeerAddresses[1]
	}

	return r.PeerAddresses[0]
}

func (r Request) GetPeerAddress() string {
	return r.PeerAddresses[r.PeerFromIndex]
}

// implement fmt.Stringer
func (r Request) String() string {
	return strings.TrimSpace(fmt.Sprintf(`SeqNum: %d, PeerAddresses: %s, PeerFromIndex: %d, DisputeTimeout: %d, TriggerTxHash: %s, TriggerTxBlkNum: %d, GuardTxHash: %s, GuardTxBlkNum: %d, GuardSender: %s`,
		r.SeqNum, r.PeerAddresses, r.PeerFromIndex, r.DisputeTimeout,
		r.TriggerTxHash, r.TriggerTxBlkNum, r.GuardTxHash, r.GuardTxBlkNum, r.GuardSender))
}
