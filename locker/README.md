# Locker

## Usage

```go
package main

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/go-kratos-ecosystem/components/v2/cache"
	redisStore "github.com/go-kratos-ecosystem/components/v2/cache/redis"
	redisLocker "github.com/go-kratos-ecosystem/components/v2/locker/redis"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// ex1: 直接用
	locker := redisLocker.NewLocker(client, redisLocker.WithName("lock"), redisLocker.WithTTL(5*time.Minute))
	_ = locker.Try(context.Background(), func() {
		// do something
	})

	// ex2: 基于缓存用
	repository := cache.NewRepository(
		redisStore.New(client, redisStore.Prefix("cache")),
	)

	_ = repository.Lock("lock", 5*time.Minute).Try(context.Background(), func() {
		// do something
	})
}

```