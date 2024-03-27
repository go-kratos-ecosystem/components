package config

import "context"

type configKey struct{}

func NewContext[T any](ctx context.Context, config T) context.Context {
	return context.WithValue(ctx, configKey{}, config)
}

func FromContext[T any](ctx context.Context) (T, bool) {
	config, ok := ctx.Value(configKey{}).(T)
	return config, ok
}
