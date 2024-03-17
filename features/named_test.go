package features

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNamedFeature(t *testing.T) {
	f := NewNamedFeature("test")

	assert.Equal(t, "test", f.Name())
}
