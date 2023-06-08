package events

import (
	"errors"
	"sync"
)

var ErrHandlerAlreadyRegistered = errors.New("handler already registered")

type Dispatcher struct {
	handlers map[string][]EventHandler
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		handlers: make(map[string][]EventHandler),
	}
}

func (dispatcher *Dispatcher) Register(eventName string, evHandler EventHandler) error {
	if _, ok := dispatcher.handlers[eventName]; ok {
		for _, handler := range dispatcher.handlers[eventName] {
			if handler == evHandler {
				return ErrHandlerAlreadyRegistered
			}
		}
	}
	dispatcher.handlers[eventName] = append(dispatcher.handlers[eventName], evHandler)
	return nil
}

func (dispatcher *Dispatcher) Dispatch(event Event) error {
	if handlers, ok := dispatcher.handlers[event.GetName()]; ok {
		wg := &sync.WaitGroup{}
		for _, handler := range handlers {
			wg.Add(1)
			go handler.Handle(event, wg)
		}
		wg.Wait()
	}
	return nil
}

func (dispatcher *Dispatcher) Remove(eventName string, handler EventHandler) error {
	if _, ok := dispatcher.handlers[eventName]; ok {
		for i, h := range dispatcher.handlers[eventName] {
			if h == handler {
				dispatcher.handlers[eventName] = append(
					dispatcher.handlers[eventName][:i],
					dispatcher.handlers[eventName][i+1:]...)
				return nil
			}
		}
	}
	return nil
}

func (dispatcher *Dispatcher) Has(eventName string, handler EventHandler) bool {
	if _, ok := dispatcher.handlers[eventName]; ok {
		for _, h := range dispatcher.handlers[eventName] {
			if h == handler {
				return true
			}
		}
	}
	return false
}

func (dispatcher *Dispatcher) Clear() {
	dispatcher.handlers = make(map[string][]EventHandler)
}
