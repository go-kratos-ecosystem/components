package pubsub

type Topic[T any] struct {
	subscriptions []*Subscription[T]
}

func NewTopic[T any]() *Topic[T] {
	return &Topic[T]{
		subscriptions: make([]*Subscription[T], 0),
	}
}

func (t *Topic[T]) Subscribe(handler func(msg T) error) *Subscription[T] {
	sub := &Subscription[T]{Topic: t, handler: handler}
	t.subscriptions = append(t.subscriptions, sub)
	return sub
}

type publishOptions struct {
	skipErrors bool
	async      bool
}

type PublishOption func(*publishOptions)

func PublishSkipErrors() PublishOption {
	return func(o *publishOptions) {
		o.skipErrors = true
	}
}

func PublishAsync() PublishOption {
	return func(o *publishOptions) {
		o.async = true
	}
}

func (t *Topic[T]) Publish(msg T, opts ...PublishOption) error {
	o := &publishOptions{}
	for _, opt := range opts {
		opt(o)
	}

	for _, sub := range t.subscriptions {
		if o.async {
			go func() {
				_ = sub.handler(msg) //nolint:govet
			}()
		} else {
			if err := sub.handler(msg); err != nil {
				if o.skipErrors {
					continue
				}

				return err
			}
		}
	}

	return nil
}

func (t *Topic[T]) PublishAsync(msg T, opts ...PublishOption) error {
	return t.Publish(msg, append(opts, PublishAsync())...)
}
