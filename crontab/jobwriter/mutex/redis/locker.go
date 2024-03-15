package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Locker struct {
	redis.Cmdable
}

func NewLocker(client redis.Cmdable) *Locker {
	return &Locker{
		Cmdable: client,
	}
}

func (m *Locker) Lock(name string, expiration time.Duration) error {
	return m.SetNX(context.Background(), name, "1", expiration).Err()
}

func (m *Locker) Unlock(name string) error {
	return m.Del(context.Background(), name).Err()
}
