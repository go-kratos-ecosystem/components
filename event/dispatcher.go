package event

import (
	"sync"
)

type Event interface {
	Event() any
}

type Listener interface {
	Listen() []Event
	Handle(event Event)
}

type Dispatcher struct {
	listeners map[any][]Listener
	mu        sync.RWMutex
	recovery  func(err any, listener Listener, event Event)
}

type Option func(*Dispatcher)

func WithRecovery(recovery func(err any, listener Listener, event Event)) Option {
	return func(d *Dispatcher) {
		d.recovery = recovery
	}
}

func NewDispatcher(opts ...Option) *Dispatcher {
	d := &Dispatcher{
		listeners: make(map[any][]Listener),
	}

	for _, opt := range opts {
		opt(d)
	}

	return d
}

func (d *Dispatcher) AddListener(listener ...Listener) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for _, l := range listener {
		for _, event := range l.Listen() {
			event := event.Event()
			if _, ok := d.listeners[event]; !ok {
				d.listeners[event] = make([]Listener, 0)
			}
			d.listeners[event] = append(d.listeners[event], l)
		}
	}
}

func (d *Dispatcher) Dispatch(event Event) {
	if listeners, ok := d.listeners[event.Event()]; ok {
		for _, listener := range listeners {
			d.handle(listener, event)
		}
	}
}

func (d *Dispatcher) DispatchAsync(event Event) {
	if listeners, ok := d.listeners[event.Event()]; ok {
		for _, listener := range listeners {
			go d.handle(listener, event)
		}
	}
}

func (d *Dispatcher) handle(listener Listener, event Event) {
	if d.recovery != nil {
		defer func() {
			if err := recover(); err != nil {
				d.recovery(err, listener, event)
			}
		}()
	}
	listener.Handle(event)
}
