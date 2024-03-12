package env

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProvider(t *testing.T) {
	ctx, err := Provider(Prod)(context.Background())
	assert.NoError(t, err)

	env, ok := FromContext(ctx)
	assert.True(t, ok)
	assert.Equal(t, Prod, env)

	assert.True(t, Is(Prod))
	assert.False(t, Is(Dev))
}
