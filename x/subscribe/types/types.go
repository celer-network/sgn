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
	ChannelId               []byte   `json:"channelId"`
	SignedSimplexStateBytes []byte   `json:"signedSimplexStateBytes"`
}

// Returns a new Number with the minprice as the price
func NewRequest(channelId []byte, seqNum uint64, peerAddresses []string, peerFromIndex uint) Request {
	return Request{
		SeqNum:        seqNum,
		ChannelId:     channelId,
		PeerAddresses: peerAddresses,
		PeerFromIndex: peerFromIndex,
	}
}

// implement fmt.Stringer
func (r Request) String() string {
	return strings.TrimSpace(fmt.Sprintf(`SeqNum: %d, ChannelId: %x, PeerAddresses: %x, PeerFromIndex: %d, SignedSimplexStateBytes: %x`,
		r.SeqNum, r.ChannelId, r.PeerAddresses, r.PeerFromIndex, r.SignedSimplexStateBytes))
}
