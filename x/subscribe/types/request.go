package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Request struct {
	SeqNum                  uint64           `json:"seqNum"`
	PeerAddresses           []string         `json:"peerAddresses"`
	PeerFromIndex           uint8            `json:"peerFromIndex"`
	DisputeTimeout          uint64           `json:"disputeTimeout"`
	RequestGuards           []sdk.AccAddress `json:"requestGuards"`
	SignedSimplexStateBytes []byte           `json:"signedSimplexStateBytes"`
	TriggerTxHash           string           `json:"triggerTxHash"`
	GuardTxHash             string           `json:"guardTxHash"`
}

func NewRequest(seqNum uint64, peerAddresses []string, peerFromIndex uint8, disputeTimeout uint64, requestGuards []sdk.AccAddress) Request {
	return Request{
		SeqNum:         seqNum,
		PeerAddresses:  peerAddresses,
		PeerFromIndex:  peerFromIndex,
		DisputeTimeout: disputeTimeout,
		RequestGuards:  requestGuards,
	}
}

// implement fmt.Stringer
func (r Request) String() string {
	return strings.TrimSpace(fmt.Sprintf(`SeqNum: %d, PeerAddresses: %s, PeerFromIndex: %d, SignedSimplexStateBytes: %x, TriggerTxHash: %s, GuardTxHash: %s`,
		r.SeqNum, r.PeerAddresses, r.PeerFromIndex, r.SignedSimplexStateBytes, r.TriggerTxHash, r.GuardTxHash))
}
