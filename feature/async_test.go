package feature

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAsync(t *testing.T) {
	f := AsyncFeature{}

	assert.True(t, f.Async())
}
