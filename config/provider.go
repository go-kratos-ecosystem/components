package config

import "context"

type Provider[T any] struct {
	config T
}

func NewProvider[T any](config T) *Provider[T] {
	return &Provider[T]{config: config}
}

func (p *Provider[T]) Bootstrap(ctx context.Context) (context.Context, error) {
	return NewContext(ctx, p.config), nil
}

func (p *Provider[T]) Terminate(ctx context.Context) (context.Context, error) {
	return ctx, nil
}
