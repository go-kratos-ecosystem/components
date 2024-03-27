package config

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type config struct {
	Host string
	Port int
}

type config2 struct {
	Host string
	Port int
}

func TestConfig(t *testing.T) {
	ctx := NewContext(context.Background(), &config{
		Host: "localhost",
		Port: 8080,
	})

	cfg, ok := FromContext[*config](ctx)
	assert.True(t, ok)
	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, 8080, cfg.Port)

	cfg2, ok := FromContext[*config2](ctx)
	assert.False(t, ok)
	assert.Nil(t, cfg2)
}
