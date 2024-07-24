# EventBus

## Usage

```go
package main

import (
	"context"
	"fmt"

	"github.com/go-kratos-ecosystem/components/v2/eventbus"
)

type IntSubscriber[T int] struct{}

var _ eventbus.Handler[int] = (*IntSubscriber[int])(nil)

func NewIntSubscriber() *IntSubscriber[int] {
	return &IntSubscriber[int]{}
}

func (s *IntSubscriber[T]) Handle(ctx context.Context, msg int) error {
	fmt.Println("IntSubscriber", msg)
	return nil
}

type Event struct {
	ID int
}

type EventSubscriber[T Event] struct{}

var _ eventbus.Handler[Event] = (*EventSubscriber[Event])(nil)

func NewEventSubscriber() *EventSubscriber[Event] {
	return &EventSubscriber[Event]{}
}

func (s *EventSubscriber[T]) Handle(ctx context.Context, msg Event) error {
	fmt.Println("EventSubscriber", msg.ID)
	return nil
}

func main() {
	// basic type
	topic1 := eventbus.NewTopic[int]()
	sub1 := topic1.Subscribe(eventbus.HandlerFunc[int](func(_ context.Context, msg int) error {
		fmt.Println("HandlerFunc", msg)
		return nil
	}))
	topic1.Subscribe(NewIntSubscriber())

	_ = topic1.Publish(context.Background(), 1)
	// Output:
	// HandlerFunc 1
	// IntSubscriber 1

	// struct type
	topic2 := eventbus.NewTopic[Event]()
	topic2.Subscribe(eventbus.HandlerFunc[Event](func(_ context.Context, msg Event) error {
		fmt.Println("HandlerFunc", msg.ID)
		return nil
	}))
	topic2.Subscribe(NewEventSubscriber())

	_ = topic2.Publish(context.Background(), Event{ID: 2})
	// Output:
	// HandlerFunc 2
	// EventSubscriber 2

	// unsubscribe
	_ = sub1.Unsubscribe()

	// async
	_ = topic1.Publish(context.Background(), 3, eventbus.WithPublishAsync())
	_ = topic1.PublishAsync(context.Background(), 4)

	// skip errors(only sync)
	_ = topic1.Publish(context.Background(), 5, eventbus.WithPublishSkipErrors())
}
```