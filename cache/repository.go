package cache

import (
	"context"
	"time"

	"github.com/go-packagist/go-kratos-components/contract/cache"
	"github.com/go-packagist/go-kratos-components/helper"
)

type Repository struct {
	cache.Store
}

func NewRepository(store cache.Store) cache.Repository {
	return &Repository{
		Store: store,
	}
}

func (r *Repository) Missing(ctx context.Context, key string) (bool, error) {
	if hatted, err := r.Store.Has(ctx, key); err != nil {
		return false, err
	} else {
		return hatted, nil
	}
}

func (r *Repository) Delete(ctx context.Context, key string) error {
	return r.Store.Forget(ctx, key)
}

func (r *Repository) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return r.Store.Put(ctx, key, value, ttl)
}

func (r *Repository) Add(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
	if addable, ok := r.Store.(cache.Addable); ok {
		return addable.Add(ctx, key, value, ttl)
	}

	if missed, err := r.Missing(ctx, key); err != nil {
		return false, err
	} else if !missed {
		if err := r.Set(ctx, key, value, ttl); err != nil {
			return false, err
		} else {
			return true, nil
		}
	}

	return false, nil
}

func (r *Repository) Remember(ctx context.Context, key string, dest interface{}, value func() interface{}, ttl time.Duration) error {
	if missed, err := r.Missing(ctx, key); err != nil {
		return err
	} else if missed {
		valued := value()

		if err := r.Set(ctx, key, valued, ttl); err != nil {
			return err
		}

		return helper.ValueOf(valued, dest)
	}

	return r.Get(ctx, key, dest)
}
