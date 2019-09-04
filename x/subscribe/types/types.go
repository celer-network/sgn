package types

import (
	"fmt"
	"strings"
)

type Subscription struct {
	Expiration uint64 `json:"expiration"`
}

// Returns a new Number with the minprice as the price
func NewSubscription(expiration uint64) Subscription {
	return Subscription{
		Expiration: expiration,
	}
}

// implement fmt.Stringer
func (s Subscription) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Expiration: %d`, s.Expiration))
}

type Request struct {
	SeqNum                  uint64   `json:"seqNum"`
	PeerAddresses           []string `json:"peerAddresses"`
	PeerFromIndex           uint     `json:"peerFromIndex"`
	SignedSimplexStateBytes []byte   `json:"signedSimplexStateBytes"`
	TxHash                  string   `json:"txHash"`
}

// Returns a new Number with the minprice as the price
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
