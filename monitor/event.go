package monitor

import (
	"github.com/ethereum/go-ethereum/core/types"
)

type Event struct {
	event interface{}
	log   types.Log
}

func NewEvent(event interface{}, log types.Log) Event {
	return Event{
		event: event,
		log:   log,
	}
}
