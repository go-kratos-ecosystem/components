package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnv(t *testing.T) {
	assert.Equal(t, Env("dev"), Dev)
	assert.Equal(t, Env("prod"), Prod)
	assert.Equal(t, Env("debug"), Debug)
	assert.Equal(t, Env("stage"), Stage)
}

func TestIs(t *testing.T) {
	SetEnv(Dev)
	assert.True(t, Is(Dev))
	assert.False(t, Is(Prod))
	assert.False(t, Is(Debug))
	assert.False(t, Is(Stage))
	assert.True(t, Is(Dev, Prod))
	assert.True(t, IsDev())
	assert.False(t, IsProd())
	assert.False(t, IsDebug())
	assert.False(t, IsStage())
	assert.True(t, IsUseString("dev"))
	assert.False(t, IsUseString("prod"))
	assert.True(t, IsUseString("dev", "prod"))
	assert.True(t, Is("dev"))
	assert.False(t, Is("prod"))

	SetEnv(Prod)
	assert.True(t, Is(Prod))
	assert.False(t, Is(Dev))
}

func TestCustom(t *testing.T) {
	SetEnv("online")

	assert.True(t, Is("online"))
	assert.False(t, Is("dev"))
}
