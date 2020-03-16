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
type EventWrapper struct {
	Name EventName `json:"name"`
	Log  types.Log `json:"log"`
}

func NewEvent(name EventName, l types.Log) EventWrapper {
	return EventWrapper{
		Name: name,
		Log:  l,
	}
}

func NewEventFromBytes(input []byte) EventWrapper {
	event := EventWrapper{}
	event.MustUnMarshal(input)
	return event
}

// Marshal event into json bytes
func (e EventWrapper) MustMarshal() []byte {
	res, err := json.Marshal(&e)
	if err != nil {
		panic(err)
	}

	return res
}

// Unmarshal json bytes to event
func (e *EventWrapper) MustUnMarshal(input []byte) {
	err := json.Unmarshal(input, e)
	if err != nil {
		panic(err)
	}
}

func (e EventWrapper) ParseEvent(ethClient *mainchain.EthClient) (res interface{}) {
	var err error
	switch e.Name {
	case InitializeCandidate:
		res, err = ethClient.Guard.ParseInitializeCandidate(e.Log)
	case Delegate:
		res, err = ethClient.Guard.ParseDelegate(e.Log)
	case ValidatorChange:
		res, err = ethClient.Guard.ParseValidatorChange(e.Log)
	case IntendWithdraw:
		res, err = ethClient.Guard.ParseIntendWithdraw(e.Log)
	case IntendSettle:
		res, err = ethClient.Ledger.ParseIntendSettle(e.Log)
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
	Nonce uint64
}

func NewPenaltyEvent(nonce uint64) PenaltyEvent {
	return PenaltyEvent{
		Nonce: nonce,
	}
}

func NewPenaltyEventFromBytes(input []byte) PenaltyEvent {
	event := PenaltyEvent{}
	event.MustUnMarshal(input)
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

// Unmarshal json bytes to penalty event
func (e *PenaltyEvent) MustUnMarshal(input []byte) {
	err := json.Unmarshal(input, e)
	if err != nil {
		panic(err)
	}
}
