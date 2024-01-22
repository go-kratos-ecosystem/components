# Logging

## Usage

```go
package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"

	"github.com/go-kratos-ecosystem/components/v2/log/logging"
)

func main() {
	app := kratos.New(
		kratos.Logger(newLogger()),
	)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func newLogger() *logging.Logger {
	logger := logging.New(log.DefaultLogger)

	logger.Register("ts", log.With(log.DefaultLogger))

	return logger
}
```