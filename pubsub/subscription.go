package pubsub

type Subscription[T any] struct {
	Topic   *Topic[T]
	handler func(msg T) error
}

func (s *Subscription[T]) Unsubscribe() {
	for i, sub := range s.Topic.subscriptions {
		if sub == s {
			s.Topic.subscriptions = append(s.Topic.subscriptions[:i], s.Topic.subscriptions[i+1:]...)
			return
		}
	}
}
