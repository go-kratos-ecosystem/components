package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/go-kratos-ecosystem/components/v2/crontab/jobwrapper/mutex"
)

type Locker struct {
	redis.UniversalClient
}

func NewLocker(client redis.UniversalClient) *Locker {
	return &Locker{
		UniversalClient: client,
	}
}

func (m *Locker) Lock(slug string, expiration time.Duration) error {
	if flag, err := m.SetNX(context.Background(), slug, "1", expiration).Result(); err != nil {
		return err
	} else if !flag {
		return mutex.ErrLocked
	}
	return nil
}

func (m *Locker) Unlock(slug string) error {
	return m.Del(context.Background(), slug).Err()
}
