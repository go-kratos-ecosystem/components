package coroutine

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWait(t *testing.T) {
	var msg string

	Wait(func() {
		time.Sleep(time.Second)
		msg = "hello"
	})

	assert.Equal(t, "hello", msg)
}
