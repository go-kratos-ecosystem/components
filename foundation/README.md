# Foundation

## Example

```go
package main

import (
	"fmt"
	"log"

	"github.com/go-kratos/kratos/v2"
	"github.com/robfig/cron/v3"

	"github.com/go-kratos-ecosystem/components/v2/coordinator"
	v2 "github.com/go-kratos-ecosystem/components/v2/crontab/v2"
	"github.com/go-kratos-ecosystem/components/v2/event"
	"github.com/go-kratos-ecosystem/components/v2/foundation"
)

type listener struct{}

var _ event.Listener = (*listener)(nil)

func (l *listener) Listen() []event.Event {
	return []event.Event{
		foundation.BootName,
		foundation.ShutName,
	}
}

func (l *listener) Handle(event event.Event, data interface{}) {
	if event.String() == foundation.BootName {
		if e, ok := data.(*foundation.BootEvent); ok {
			fmt.Println("bootstrap done, and the app start time: ", e.Time)
		}
	}

	if event.String() == foundation.ShutName {
		if e, ok := data.(*foundation.ShutEvent); ok {
			fmt.Println("shutdown done, and the app end time: ", e.Time)
		}
	}
}

func main() {
	m := coordinator.NewManager()
	d := event.NewDispatcher()

	go func() {
		if <-m.Until(foundation.BootName).Done(); true {
			fmt.Println("bootstrap done")
		}
	}()

	go func() {
		if <-m.Until(foundation.ShutName).Done(); true {
			fmt.Println("shutdown done")
		}
	}()

	d.AddListener(&listener{})

	app := kratos.New(
		kratos.Server(newCrontabServer()),
		kratos.BeforeStart(foundation.NewBoot(
			foundation.WithManager(m),
			foundation.WithDispatcher(d),
		)),
		kratos.AfterStop(foundation.NewShut(
			foundation.WithManager(m),
			foundation.WithDispatcher(d),
		)),
	)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func newCrontabServer() *v2.Server {
	srv := v2.NewServer(
		cron.New(cron.WithSeconds()),
	)

	srv.AddFunc("* * * * * *", func() { //nolint:errcheck
		println("hello")
	})

	return srv
}
```

output:

```bash
bootstrap done
bootstrap done, and the app start time:  2024-02-17 17:21:52.309688 +0800 CST m=+0.003057710
hello
hello
hello
^Cshutdown done, and the app end time:  2024-02-17 17:21:55.951264 +0800 CST m=+3.644671418
```