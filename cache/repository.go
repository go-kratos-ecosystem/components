package cache

import (
	"context"
	"time"

	"github.com/go-packagist/go-kratos-components/contracts/cache"
)

type Repository struct {
	cache.Store
}

func New(store cache.Store) cache.Repository {
	return &Repository{
		Store: store,
	}
}

func (r *Repository) Missing(ctx context.Context, key string) bool {
	return !r.Store.Has(ctx, key)
}

func (r *Repository) Delete(ctx context.Context, key string) error {
	return r.Store.Forget(ctx, key)
}

func (r *Repository) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return r.Store.Put(ctx, key, value, ttl)
}

func (r *Repository) Add(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if addable, ok := r.Store.(cache.Addable); ok {
		return addable.Add(ctx, key, value, ttl)
	}

	if r.Missing(ctx, key) {
		return r.Set(ctx, key, value, ttl)
	}

	return cache.ErrKeyAlreadyExists
}
