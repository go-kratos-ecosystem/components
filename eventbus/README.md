# EventBus

## Usage

```go
package main

import (
	"context"
	"fmt"

	"github.com/go-kratos-ecosystem/components/v2/eventbus"
)

type IntListener[T int] struct{}

var _ eventbus.Handler[int] = (*IntListener[int])(nil)

func NewIntListener() *IntListener[int] {
	return &IntListener[int]{}
}

func (s *IntListener[T]) Handle(ctx context.Context, msg int) error {
	fmt.Println("IntListener", msg)
	return nil
}

type Event struct {
	ID int
}

type EventListener[T Event] struct{}

var _ eventbus.Handler[Event] = (*EventListener[Event])(nil)

func NewEventListener() *EventListener[Event] {
	return &EventListener[Event]{}
}

func (s *EventListener[T]) Handle(ctx context.Context, msg Event) error {
	fmt.Println("EventListener", msg.ID)
	return nil
}

func main() {
	// basic type
	event1 := eventbus.NewEvent[int]()
	listener1 := event1.On(eventbus.HandlerFunc[int](func(_ context.Context, msg int) error {
		fmt.Println("HandlerFunc", msg)
		return nil
	}))
	event1.On(NewIntListener())

	_ = event1.Emit(context.Background(), 1)
	// Output:
	// HandlerFunc 1
	// IntListener 1

	// struct type
	event2 := eventbus.NewEvent[Event]()
	event2.On(eventbus.HandlerFunc[Event](func(_ context.Context, msg Event) error {
		fmt.Println("HandlerFunc", msg.ID)
		return nil
	}))
	event2.On(NewEventListener())

	_ = event2.Emit(context.Background(), Event{ID: 2})
	// Output:
	// HandlerFunc 2
	// EventListener 2

	// off
	_ = listener1.Off()

	// async
	_ = event1.Emit(context.Background(), 3, eventbus.WithEmitAsync())
	_ = event1.EmitAsync(context.Background(), 4)

	// skip errors(only sync)
	_ = event1.Emit(context.Background(), 5, eventbus.WithEmitSkipErrors())
}

```