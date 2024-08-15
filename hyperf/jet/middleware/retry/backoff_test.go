package retry

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBackoff_NoBackoff(t *testing.T) {
	backoff := NoBackoff()
	for i := 1; i <= 10; i++ {
		assert.Equal(t, time.Duration(0), backoff(i))
	}
}

func TestBackoff_LinearBackoff(t *testing.T) {
	backoff := LinearBackoff(time.Second)
	for i := 1; i <= 10; i++ {
		assert.Equal(t, time.Duration(i)*time.Second, backoff(i))
	}
}

func TestBackoff_ExponentialBackoff(t *testing.T) {
	backoff := ExponentialBackoff(time.Second)
	for i := 1; i <= 10; i++ {
		assert.Equal(t, time.Duration(1<<uint(i))*time.Second, backoff(i))
	}
}

func TestBackoff_ConstantBackoff(t *testing.T) {
	backoff := ConstantBackoff(time.Second)
	for i := 1; i <= 10; i++ {
		assert.Equal(t, time.Second, backoff(i))
	}
}
