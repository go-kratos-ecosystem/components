package bootstrap

import (
	"context"
)

type Provider func(ctx context.Context) (context.Context, error)

type Context struct {
	context.Context
	providers []Provider
}

type ContextOption func(*Context)

func WithContext(ctx context.Context) ContextOption {
	return func(c *Context) {
		c.Context = ctx
	}
}

func WithProviders(providers ...Provider) ContextOption {
	return func(c *Context) {
		c.Register(providers...)
	}
}

func NewContext(opts ...ContextOption) *Context {
	ctx := &Context{}

	for _, opt := range opts {
		opt(ctx)
	}

	ctx.init()

	return ctx
}

func (c *Context) init() {
	if c.Context == nil {
		c.Context = context.Background()
	}
}

func (c *Context) Register(providers ...Provider) {
	c.providers = append(c.providers, providers...)
}

func (c *Context) Boot() (ctx context.Context, err error) {
	ctx = c.Context
	for _, provider := range c.providers {
		if ctx, err = provider(ctx); err != nil {
			return
		}
	}
	return
}
