package monitor

import (
	"encoding/json"

	"github.com/celer-network/sgn/mainchain"
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
	Name EventName `json:"name"`
	Log  types.Log `json:"log"`
}

func NewEvent(name EventName, l types.Log) Event {
	return Event{
		Name: name,
		Log:  l,
	}
}

func NewEventFromBytes(input []byte) Event {
	event := Event{}
	event.MustUnMarshal(input)
	return event
}

// Marshal event into json bytes
func (e Event) MustMarshal() []byte {
	res, err := json.Marshal(&e)
	if err != nil {
		panic(err)
	}

	return res
}

// Unmarshal json bytes to event
func (e *Event) MustUnMarshal(input []byte) {
	err := json.Unmarshal(input, e)
	if err != nil {
		panic(err)
	}
}

func (e Event) ParseEvent(ethClient *mainchain.EthClient) (res interface{}) {
	var err error
	switch e.Name {
	case InitializeCandidate:
		res, err = ethClient.Guard.ParseInitializeCandidate(e.Log)
	default:
		panic("Unsupported event")
	}

	if err != nil {
		panic(err)
	}

	switch tmp := res.(type) {
	case *mainchain.GuardInitializeCandidate:
		tmp.Raw = e.Log
		res = tmp
	case *mainchain.CelerLedgerIntendSettle:
		tmp.Raw = e.Log
		res = tmp
	}

	return
}

type PenaltyEvent struct {
	nonce uint64
}

func NewPenaltyEvent(nonce uint64) PenaltyEvent {
	return PenaltyEvent{
		nonce: nonce,
	}
}

func NewPenaltyEventFromBytes(input []byte) PenaltyEvent {
	event := PenaltyEvent{}
	err := json.Unmarshal(input, event)
	if err != nil {
		panic(err)
	}
	return event
}

// Marshal event into json bytes
func (e PenaltyEvent) MustMarshal() []byte {
	res, err := json.Marshal(&e)
	if err != nil {
		panic(err)
	}

	return res
}
