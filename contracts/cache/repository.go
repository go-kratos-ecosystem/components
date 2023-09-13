package cache

import "context"

type Repository interface {
	Store

	Missing(ctx context.Context, key string) bool
	Delete(ctx context.Context, key string) error
}
