package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Subscription struct {
	Deposit     sdk.Int `json:"deposit"`
	Spend       sdk.Int `json:"spend"`
	Subscribing bool    `json:"subscribing"`
}

// Returns a new Number with the minprice as the price
func NewSubscription(deposit sdk.Int) Subscription {
	return Subscription{
		Deposit:     deposit,
		Spend:       sdk.NewInt(0),
		Subscribing: false,
	}
}

// implement fmt.Stringer
func (s Subscription) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Deposit: %v, Spend: %v, Subscribing: %b`, s.Deposit, s.Spend, s.Subscribing))
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
