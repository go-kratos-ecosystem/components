package cache

import (
	"context"

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
