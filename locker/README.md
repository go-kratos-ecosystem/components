# Locker

## Usage

```go
package main

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	redisLocker "github.com/go-kratos-ecosystem/components/v2/locker/redis"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
    })

	locker := redisLocker.NewLocker(client, "lock", 5*time.Minute)
	_ = locker.Try(context.Background(), func() error {
		// do something
		return nil
	})
}
```