package timezone

import (
	"context"
	"time"
)

type contextKey struct{}

func NewContext(ctx context.Context, local *time.Location) context.Context {
	return context.WithValue(ctx, contextKey{}, local)
}

func FromContext(ctx context.Context) (*time.Location, bool) {
	local, ok := ctx.Value(contextKey{}).(*time.Location)
	return local, ok
}
