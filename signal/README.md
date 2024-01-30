# Signal Server

## Example

```go
package main

import (
	"os"
	"syscall"

	"github.com/go-kratos/kratos/v2"

	"github.com/go-kratos-ecosystem/components/v2/signal"
)

func main() {
	app := kratos.New(
		kratos.Server(newSignalServer()),
	)

	if err := app.Run(); err != nil {
		panic(err)
	}
}

func newSignalServer() *signal.Server {
	srv := signal.NewServer(
		signal.WithRecovery(signal.DefaultRecovery),
	)

	srv.Register(&exampleHandler{}, &example2Handler{})

	return srv
}

type exampleHandler struct{}

func (h *exampleHandler) Listen() []os.Signal {
	return []os.Signal{syscall.SIGUSR1, syscall.SIGUSR2}
}

func (h *exampleHandler) Handle(sig os.Signal) {
	println("exampleHandler signal:", sig)
}

type example2Handler struct{
	signal.AsyncFeature // async feature
}

func (h *example2Handler) Listen() []os.Signal {
	return []os.Signal{syscall.SIGUSR1}
}

func (h *example2Handler) Handle(os.Signal) {
	panic("example2Handler panic")
}
```

Send signal:

```bash
$ kill -SIGUSR2 42750
$ kill -SIGUSR1 42750
```

Output:

```bash
INFO msg=[Signal] server starting
exampleHandler signal: (0x104ff0240,0x1051875b8)
exampleHandler signal: (0x104ff0240,0x1051875b0)
ERROR msg=[Signal] handler panic (user defined signal 1): example2Handler panic
```