package cache

import (
	"context"
	"time"
)

type Store interface {
	Has(ctx context.Context, key string) (bool, error)

	Get(ctx context.Context, key string, dest interface{}) error

	Put(ctx context.Context, key string, value interface{}, ttl time.Duration) error

	Increment(ctx context.Context, key string, value int) (int, error)

	Decrement(ctx context.Context, key string, value int) (int, error)

	Forever(ctx context.Context, key string, value interface{}) error

	Forget(ctx context.Context, key string) error

	Flush(ctx context.Context) error

	GetPrefix() string
}

type Addable interface {
	Add(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error)
}
