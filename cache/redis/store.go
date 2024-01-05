package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/go-kratos-ecosystem/components/v2/cache"
	"github.com/go-kratos-ecosystem/components/v2/serializer"
	"github.com/go-kratos-ecosystem/components/v2/serializer/json"
)

type Store struct {
	redis redis.Cmdable

	opts *options
}

type options struct {
	prefix     string
	serializer serializer.Serializable
}

type Option func(*options)

func Prefix(prefix string) Option {
	return func(o *options) {
		if prefix != "" {
			o.prefix = prefix + ":"
		}
	}
}

func Serializer(serializer serializer.Serializable) Option {
	return func(o *options) {
		o.serializer = serializer
	}
}

var (
	_ cache.Store = (*Store)(nil)
)

func New(redis redis.Cmdable, opts ...Option) *Store {
	opt := &options{
		serializer: json.Serializer,
	}

	for _, o := range opts {
		o(opt)
	}

	return &Store{
		redis: redis,
		opts:  opt,
	}
}

func (s *Store) Has(ctx context.Context, key string) (bool, error) {
	if r := s.redis.Exists(ctx, s.opts.prefix+key); r.Err() != nil {
		return false, r.Err()
	} else {
		return r.Val() > 0, nil
	}
}

func (s *Store) Get(ctx context.Context, key string, dest interface{}) error {
	if r := s.redis.Get(ctx, s.opts.prefix+key); r.Err() != nil {
		return r.Err()
	} else {
		return s.opts.serializer.Unserialize([]byte(r.Val()), dest)
	}
}

func (s *Store) Put(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
	if valued, err := s.opts.serializer.Serialize(value); err != nil {
		return false, err
	} else if r := s.redis.Set(ctx, s.opts.prefix+key, valued, ttl); r.Err() != nil {
		return false, r.Err()
	} else {
		return r.Val() == "OK", nil
	}
}

func (s *Store) Increment(ctx context.Context, key string, value int) (int, error) {
	if r := s.redis.IncrBy(ctx, s.opts.prefix+key, int64(value)); r.Err() != nil {
		return 0, r.Err()
	} else {
		return int(r.Val()), nil
	}
}

func (s *Store) Decrement(ctx context.Context, key string, value int) (int, error) {
	if r := s.redis.DecrBy(ctx, s.opts.prefix+key, int64(value)); r.Err() != nil {
		return 0, r.Err()
	} else {
		return int(r.Val()), nil
	}
}

func (s *Store) Forever(ctx context.Context, key string, value interface{}) (bool, error) {
	if valued, err := s.opts.serializer.Serialize(value); err != nil {
		return false, err
	} else if r := s.redis.Set(ctx, s.opts.prefix+key, valued, redis.KeepTTL); r.Err() != nil {
		return false, r.Err()
	} else {
		return r.Val() == "OK", nil
	}
}

func (s *Store) Forget(ctx context.Context, key string) (bool, error) {
	if r := s.redis.Del(ctx, s.opts.prefix+key); r.Err() != nil {
		return false, r.Err()
	} else {
		return r.Val() > 0, nil
	}
}

func (s *Store) Flush(ctx context.Context) (bool, error) {
	if r := s.redis.FlushAll(ctx); r.Err() != nil {
		return false, r.Err()
	} else {
		return r.Val() == "OK", nil
	}
}

func (s *Store) GetPrefix() string {
	return s.opts.prefix
}
