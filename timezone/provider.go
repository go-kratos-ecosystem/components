package timezone

import (
	"context"
	"time"
)

type options struct {
	local string
}

type Option func(o *options)

func Local(name string) Option {
	return func(o *options) {
		o.local = name
	}
}

func Provider(opts ...Option) func(ctx context.Context) (context.Context, error) {
	op := options{
		local: "UTC",
	}

	for _, opt := range opts {
		opt(&op)
	}

	return func(ctx context.Context) (context.Context, error) {
		location, err := time.LoadLocation(op.local)
		if err != nil {
			return ctx, err
		}

		time.Local = location

		return NewContext(ctx, location), nil
	}
}
