package eventbus

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var ctx = context.Background()

func TestEventBus_Func(t *testing.T) {
	topic := NewTopic[int]()
	ch := make(chan int, 1)

	assert.Equal(t, topic.SubscriptionsCount(), 0)
	sub := topic.Subscribe(HandlerFunc[int](func(_ context.Context, msg int) error {
		ch <- msg
		return nil
	}))
	assert.Len(t, topic.Subscriptions(), 1)
	assert.Equal(t, 1, topic.SubscriptionsCount())
	assert.NoError(t, topic.Publish(ctx, 1))
	assert.Equal(t, 1, <-ch)

	assert.NoError(t, sub.Unsubscribe())
	assert.Len(t, topic.Subscriptions(), 0)
	assert.Equal(t, 0, topic.SubscriptionsCount())
	assert.NoError(t, topic.Publish(ctx, 2))
	assert.Len(t, ch, 0)

	topic.UnsubscribeAll()
	assert.NoError(t, topic.Publish(ctx, 3))
	assert.Len(t, ch, 0)
}

type Event struct {
	ID int
}

func TestEventBus_Struct(t *testing.T) {
	topic := NewTopic[Event]()
	ch := make(chan int, 1)

	sub := topic.Subscribe(HandlerFunc[Event](func(_ context.Context, msg Event) error {
		ch <- msg.ID
		return nil
	}))
	assert.NoError(t, topic.Publish(ctx, Event{ID: 1}))
	assert.Equal(t, 1, <-ch)

	assert.NoError(t, sub.Unsubscribe())
	assert.NoError(t, topic.Publish(ctx, Event{ID: 2}))
	assert.Len(t, ch, 0)

	assert.ErrorIs(t, sub.Unsubscribe(), ErrSubscriptionNotFound)
}

func TestEventBus_SkipErrors(t *testing.T) {
	topic := NewTopic[int]()
	ch := make(chan int, 1)

	_ = topic.Subscribe(HandlerFunc[int](func(context.Context, int) error {
		return assert.AnError
	}))
	_ = topic.Subscribe(HandlerFunc[int](func(_ context.Context, msg int) error {
		ch <- msg
		return nil
	}))

	// no skip errors
	assert.ErrorIs(t, topic.Publish(ctx, 1), assert.AnError)
	assert.Len(t, ch, 0)

	// skip errors
	assert.NoError(t, topic.Publish(ctx, 2, WithPublishSkipErrors()))
	assert.Equal(t, 2, <-ch)
}

func TestEventBus_Async(t *testing.T) {
	topic := NewTopic[int]()
	ch := make(chan int, 2)
	wg := sync.WaitGroup{}

	_ = topic.Subscribe(HandlerFunc[int](func(_ context.Context, msg int) error {
		defer wg.Done()
		time.Sleep(time.Millisecond * 100)
		ch <- 10 * msg
		return nil
	}))
	_ = topic.Subscribe(HandlerFunc[int](func(_ context.Context, msg int) error {
		defer wg.Done()
		time.Sleep(time.Millisecond * 100)
		ch <- 100 * msg
		return nil
	}))

	// sync
	wg.Add(2)
	start := time.Now()
	assert.NoError(t, topic.Publish(ctx, 1))
	wg.Wait()
	assert.Equal(t, 10, <-ch)
	assert.Equal(t, 100, <-ch)
	assert.Len(t, ch, 0)
	assert.GreaterOrEqual(t, time.Since(start).Milliseconds(), int64(200))

	// async
	wg.Add(2)
	start = time.Now()
	assert.NoError(t, topic.Publish(ctx, 2, WithPublishAsync()))
	wg.Wait()
	assert.Len(t, ch, 2)
	assert.Contains(t, []int{20, 200}, <-ch)
	assert.Contains(t, []int{20, 200}, <-ch)
	assert.GreaterOrEqual(t, time.Since(start).Milliseconds(), int64(100))
	assert.Less(t, time.Since(start).Milliseconds(), int64(150))
}
