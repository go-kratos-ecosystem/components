package redis

import (
	"context"
	"time"

	serializerContract "github.com/go-packagist/go-kratos-components/contract"
	"github.com/go-packagist/go-kratos-components/contract/cache"
	"github.com/go-packagist/go-kratos-components/serializer"
	"github.com/redis/go-redis/v9"
)

type options struct {
	prefix     string
	redis      redis.Cmdable
	serializer serializerContract.Serializable
}

func (o *options) setup() {
	if o.redis == nil {
		panic("redis is nil")
	}

	if o.serializer == nil {
		o.serializer = serializer.JsonSerializer
	}
}

type Option func(*options)

func Prefix(prefix string) Option {
	return func(o *options) {
		if prefix != "" {
			o.prefix = prefix + ":"
		}
	}
}

func Redis(redis redis.Cmdable) Option {
	return func(o *options) {
		o.redis = redis
	}
}

func Serializer(serializer serializerContract.Serializable) Option {
	return func(o *options) {
		o.serializer = serializer
	}
}

type Store struct {
	opt *options
}

var _ cache.Store = (*Store)(nil)
var _ cache.Addable = (*Store)(nil)

func New(opts ...Option) cache.Store {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}

	o.setup()

	return &Store{
		opt: o,
	}
}

func (s *Store) Has(ctx context.Context, key string) bool {
	if result := s.opt.redis.Exists(ctx, s.opt.prefix+key); result.Err() == nil && result.Val() == 1 {
		return true
	}

	return false
}

func (s *Store) Get(ctx context.Context, key string, dest interface{}) error {
	result := s.opt.redis.Get(ctx, s.opt.prefix+key)

	if result.Err() != nil {
		return result.Err()
	}

	return s.opt.serializer.Unserialize([]byte(result.Val()), dest)
}

func (s *Store) Put(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if data, err := s.opt.serializer.Serialize(value); err != nil {
		return err
	} else {
		return s.opt.redis.SetEx(ctx, s.opt.prefix+key, data, ttl).Err()
	}
}

func (s *Store) Increment(ctx context.Context, key string, value int) (int, error) {
	if result := s.opt.redis.IncrBy(ctx, s.opt.prefix+key, int64(value)); result.Err() != nil {
		return 0, result.Err()
	} else {
		return int(result.Val()), nil
	}
}

func (s *Store) Decrement(ctx context.Context, key string, value int) (int, error) {
	if result := s.opt.redis.DecrBy(ctx, s.opt.prefix+key, int64(value)); result.Err() != nil {
		return 0, result.Err()
	} else {
		return int(result.Val()), nil
	}
}

func (s *Store) Forever(ctx context.Context, key string, value interface{}) error {
	if data, err := s.opt.serializer.Serialize(value); err != nil {
		return err
	} else {
		return s.opt.redis.Set(ctx, s.opt.prefix+key, data, redis.KeepTTL).Err()
	}
}

func (s *Store) Forget(ctx context.Context, key string) error {
	return s.opt.redis.Del(ctx, s.opt.prefix+key).Err()
}

func (s *Store) Flush(ctx context.Context) error {
	return s.opt.redis.FlushAll(ctx).Err()
}

func (s *Store) GetPrefix() string {
	return s.opt.prefix
}

func (s *Store) Add(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
	if data, err := s.opt.serializer.Serialize(value); err != nil {
		return false, err
	} else {
		result := s.opt.redis.SetNX(ctx, s.opt.prefix+key, data, ttl)
		return result.Val(), result.Err()
	}
}
