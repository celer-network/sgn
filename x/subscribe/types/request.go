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
	RequestGuards           []sdk.AccAddress `json:"requestGuards"`
	SignedSimplexStateBytes []byte           `json:"signedSimplexStateBytes"`
	TriggerTxHash           string           `json:"triggerTxHash"`
	GuardTxHash             string           `json:"guardTxHash"`
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
	return strings.TrimSpace(fmt.Sprintf(`SeqNum: %d, PeerAddresses: %s, PeerFromIndex: %d, SignedSimplexStateBytes: %x, TriggerTxHash: %s, GuardTxHash: %s`,
		r.SeqNum, r.PeerAddresses, r.PeerFromIndex, r.SignedSimplexStateBytes, r.TriggerTxHash, r.GuardTxHash))
}
