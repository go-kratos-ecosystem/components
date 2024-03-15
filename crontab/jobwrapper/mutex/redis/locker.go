package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/go-kratos-ecosystem/components/v2/crontab/jobwrapper/mutex"
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
	if flag, err := m.SetNX(context.Background(), name, "1", expiration).Result(); err != nil {
		return err
	} else if !flag {
		return mutex.ErrLocked
	}
	return nil
}

func (m *Locker) Unlock(name string) error {
	return m.Del(context.Background(), name).Err()
}
