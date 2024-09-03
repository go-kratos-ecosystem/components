# Logging - Hyperf jet middleware

logging middleware for Hyperf jet.

## Usage Example

```go
package main

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/go-kratos-ecosystem/components/v2/hyperf/jet"
	"github.com/go-kratos-ecosystem/components/v2/hyperf/jet/middleware/logging"
)

func main() {
	client, err := jet.NewClient(
		jet.WithTransporter(nil),
		jet.WithService("Example/User/MoneyService"),
	)
	if err != nil {
		panic(err)
	}

	// base usage
	client.Use(logging.New()) // use github.com/go-kratos/kratos/v2/log.DefaultLogger

	// with options
	client.Use(logging.New(
		logging.Logger(&customLogger{}),
	))

	// call service
	client.Invoke(context.Background(), "service", []any{"..."}, nil)
}

type customLogger struct{}

var _ log.Logger = (*customLogger)(nil)

func (c *customLogger) Log(level log.Level, keyvals ...any) error {
	// custom log
	return nil
}

```
