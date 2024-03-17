package features

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExpireFeature(t *testing.T) {
	f := NewExpireFeature(100 * time.Millisecond)

	assert.Equal(t, 100*time.Millisecond, f.Expiration())
}
