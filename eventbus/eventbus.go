package eventbus

import (
	"context"
	"errors"
	"sync"
)

var ErrListenerNotFound = errors.New("[event_bus] listener not found")

type Handler[T any] interface {
	Handle(ctx context.Context, msg T) error
}

type HandlerFunc[T any] func(ctx context.Context, msg T) error

func (f HandlerFunc[T]) Handle(ctx context.Context, msg T) error {
	return f(ctx, msg)
}

type Event[T any] struct {
	listeners []*Listener[T]
	mu        sync.RWMutex
}

type Option[T any] func(*Event[T])

func NewEvent[T any]() *Event[T] {
	return &Event[T]{
		listeners: make([]*Listener[T], 0),
	}
}

func (t *Event[T]) On(handler Handler[T]) *Listener[T] {
	listener := newListener(t, handler)

	t.mu.Lock()
	defer t.mu.Unlock()

	t.listeners = append(t.listeners, listener)
	return listener
}

func (t *Event[T]) Off(listener *Listener[T]) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	for i, s := range t.listeners {
		if s == listener {
			t.listeners = append(t.listeners[:i], t.listeners[i+1:]...)
			return nil
		}
	}

	return ErrListenerNotFound
}

type emitOptions struct {
	async      bool
	skipErrors bool
}

type EmitOption func(*emitOptions)

func WithEmitAsync() EmitOption {
	return func(p *emitOptions) {
		p.async = true
	}
}

func WithEmitSkipErrors() EmitOption {
	return func(p *emitOptions) {
		p.skipErrors = true
	}
}

func (t *Event[T]) Emit(ctx context.Context, msg T, opts ...EmitOption) error {
	p := &emitOptions{}
	for _, opt := range opts {
		opt(p)
	}

	t.mu.RLock()
	defer t.mu.RUnlock()

	for _, listener := range t.listeners {
		if listener == nil || listener.handler == nil {
			continue
		}

		if p.async {
			go func(listener *Listener[T]) {
				_ = listener.handler.Handle(ctx, msg)
			}(listener)
		} else {
			if err := listener.handler.Handle(ctx, msg); err != nil {
				if !p.skipErrors {
					return err
				}
			}
		}
	}

	return nil
}

func (t *Event[T]) EmitAsync(ctx context.Context, msg T, opts ...EmitOption) error {
	return t.Emit(ctx, msg, append(opts, WithEmitAsync())...)
}

func (t *Event[T]) OffAll() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.listeners = make([]*Listener[T], 0)
}

func (t *Event[T]) Listeners() []*Listener[T] {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.listeners
}

func (t *Event[T]) ListenersCount() int {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return len(t.listeners)
}

type Listener[T any] struct {
	topic   *Event[T]
	handler Handler[T]
}

func newListener[T any](topic *Event[T], handler Handler[T]) *Listener[T] {
	return &Listener[T]{
		topic:   topic,
		handler: handler,
	}
}

func (s *Listener[T]) Off() error {
	return s.topic.Off(s)
}
