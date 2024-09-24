# Cache

## Usage Example

```go
package main

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/go-kratos-ecosystem/components/v2/cache"
	redisStore "github.com/go-kratos-ecosystem/components/v2/cache/redis"
)

var ctx = context.Background()

type User struct {
	Name string
	Age  int
}

func main() {
	// 创建个 Redis 连接客户端
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer rdb.Close()

	// create a redis store
	store := redisStore.New(rdb, redisStore.Prefix("example:cache"))

	// create a cache repository
	repository := cache.NewRepository(store)

	// set cache
	ok, err := repository.Set(ctx, "key", User{
		Name: "example",
		Age:  18,
	}, time.Second*10)
	if err != nil {
		log.Fatal(err)
	}
	_ = ok

	// get cache
	var user User
	err = repository.Get(ctx, "key", &user)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("user: %+v", user)
}
```