package env

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProvider(t *testing.T) {
	assert.False(t, Is(Debug))

	p := NewProvider(Debug)
	ctx, err := p.Bootstrap(context.Background())
	assert.NoError(t, err)

	e1, ok := FromContext(ctx)
	assert.True(t, ok)
	assert.Equal(t, Debug, e1)

	ctx2, err := p.Terminate(ctx)
	assert.NoError(t, err)

	e2, ok := FromContext(ctx2)
	assert.True(t, ok)
	assert.Equal(t, Debug, e2)
}
