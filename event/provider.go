package event

import "context"

func Provider(d *Dispatcher) func(ctx context.Context) (context.Context, error) {
	return func(ctx context.Context) (context.Context, error) {
		return NewContext(ctx, d), nil
	}
}
