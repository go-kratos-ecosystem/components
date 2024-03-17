package contexts

import (
	"context"
)

type Func func(ctx context.Context) (context.Context, error)

// Pipe returns a Provider that chains the provided Providers.
func Pipe(ctx context.Context, fns ...Func) (context.Context, error) {
	var err error
	for _, fn := range fns {
		if fn != nil {
			if ctx, err = fn(ctx); err != nil {
				return ctx, err
			}
		}
	}
	return ctx, nil
}

// Chain is a reverse Pipe.
func Chain(ctx context.Context, fns ...Func) (context.Context, error) {
	var err error
	for i := len(fns) - 1; i >= 0; i-- {
		if fns[i] != nil {
			if ctx, err = fns[i](ctx); err != nil {
				return ctx, err
			}
		}
	}
	return ctx, nil
}
