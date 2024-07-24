package eventbus

import (
	"context"
	"errors"
	"sync"
)

var ErrSubscriptionNotFound = errors.New("[event_bus] subscription not found")

type Handler[T any] interface {
	Handle(ctx context.Context, msg T) error
}

type HandlerFunc[T any] func(ctx context.Context, msg T) error

func (f HandlerFunc[T]) Handle(ctx context.Context, msg T) error {
	return f(ctx, msg)
}

type Topic[T any] struct {
	subscriptions []*Subscription[T]
	mu            sync.RWMutex
}

type Option[T any] func(*Topic[T])

func NewTopic[T any]() *Topic[T] {
	return &Topic[T]{
		subscriptions: make([]*Subscription[T], 0),
	}
}

func (t *Topic[T]) Subscribe(handler Handler[T]) *Subscription[T] {
	sub := newSubscription(t, handler)

	t.mu.Lock()
	defer t.mu.Unlock()

	t.subscriptions = append(t.subscriptions, sub)
	return sub
}

func (t *Topic[T]) Unsubscribe(sub *Subscription[T]) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	for i, s := range t.subscriptions {
		if s == sub {
			t.subscriptions = append(t.subscriptions[:i], t.subscriptions[i+1:]...)
			return nil
		}
	}

	return ErrSubscriptionNotFound
}

type publishOptions struct {
	async      bool
	skipErrors bool
}

type PublishOption func(*publishOptions)

func WithAsync() PublishOption {
	return func(p *publishOptions) {
		p.async = true
	}
}

func WithSkipErrors() PublishOption {
	return func(p *publishOptions) {
		p.skipErrors = true
	}
}

func (t *Topic[T]) Publish(ctx context.Context, msg T, opts ...PublishOption) error {
	p := &publishOptions{}
	for _, opt := range opts {
		opt(p)
	}

	t.mu.RLock()
	defer t.mu.RUnlock()

	for _, sub := range t.subscriptions {
		if sub == nil || sub.handler == nil {
			continue
		}

		if p.async {
			go func(sub *Subscription[T]) {
				_ = sub.handler.Handle(ctx, msg)
			}(sub)
		} else {
			if err := sub.handler.Handle(ctx, msg); err != nil {
				if !p.skipErrors {
					return err
				}
			}
		}
	}

	return nil
}

func (t *Topic[T]) PublishAsync(ctx context.Context, msg T, opts ...PublishOption) error {
	return t.Publish(ctx, msg, append(opts, WithAsync())...)
}

func (t *Topic[T]) UnsubscribeAll() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.subscriptions = make([]*Subscription[T], 0)
}

func (t *Topic[T]) Subscriptions() []*Subscription[T] {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.subscriptions
}

func (t *Topic[T]) SubscriptionsCount() int {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return len(t.subscriptions)
}

type Subscription[T any] struct {
	topic   *Topic[T]
	handler Handler[T]
}

func newSubscription[T any](topic *Topic[T], handler Handler[T]) *Subscription[T] {
	return &Subscription[T]{
		topic:   topic,
		handler: handler,
	}
}

func (s *Subscription[T]) Unsubscribe() error {
	return s.topic.Unsubscribe(s)
}
