# Crontab

## Usage

```go
package main

import (
	"fmt"

	"github.com/go-kratos/kratos/v2"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"

	"github.com/go-packagist/go-kratos-components/crontab"
	redisMutex "github.com/go-packagist/go-kratos-components/crontab/mutex/redis"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})

	if err := kratos.New(
		kratos.Server(
			NewCrontabServer(rdb),
		),
	).Run(); err != nil {
		panic(err)
	}
}

func NewCrontabServer(rdb redis.Cmdable) *crontab.Server {
	c := cron.New(
		cron.WithSeconds(),
	)

	c.AddFunc("*/1 * * * * *", func() {
		fmt.Println("Every hour on the half hour")
	})

	return crontab.NewServer(
		c,
		crontab.WithName("crontab:server"),
		crontab.WithDebug(),
		crontab.WithMutex(
			redisMutex.New(rdb),
		),
	)
}
```