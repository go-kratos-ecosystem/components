package redis

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	"github.com/go-kratos-ecosystem/components/v2/locker"
)

func createRedis() redis.Cmdable {
	return redis.NewClient(&redis.Options{
		Addr: ":6379",
	})
}

func TestRedis_Lock(t *testing.T) {
	r := New(createRedis())
	var wg sync.WaitGroup
	var s int64

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := r.Lock("test", 5*time.Second).Try(context.Background(), func() error {
				time.Sleep(time.Second)
				return nil
			})
			if err != nil {
				assert.True(t, errors.Is(err, locker.ErrLocked))
			} else {
				atomic.AddInt64(&s, 1)
			}
		}()
	}
	wg.Wait()
	assert.True(t, s > 0)
}
