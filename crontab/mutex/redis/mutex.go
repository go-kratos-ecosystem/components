package redis

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/go-kratos-ecosystem/components/v2/crontab"
)

type Mutex struct {
	redis.Cmdable
	expired  time.Duration
	serverid string
}

var _ crontab.Mutex = (*Mutex)(nil)

type Option func(*Mutex)

func WithExpired(expired time.Duration) Option {
	return func(m *Mutex) {
		m.expired = expired
	}
}

func New(redis redis.Cmdable, opts ...Option) *Mutex {
	m := &Mutex{
		Cmdable:  redis,
		expired:  time.Second * 60, //nolint:gomnd
		serverid: uuid.New().String(),
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func (m *Mutex) Lock(ctx context.Context, name string) error {
	if result := m.SetNX(ctx, name, m.serverid, m.expired); result.Err() != nil {
		return result.Err()
	} else if !result.Val() {
		if val, err := m.get(ctx, name); err != nil {
			return err
		} else if val == m.serverid {
			if err := m.refresh(ctx, name); err != nil {
				return err
			}

			return nil
		}

		return crontab.ErrAnotherServerRunning
	}

	return nil
}

func (m *Mutex) get(ctx context.Context, name string) (string, error) {
	return m.Get(ctx, name).Result()
}

func (m *Mutex) refresh(ctx context.Context, name string) error {
	return m.Expire(ctx, name, m.expired).Err()
}

func (m *Mutex) Unlock(ctx context.Context, name string) error {
	if result := m.Del(ctx, name); result.Err() != nil {
		return result.Err()
	}

	return nil
}
