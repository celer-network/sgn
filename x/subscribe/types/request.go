package types

import (
	"fmt"
	"strings"
)

type Request struct {
	SeqNum                  uint64   `json:"seqNum"`
	PeerAddresses           []string `json:"peerAddresses"`
	PeerFromIndex           uint     `json:"peerFromIndex"`
	SignedSimplexStateBytes []byte   `json:"signedSimplexStateBytes"`
	TxHash                  string   `json:"txHash"`
}

func NewRequest(seqNum uint64, peerAddresses []string, peerFromIndex uint) Request {
	return Request{
		SeqNum:        seqNum,
		PeerAddresses: peerAddresses,
		PeerFromIndex: peerFromIndex,
	}
}

// implement fmt.Stringer
func (r Request) String() string {
	return strings.TrimSpace(fmt.Sprintf(`SeqNum: %d, PeerAddresses: %x, PeerFromIndex: %d, SignedSimplexStateBytes: %x, TxHash: %s`,
		r.SeqNum, r.PeerAddresses, r.PeerFromIndex, r.SignedSimplexStateBytes, r.TxHash))
}
