package redis

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func newRedis(t *testing.T) redis.UniversalClient {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	t.Cleanup(func() {
		_ = rdb.Close()
	})
	return rdb
}

func TestLocker_Get(t *testing.T) {
	rdb := newRedis(t)

	locker := NewLocker(rdb,
		WithName("locker:test:get"),
		WithTTL(time.Second*100),
	)

	owner1, ok1, err := locker.Get(context.Background())
	assert.NoError(t, err)
	assert.True(t, ok1)
	assert.NotEmpty(t, owner1.Name())

	_, ok2, err := locker.Get(context.Background())
	assert.NoError(t, err)
	assert.False(t, ok2)

	owner1.Release(context.Background()) //nolint:errcheck

	owner3, ok3, err := locker.Get(context.Background())
	assert.NoError(t, err)
	t.Logf("owner3: %v, ok3: %v", owner3, ok3)
	owner3.Release(context.Background()) //nolint:errcheck
}
