package context

import (
	"context"
)

type Provider func(ctx context.Context) (context.Context, error)

func Chain(ctx context.Context, providers ...Provider) (context.Context, error) {
	var err error
	for _, provider := range providers {
		if provider != nil {
			if ctx, err = provider(ctx); err != nil {
				return ctx, err
			}
		}
	}
	return ctx, nil
}
