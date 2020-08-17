package monitor

import (
	"encoding/json"
	"strings"

	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/viper"
)

type EventName string

const (
	UpdateSidechainAddr   EventName = "UpdateSidechainAddr"
	ConfirmParamProposal  EventName = "ConfirmParamProposal"
	Delegate              EventName = "Delegate"
	CandidateUnbonded     EventName = "CandidateUnbonded"
	ValidatorChange       EventName = "ValidatorChange"
	UpdateCommissionRate  EventName = "UpdateCommissionRate"
	IntendSettle          EventName = "IntendSettle"
	IntendWithdraw        EventName = "IntendWithdraw"
	IntendWithdrawDpos    EventName = "IntendWithdrawDpos"
	IntendWithdrawChannel EventName = "IntendWithdrawChannel"
)

// Wrapper for ethereum Event
type EventWrapper struct {
	Name EventName `json:"name"`
	Log  types.Log `json:"log"`
}

func NewEvent(name EventName, l types.Log) *EventWrapper {
	return &EventWrapper{
		Name: name,
		Log:  l,
	}
}

func NewEventFromBytes(input []byte) *EventWrapper {
	event := &EventWrapper{}
	event.MustUnMarshal(input)
	return event
}

// Marshal event into json bytes
func (e *EventWrapper) MustMarshal() []byte {
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

func (e *EventWrapper) ParseEvent(ethClient *mainchain.EthClient) interface{} {
	var res interface{}
	var err error
	switch e.Name {
	case UpdateSidechainAddr:
		res, err = ethClient.SGN.ParseUpdateSidechainAddr(e.Log)
	case ConfirmParamProposal:
		res, err = ethClient.DPoS.ParseConfirmParamProposal(e.Log)
	case Delegate:
		res, err = ethClient.DPoS.ParseDelegate(e.Log)
	case CandidateUnbonded:
		res, err = ethClient.DPoS.ParseCandidateUnbonded(e.Log)
	case ValidatorChange:
		res, err = ethClient.DPoS.ParseValidatorChange(e.Log)
	case UpdateCommissionRate:
		res, err = ethClient.DPoS.ParseUpdateCommissionRate(e.Log)
	case IntendWithdrawDpos:
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

	return res
}

func eventCheckInterval(name EventName) uint64 {
	m := viper.GetStringMap(common.FlagEthCheckInterval)
	if m[string(name)] != nil {
		return uint64(m[string(name)].(float64))
	} else if m[strings.ToLower(string(name))] != nil {
		return uint64(m[strings.ToLower(string(name))].(float64))
	}
	return 0
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
