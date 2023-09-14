package cache

import (
	"context"
	"testing"
	"time"

	redisCache "github.com/go-packagist/go-kratos-components/cache/redis"
	"github.com/go-packagist/go-kratos-components/contracts/cache"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

var ctx = context.Background()

func createRedisRepository() cache.Repository {
	return NewRepository(
		redisCache.New(
			redisCache.Prefix("repository"),
			redisCache.Redis(
				redis.NewClient(&redis.Options{
					Addr: "127.0.0.1:6379",
				}),
			),
		),
	)
}

func TestRepository_Add(t *testing.T) {
	r := createRedisRepository()

	added, err := r.Add(ctx, "test", 1, time.Second*10)
	assert.NoError(t, err)
	assert.True(t, added)

	added2, err2 := r.Add(ctx, "test", 1, time.Second*10)
	assert.NoError(t, err2)
	assert.False(t, added2)
}

func TestRepository_Remember(t *testing.T) {
	r := createRedisRepository()

	var value string
	err1 := r.Remember(ctx, "remember", &value, func() interface{} {
		return "test"
	}, time.Second*10)

	assert.NoError(t, err1)
	assert.Equal(t, "test", value)

	err2 := r.Remember(ctx, "remember", &value, func() interface{} {
		return "test2"
	}, time.Second*10)
	assert.NoError(t, err2)
	assert.Equal(t, "test", value)
}
