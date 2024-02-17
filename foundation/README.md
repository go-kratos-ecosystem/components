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
		foundation.BootstrapName,
		foundation.ShutdownName,
	}
}

func (l *listener) Handle(event event.Event, data interface{}) {
	if event.String() == foundation.BootstrapName {
		if e, ok := data.(*foundation.BootstrapEvent); ok {
			fmt.Println("bootstrap done, and the app start time: ", e.Time)
		}
	}

	if event.String() == foundation.ShutdownName {
		if e, ok := data.(*foundation.ShutdownEvent); ok {
			fmt.Println("shutdown done, and the app end time: ", e.Time)
		}
	}
}

func main() {
	m := coordinator.NewManager()
	d := event.NewDispatcher()

	go func() {
		if <-m.Until(foundation.BootstrapName).Done(); true {
			fmt.Println("bootstrap done")
		}
	}()

	go func() {
		if <-m.Until(foundation.ShutdownName).Done(); true {
			fmt.Println("shutdown done")
		}
	}()

	d.AddListener(&listener{})

	app := kratos.New(
		kratos.Server(newCrontabServer()),
		kratos.BeforeStart(foundation.NewBootstrap(
			foundation.WithManager(m),
			foundation.WithDispatcher(d),
		)),
		kratos.AfterStop(foundation.NewShutdown(
			foundation.WithManager(m),
			foundation.WithDispatcher(d),
		)),
	)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
	
	time.Sleep(time.Second)
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
bootstrap done, and the app start time:  2024-02-17 17:42:45.626551 +0800 CST m=+0.002061459
bootstrap done
hello
hello
^Cshutdown done, and the app end time:  2024-02-17 17:42:47.121271 +0800 CST m=+1.496789626
shutdown done
```