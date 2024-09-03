package redis

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/go-kratos-ecosystem/components/v2/locker"
)

// releaseLockScript is a Lua script to release a lock in an atomic way.
// KEYS[1] is the lock key
// ARGV[1] is the lock value
const releaseLockScript = `if redis.call("get",KEYS[1]) == ARGV[1] then
    return redis.call("del",KEYS[1])
else
    return 0
end`

type Locker struct {
	redis   redis.UniversalClient
	name    string        // lock key
	timeout time.Duration // lock ttl

	owner string        // lock value
	sleep time.Duration // sleep duration
}

type Option func(*Locker)

func WithOwner(owner string) Option {
	return func(l *Locker) {
		l.owner = owner
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(l *Locker) {
		l.timeout = timeout
	}
}

// Deprecated: Use WithTimeout instead.
func WithSeconds(timeout time.Duration) Option {
	return func(l *Locker) {
		l.timeout = timeout
	}
}

func WithSleep(sleep time.Duration) Option {
	return func(l *Locker) {
		l.sleep = sleep
	}
}

var _ locker.Locker = (*Locker)(nil)

func NewLocker(redis redis.UniversalClient, name string, timeout time.Duration, opts ...Option) *Locker {
	l := &Locker{
		redis:   redis,
		name:    name,
		timeout: timeout,
		owner:   uuid.New().String(),
		sleep:   time.Millisecond * 100, //nolint:mnd
	}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

func (l *Locker) acquire(ctx context.Context) (bool, error) {
	return l.redis.SetNX(ctx, l.name, l.owner, l.timeout).Result()
}

func (l *Locker) Try(ctx context.Context, fn func() error) error {
	if ok, err := l.acquire(ctx); err != nil {
		return err
	} else if !ok {
		return locker.ErrLocked
	}

	defer l.Release(ctx) //nolint:errcheck
	return fn()
}

func (l *Locker) Until(ctx context.Context, timeout time.Duration, fn func() error) error {
	starting := time.Now()

	for {
		if ok, err := l.acquire(ctx); err != nil {
			return err
		} else if ok {
			break
		}

		if time.Since(starting) > timeout {
			return locker.ErrTimeout
		}

		time.Sleep(l.sleep)
	}

	defer l.Release(ctx) //nolint:errcheck
	return fn()
}

func (l *Locker) Release(ctx context.Context) (bool, error) {
	if val, err := l.redis.Eval(ctx, releaseLockScript, []string{l.name}, l.owner).Result(); err != nil {
		return false, err
	} else if val == int64(0) {
		return false, nil
	}

	return true, nil
}

func (l *Locker) ForceRelease(ctx context.Context) error {
	return l.redis.Del(ctx, l.name).Err()
}

func (l *Locker) Owner() string {
	return l.owner
}

func (l *Locker) LockedOwner(ctx context.Context) (string, error) {
	val, err := l.redis.Get(ctx, l.name).Result()
	if err != nil {
		return "", err
	} else if val == "" {
		return "", locker.ErrNotLocked
	}

	return val, nil
}
