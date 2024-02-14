package event

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-kratos-ecosystem/components/v2/feature"
)

type result struct {
	event Event
	data  interface{}
	err   interface{}
}

var recv = make(chan result, 1)

type testListener struct {
	feature.AsyncFeature
}

func newTestListener() *testListener {
	return &testListener{}
}

func (l *testListener) Listen() []Event {
	return []Event{
		"test",
		"test2",
	}
}

func (l *testListener) Handle(event Event, data interface{}) {
	if s, ok := data.(string); ok {
		recv <- result{
			event: event,
			data:  s,
		}
	} else {
		panic("invalid data")
	}
}

type test2Listener struct{}

func (l *test2Listener) Listen() []Event {
	return []Event{
		"test3",
	}
}

func (l *test2Listener) Handle(event Event, data interface{}) {
	recv <- result{
		event: event,
		data:  data,
	}
}

func TestDispatcher(t *testing.T) {
	var (
		d = NewDispatcher(
			WithRecovery(func(err interface{}, _ Listener, event Event, data interface{}) {
				recv <- result{
					event: event,
					data:  data,
					err:   err,
				}
			}),
		)
		l = newTestListener()
	)

	d.AddListener(l, &test2Listener{})
	assert.True(t, l.Async())

	d.Dispatch("test", "test data")
	r1 := <-recv
	assert.Equal(t, Event("test"), r1.event)
	assert.Equal(t, "test data", r1.data)

	d.Dispatch("test2", "test2 data")
	r2 := <-recv
	assert.Equal(t, Event("test2"), r2.event)
	assert.Equal(t, "test2 data", r2.data)

	d.Dispatch("test", 111)
	r3 := <-recv
	assert.Equal(t, Event("test"), r3.event)
	assert.Equal(t, 111, r3.data)
	assert.Equal(t, "invalid data", r3.err)

	d.Dispatch("test3", "test3 data")
	r4 := <-recv
	assert.Equal(t, Event("test3"), r4.event)
	assert.Equal(t, "test3 data", r4.data)
}
