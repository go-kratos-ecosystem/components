# Crontab

## Usage

```go
package main

import (
	"log"

	"github.com/go-kratos/kratos/v2"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"

	"github.com/go-packagist/go-kratos-components/crontab"
	redisMutex "github.com/go-packagist/go-kratos-components/crontab/mutex/redis"
)

func main() {
	c := cron.New(
		cron.WithSeconds(),
	)

	c.AddFunc("* * * * * *", func() {
		log.Println("Hello world")
	})

	app := kratos.New(
		kratos.Server(
			crontab.NewServer(c,
				crontab.WithMutex(redisMutex.New(redis.NewClient(&redis.Options{
					Addr: "localhost:6379",
				}))),
				crontab.WithDebug(),
			),
		),
	)

	app.Run()
}
```

output:

```bash
2023/12/25 14:25:56 crontab: server started
2023/12/25 14:25:57 Hello world
2023/12/25 14:25:58 Hello world
2023/12/25 14:25:59 Hello world
2023/12/25 14:26:00 Hello world
```