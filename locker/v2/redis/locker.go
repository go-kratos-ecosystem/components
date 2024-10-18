package redis

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/go-kratos-ecosystem/components/v2/locker/v2"
)

// releaseScript is a Lua script to release a lock in an atomic way.
//
//	KEYS[1] is the lock key
//	ARGV[1] is the lock value
const releaseScript = `if redis.call("get",KEYS[1]) == ARGV[1] then
    return redis.call("del",KEYS[1])
else
    return 0
end`

type Locker struct {
	redis redis.UniversalClient
	name  string
	ttl   time.Duration
}

type Option func(*Locker)

func WithName(name string) Option {
	return func(l *Locker) {
		l.name = name
	}
}

func WithTTL(ttl time.Duration) Option {
	return func(l *Locker) {
		l.ttl = ttl
	}
}

var _ locker.Locker = (*Locker)(nil)

func NewLocker(redis redis.UniversalClient, opts ...Option) *Locker {
	l := &Locker{
		redis: redis,
		name:  uuid.New().String(),
		ttl:   time.Second * 10, //nolint:mnd
	}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

func (l *Locker) Try(ctx context.Context, fn func()) (bool, error) {
	owner, ok, err := l.Get(ctx)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	defer func() {
		_, _ = owner.Release(ctx)
	}()

	fn()

	return true, nil
}

func (l *Locker) Until(ctx context.Context, timeout time.Duration, fn func()) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var owner locker.Owner

	for {
		owner = locker.NewOwner(l)
		if ok, err := l.acquire(ctx, owner); err != nil {
			return false, err
		} else if ok {
			break
		}

		select {
		case <-ctx.Done():
			return false, nil
		default:
			time.Sleep(time.Millisecond * 100) //nolint:mnd // todo: configurable
		}
	}

	defer func() {
		_, _ = owner.Release(ctx)
	}()

	fn()

	return true, nil
}

func (l *Locker) Get(ctx context.Context) (locker.Owner, bool, error) {
	owner := locker.NewOwner(l)
	ok, err := l.acquire(ctx, owner)
	if err != nil {
		return nil, false, err
	}

	return owner, ok, nil
}

func (l *Locker) Release(ctx context.Context, owner locker.Owner) (bool, error) {
	if val, err := l.redis.Eval(ctx, releaseScript, []string{l.name}, owner.Name()).Result(); err != nil {
		return false, err
	} else if val == int64(0) {
		return false, nil
	}

	return true, nil
}

func (l *Locker) ForceRelease(ctx context.Context) error {
	return l.redis.Del(ctx, l.name).Err()
}

func (l *Locker) LockedOwner(ctx context.Context) (locker.Owner, error) {
	val, err := l.redis.Get(ctx, l.name).Result()
	if err != nil {
		return nil, err
	} else if val == "" {
		return nil, locker.ErrNotLocked
	}

	return locker.NewOwner(l, locker.WithName(val)), nil
}

func (l *Locker) acquire(ctx context.Context, owner locker.Owner) (bool, error) {
	return l.redis.SetNX(ctx, l.name, owner.Name(), l.ttl).Result()
}
