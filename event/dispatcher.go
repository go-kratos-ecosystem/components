package event

import (
	"sync"

	"github.com/go-kratos-ecosystem/components/v2/features"
)

type Event interface {
	Name() string
}

type Listener interface {
	Listen() []Event
	Handle(event Event)
}

type Dispatcher struct {
	listeners map[string][]Listener
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
		listeners: make(map[string][]Listener),
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
			name := event.Name()
			if _, ok := d.listeners[name]; !ok {
				d.listeners[name] = make([]Listener, 0)
			}
			d.listeners[name] = append(d.listeners[name], l)
		}
	}
}

func (d *Dispatcher) Dispatch(event Event) {
	if listeners, ok := d.listeners[event.Name()]; ok {
		for _, listener := range listeners {
			if l, ok := listener.(features.Asyncable); ok && l.Async() {
				go d.handle(listener, event)
				continue
			}
			d.handle(listener, event)
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
