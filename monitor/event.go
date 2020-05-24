package monitor

import (
	"encoding/json"

	"github.com/celer-network/sgn/mainchain"
	"github.com/ethereum/go-ethereum/core/types"
)

type EventName string

const (
	UpdateSidechainAddr   EventName = "UpdateSidechainAddr"
	ConfirmParamProposal  EventName = "ConfirmParamProposal"
	Delegate              EventName = "Delegate"
	CandidateUnbonded     EventName = "CandidateUnbonded"
	ValidatorChange       EventName = "ValidatorChange"
	IntendWithdrawSgn     EventName = "IntendWithdraw"
	IntendWithdrawChannel EventName = "IntendWithdrawChannel"
	IntendSettle          EventName = "IntendSettle"
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
	case UpdateSidechainAddr:
		res, err = ethClient.SGN.ParseUpdateSidechainAddr(e.Log)
	case ConfirmParamProposal:
		res, err = ethClient.DPoS.ParseConfirmParamProposal(e.Log)
	case Delegate:
		res, err = ethClient.DPoS.ParseDelegate(e.Log)
	case ValidatorChange:
		res, err = ethClient.DPoS.ParseValidatorChange(e.Log)
	case IntendWithdrawSgn:
		res, err = ethClient.DPoS.ParseIntendWithdraw(e.Log)
	case IntendWithdrawChannel:
		res, err = ethClient.Ledger.ParseIntendWithdraw(e.Log)
	case IntendSettle:
		res, err = ethClient.Ledger.ParseIntendSettle(e.Log)
	default:
		panic("Unsupported event")
	}

	if err != nil {
		panic(err)
	}

	switch tmp := res.(type) {
	case *mainchain.SGNUpdateSidechainAddr:
		tmp.Raw = e.Log
		res = tmp
	case *mainchain.CelerLedgerIntendSettle:
		tmp.Raw = e.Log
		res = tmp
	}

	return
}

type PenaltyEvent struct {
	Nonce      uint64 `json:"nonce"`
	RetryCount uint64 `json:"retryCount"`
}

func NewPenaltyEvent(nonce uint64) PenaltyEvent {
	return PenaltyEvent{
		Nonce:      nonce,
		RetryCount: 0,
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
