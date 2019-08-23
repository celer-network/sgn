package types

import (
	"fmt"
	"strings"
)

type Guardian struct {
	Balance uint64 `json:"balance"`
}

// Returns a new Number with the minprice as the price
func NewGuardian() Guardian {
	return Guardian{}
}

// implement fmt.Stringer
func (g Guardian) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Balance: %d`, g.Balance))
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
