package monitor

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/core/types"
)

type EventName string

const (
	InitializeCandidate EventName = "InitializeCandidate"
	Delegate            EventName = "Delegate"
	ValidatorChange     EventName = "ValidatorChange"
	IntendWithdraw      EventName = "IntendWithdraw"
	IntendSettle        EventName = "IntendSettle"
)

// Wrapper for ethereum Event
type Event struct {
	name EventName `json:"name"`
	log  types.Log `json:"log"`
}

func NewEvent(name EventName, l types.Log) Event {
	return Event{
		name: name,
		log:  l,
	}
}

func (e Event) MustMarshal() []byte {
	res, err := json.Marshal(&e)
	if err != nil {
		panic(err)
	}

	return res
}
func (e *Event) MustUnMarshal(input []byte) {
	err := json.Unmarshal(input, e)
	if err != nil {
		panic(err)
	}
}

type PenaltyEvent struct {
	nonce uint64
}

func NewPenaltyEvent(nonce uint64) PenaltyEvent {
	return PenaltyEvent{
		nonce: nonce,
	}
}
