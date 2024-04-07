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

type BootstrapFunc func(context.Context) (context.Context, error)

func (f BootstrapFunc) Bootstrap(ctx context.Context) (context.Context, error) {
	return f(ctx)
}

func (f BootstrapFunc) Terminate(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

type TerminateFunc func(context.Context) (context.Context, error)

func (f TerminateFunc) Bootstrap(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (f TerminateFunc) Terminate(ctx context.Context) (context.Context, error) {
	return f(ctx)
}
