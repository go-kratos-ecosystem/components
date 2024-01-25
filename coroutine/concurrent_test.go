package coroutine

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConcurrent(t *testing.T) {
	var (
		start  = time.Now()
		buffer bytes.Buffer
	)

	Concurrent(2, func() {
		time.Sleep(1 * time.Second)
		buffer.WriteString("1")
	}, func() {
		time.Sleep(2 * time.Second)
		buffer.WriteString("2")
	}, func() {
		time.Sleep(3 * time.Second)
		buffer.WriteString("3")
	})

	assert.Equal(t, "123", buffer.String())
	assert.True(t, time.Since(start) > 4*time.Second)
	assert.True(t, time.Since(start) < 6*time.Second)
}
