package foundation

import "context"

type Provider interface {
	Bootstrap(context.Context) (context.Context, error)
	Terminate(context.Context) (context.Context, error)
}

type BaseProvider struct{}

func (*BaseProvider) Bootstrap(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (*BaseProvider) Terminate(ctx context.Context) (context.Context, error) {
	return ctx, nil
}
