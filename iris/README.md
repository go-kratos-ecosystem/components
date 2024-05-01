# iris server

- https://github.com/kataras/iris

## Example

```go
package main

import (
	"context"
	"time"

	"github.com/kataras/iris/v12"
	
	iriss "github.com/go-kratos-ecosystem/components/v2/iris"
)

func main() {
	srv := iriss.NewServer(
		iris.New(),
		iriss.Addr(":8002"),
		iriss.WithConfigurators(
			iris.WithTimeout(10*time.Second),
		),
	)
	srv.Get("/ping", func(ctx iris.Context) {
		_, _ = ctx.WriteString("pong")
	})

	if err := srv.Start(context.Background()); err != nil {
		panic(err)
	}
}
```