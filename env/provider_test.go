package env

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProvider(t *testing.T) {
	assert.NoError(t, Provider(Dev)(context.Background()))

	assert.False(t, Is(Prod))
	assert.True(t, Is(Dev))
}
