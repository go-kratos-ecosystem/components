# Crontab

## Example

```gopackage main

import (
	"log"
	"time"

	"github.com/go-kratos/kratos/v2"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"

	"github.com/go-kratos-ecosystem/components/v2/crontab"
	"github.com/go-kratos-ecosystem/components/v2/crontab/jobwrapper/mutex"
	redisLocker "github.com/go-kratos-ecosystem/components/v2/crontab/jobwrapper/mutex/redis"
)

type mutexJob struct {
	mutex.MutexJob // 可以偷懒实现 IsMutexJob 方法
}

// Name 分布式锁名称，同名的任务只有一个会执行
func (m *mutexJob) Name() string {
	return "mutexJob"
}

// Run 任务执行的逻辑
func (m *mutexJob) Run() {
	log.Println("mutexJob running")
}

// Expiration 自定义锁的过期时间，可选
func (m *mutexJob) Expiration() time.Duration {
	return time.Second * 5
}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	c := cron.New(
		cron.WithSeconds(),
		cron.WithChain(
			mutex.SkipIfStillMutexRunning(
				mutex.WithPrefix("crontab:"),                 // 分布式锁的前缀，默认值：`crontab:`，用于区分同一锁实例（如：redis）下的不同服务
				mutex.WithLogger(cron.DefaultLogger),         // 日志，默认是 cron.DefaultLogger
				mutex.WithLocker(redisLocker.NewLocker(rdb)), // 锁的实现，需要实现 mutex.Locker 接口，此处是 redis 实现
				mutex.WithExpiration(2*time.Hour),            // 设置锁的默认过期时间，默认值：1 hour，时间必须大于任务执行的时间；如果针对单个 Job 自定义，请实现 mutex.ExpirableJob 接口(时间必须大于任务执行的时间）
			),
		),
	)

	// 默认 Job，不受分布式锁的影响
	_, _ = c.AddFunc("* * * * * *", func() {
		log.Println("Hello world")
	})

	// MutexJob，受分布式锁的影响，job只需要实现 mutex.Job 接口即可
	_, _ = c.AddJob("*/2 * * * * *", &mutexJob{})

	app := kratos.New(
		kratos.Server(
			crontab.NewServer(c),
		),
	)

	err := app.Run()
	if err != nil {
		panic(err)
	}
}
```