package feature

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRetry(t *testing.T) {
	assert.Equal(t, 3, NewRetryFeature(3).Retries())
	assert.Equal(t, 2, NewRetryFeature(2).Retries())
}
