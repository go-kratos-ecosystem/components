package cache

import (
	"context"
	"errors"
	"time"
)

var (
	ErrKeyAlreadyExists = errors.New("cache: key already exists")
)

type Repository interface {
	Store
	Addable

	Missing(ctx context.Context, key string) bool
	Delete(ctx context.Context, key string) error
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Remember(ctx context.Context, key string, dest interface{}, value func() interface{}, ttl time.Duration) error
}
