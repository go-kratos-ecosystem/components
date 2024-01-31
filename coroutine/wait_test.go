package coroutine

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWait(t *testing.T) {
	var buffer bytes.Buffer

	Wait(func() {
		time.Sleep(100 * time.Millisecond)
		buffer.WriteString("hello")
	}, func() {
		time.Sleep(200 * time.Millisecond)
		buffer.WriteString(" world")
	})

	assert.Equal(t, "hello world", buffer.String())
}
