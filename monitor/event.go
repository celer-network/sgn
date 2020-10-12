package monitor

import (
	"encoding/json"

	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/iancoleman/strcase"
	"github.com/spf13/viper"
)

type EventName string

const (
	UpdateSidechainAddr    EventName = "UpdateSidechainAddr"
	AddSubscriptionBalance EventName = "AddSubscriptionBalance"
	ConfirmParamProposal   EventName = "ConfirmParamProposal"
	UpdateDelegatedStake   EventName = "UpdateDelegatedStake"
	CandidateUnbonded      EventName = "CandidateUnbonded"
	ValidatorChange        EventName = "ValidatorChange"
	UpdateCommissionRate   EventName = "UpdateCommissionRate"
	IntendSettle           EventName = "IntendSettle"
	IntendWithdraw         EventName = "IntendWithdraw"
	IntendWithdrawChannel  EventName = "IntendWithdrawChannel"
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
	case AddSubscriptionBalance:
		res, err = ethClient.SGN.ParseAddSubscriptionBalance(e.Log)
	case ConfirmParamProposal:
		res, err = ethClient.DPoS.ParseConfirmParamProposal(e.Log)
	case UpdateDelegatedStake:
		res, err = ethClient.DPoS.ParseUpdateDelegatedStake(e.Log)
	case CandidateUnbonded:
		res, err = ethClient.DPoS.ParseCandidateUnbonded(e.Log)
	case ValidatorChange:
		res, err = ethClient.DPoS.ParseValidatorChange(e.Log)
	case UpdateCommissionRate:
		res, err = ethClient.DPoS.ParseUpdateCommissionRate(e.Log)
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

func getEventCheckInterval(name EventName) uint64 {
	m := viper.GetStringMap(common.FlagEthCheckInterval)
	eventNameInConfig := strcase.ToSnake(string(name))
	if m[eventNameInConfig] != nil {
		return uint64(m[eventNameInConfig].(int64))
	}
	// If not specified, use the default value of 0
	return 0
}

type PenaltyEvent struct {
	Nonce      uint64 `json:"nonce"`
	RetryCount uint64 `json:"retry_count"`
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
