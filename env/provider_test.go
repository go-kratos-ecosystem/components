package env

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProvider(t *testing.T) {
	assert.True(t, Is(Prod))
	assert.False(t, Is(Dev))

	assert.NoError(t, Provider(Dev)(context.Background()))

	assert.False(t, Is(Prod))
	assert.True(t, Is(Dev))
}
