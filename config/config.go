package config

import "context"

type configKey struct{}

func NewContext[T any](ctx context.Context, config T) (context.Context, error) {
	return context.WithValue(ctx, configKey{}, config), nil
}

func FromContext[T any](ctx context.Context) (T, bool) {
	config, ok := ctx.Value(configKey{}).(T)
	return config, ok
}

func Provider[T any](config T) func(ctx context.Context) (context.Context, error) {
	return func(ctx context.Context) (context.Context, error) {
		return NewContext(ctx, config)
	}
}
