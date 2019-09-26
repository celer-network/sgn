package monitor

import (
	"log"

	"github.com/ethereum/go-ethereum/core/types"
)

// Wrapper for ethereum Event
type Event struct {
	event interface{}
	log   types.Log
}

func NewEvent(event interface{}, l types.Log) Event {
	log.Printf("New event", event)

	return Event{
		event: event,
		log:   l,
	}
}
