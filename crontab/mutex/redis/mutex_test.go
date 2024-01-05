package redis

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	"github.com/go-packagist/go-kratos-components/crontab"
)

var (
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	ctx = context.Background()
)

func TestMutex(t *testing.T) {
	rdb.Ping(context.Background())
	m1 := New(rdb, WithExpired(time.Second*1))
	m2 := New(rdb, WithExpired(time.Second*1))

	assert.NoError(t, m1.Lock(ctx, "test"))
	assert.NoError(t, m1.Lock(ctx, "test"))
	assert.EqualError(t, m2.Lock(ctx, "test"), crontab.ErrAnotherServerRunning.Error())

	time.Sleep(time.Second * 2)
	assert.NoError(t, m2.Lock(ctx, "test"))
	assert.EqualError(t, m1.Lock(ctx, "test"), crontab.ErrAnotherServerRunning.Error())

	assert.NoError(t, m2.Unlock(ctx, "test"))
	assert.NoError(t, m1.Lock(ctx, "test"))
	assert.EqualError(t, m2.Lock(ctx, "test"), crontab.ErrAnotherServerRunning.Error())
}
