package redis

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	"github.com/go-kratos-ecosystem/components/v2/locker"
)

var ctx = context.Background()

func newRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

func TestLocker(t *testing.T) {
	client := newRedisClient()
	var (
		wg         sync.WaitGroup
		i          int64
		try, until atomic.Bool
		owner      = uuid.New().String()
	)

	redisLocker := NewLocker(
		client, "locker:test", time.Second*100,
		WithTimeout(time.Second*100),
		WithOwner(owner),
		WithSleep(time.Millisecond*100),
	)

	// Owner
	assert.Equal(t, owner, redisLocker.Owner())

	// Try
	wg.Add(10)
	for j := 0; j < 10; j++ {
		go func() {
			defer wg.Done()
			err := redisLocker.Try(ctx, func() error {
				time.Sleep(time.Millisecond * 500)
				atomic.AddInt64(&i, 1)
				return nil
			})
			if err != nil && errors.Is(err, locker.ErrLocked) {
				try.Store(true)
			}
		}()
	}
	wg.Wait()
	assert.True(t, try.Load())
	assert.Equal(t, int64(1), i)

	// Until
	wg.Add(10)
	start := time.Now()
	for j := 0; j < 10; j++ {
		go func() {
			defer wg.Done()
			err := redisLocker.Until(ctx, time.Second*5, func() error {
				time.Sleep(time.Millisecond * 500)
				atomic.AddInt64(&i, 1)
				return nil
			})
			if err != nil && errors.Is(err, locker.ErrTimeout) {
				until.Store(true)
			}
		}()
	}
	wg.Wait()
	assert.False(t, until.Load())
	assert.Equal(t, int64(11), i)
	assert.True(t, time.Since(start) > time.Second*5)
}
