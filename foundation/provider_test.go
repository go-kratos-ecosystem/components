package foundation

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type providerFuncKey struct{}

func TestProvider_BootstrapFunc(t *testing.T) {
	p := BootstrapFunc(func(ctx context.Context) (context.Context, error) {
		return context.WithValue(ctx, providerFuncKey{}, "test"), nil
	})

	ctx := context.Background()
	ctx, err := p.Bootstrap(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "test", ctx.Value(providerFuncKey{}))

	ctx, err = p.Terminate(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "test", ctx.Value(providerFuncKey{}))
}

func TestProvider_BootstrapFuncError(t *testing.T) {
	p := TerminateFunc(func(ctx context.Context) (context.Context, error) {
		return context.WithValue(ctx, providerFuncKey{}, "test"), nil
	})

	ctx := context.Background()
	ctx, err := p.Bootstrap(ctx)
	assert.NoError(t, err)
	assert.Nil(t, ctx.Value(providerFuncKey{}))

	ctx, err = p.Terminate(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "test", ctx.Value(providerFuncKey{}))
}
