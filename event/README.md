# Event

## Usage

```go
package main

import (
	"fmt"
	"time"

	"github.com/go-kratos-ecosystem/components/v2/event"
	"github.com/go-kratos-ecosystem/components/v2/feature"
)

type testListener struct {
	feature.AsyncFeature  // async feature
}

func (l *testListener) Listen() []event.Event {
	return []event.Event{
		"test",
		"test2",
	}
}

func (l *testListener) Handle(e event.Event, data interface{}) {
	if s, ok := data.(string); ok {
		fmt.Println("event:", e, "data:", s)
	} else {
		panic("invalid data")
	}
}

type test2Listener struct{}

func (l *test2Listener) Listen() []event.Event {
	return []event.Event{
		"test",
	}
}

func (l *test2Listener) Handle(event event.Event, data interface{}) {
	fmt.Println("event:", event, "data:", data)
}

func main() {
	d := event.NewDispatcher(event.WithRecovery(func(err interface{}, listener event.Listener, event event.Event, data interface{}) { //nolint:lll
		fmt.Println("err:", err, "listener:", listener, "event:", event, "data:", data)
	}))

	d.AddListener(new(testListener), new(test2Listener))

	d.Dispatch("test", "test data")
	d.Dispatch("test2", "test2 data")
	d.Dispatch("test3", "test3 data")

	time.Sleep(time.Second)
}
```

Output:

```bash
event: test data: test data
event: test2 data: test2 data
event: test data: test data
```
