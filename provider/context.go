package provider

import (
	"context"
)

type Provider func(ctx context.Context) (context.Context, error)

// Pipe returns a Provider that chains the provided Providers.
func Pipe(ctx context.Context, providers ...Provider) (context.Context, error) {
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

// Chain is a reverse Pipe.
func Chain(ctx context.Context, providers ...Provider) (context.Context, error) {
	var err error
	for i := len(providers) - 1; i >= 0; i-- {
		if providers[i] != nil {
			if ctx, err = providers[i](ctx); err != nil {
				return ctx, err
			}
		}
	}
	return ctx, nil
}
