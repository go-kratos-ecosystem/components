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

func TestLocker(t *testing.T) {
	rdb := newRedis(t)

	locker := NewLocker(rdb,
		WithName("locker:test"),
		WithTTL(time.Second*100),
	)

	owner, ok, err := locker.Get(context.Background())
	assert.NoError(t, err)
	assert.True(t, ok)
	defer owner.Release(context.Background())
}
