package event

import (
	"sync"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/go-kratos-ecosystem/components/v2/feature"
)

type Event string

func (e Event) String() string {
	return string(e)
}

type Listener interface {
	Listen() []Event
	Handle(event Event, data interface{})
}

type RecoveryHandler func(err interface{}, listener Listener, event Event, data interface{})

type Dispatcher struct {
	listeners map[Event][]Listener
	rw        sync.RWMutex

	recovery RecoveryHandler
}

type Option func(*Dispatcher)

func WithRecovery(handler RecoveryHandler) Option {
	return func(d *Dispatcher) {
		if handler != nil {
			d.recovery = handler
		}
	}
}

var DefaultRecovery RecoveryHandler = func(err interface{}, listener Listener, event Event, data interface{}) {
	log.Errorf("[Event] handler panic listener: %v, event: %s, data: %v, err: %v", listener, event, data, err)
}

func NewDispatcher(opts ...Option) *Dispatcher {
	d := &Dispatcher{
		listeners: make(map[Event][]Listener),
	}

	for _, opt := range opts {
		opt(d)
	}

	return d
}

func (d *Dispatcher) AddListener(listener ...Listener) {
	d.rw.Lock()
	defer d.rw.Unlock()

	for _, l := range listener {
		for _, event := range l.Listen() {
			if _, ok := d.listeners[event]; !ok {
				d.listeners[event] = make([]Listener, 0)
			}

			d.listeners[event] = append(d.listeners[event], l)
		}
	}
}

func (d *Dispatcher) Dispatch(event Event, data interface{}) {
	d.rw.RLock()
	defer d.rw.RUnlock()

	if listeners, ok := d.listeners[event]; ok {
		for _, listener := range listeners {
			// if support Asyncable
			if l, ok := listener.(feature.Asyncable); ok && l.Async() {
				go d.handle(listener, event, data)
			} else {
				d.handle(listener, event, data)
			}
		}
	}
}

func (d *Dispatcher) handle(listener Listener, event Event, data interface{}) {
	if d.recovery != nil {
		defer func() {
			if err := recover(); err != nil {
				d.recovery(err, listener, event, data)
			}
		}()
	}

	listener.Handle(event, data)
}
