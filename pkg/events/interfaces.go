package events

import (
	"sync"
	"time"
)

type Event interface {
	GetName() string
	GetDateTime() time.Time
	GetPayload() any
}

type EventHandler interface {
	Handle(e Event, wg *sync.WaitGroup)
}

type EventDispatcher interface {
	Register(eventName string, h EventHandler) error
	Dispatch(e Event) error
	Remove(eventName string, h EventHandler) error
	Has(eventName string, h EventHandler) bool
	Clear() error
}
