package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestManager(t *testing.T) {
	var r1, r2 Repository

	manager := NewManager(r1)
	manager.Register("r2", r2)

	assert.Equal(t, r1, manager.Driver())
	assert.Equal(t, r2, manager.Driver("r2"))
	assert.Panics(t, func() {
		manager.Driver("r3")
	})
}
