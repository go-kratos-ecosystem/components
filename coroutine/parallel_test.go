package coroutine

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParallel(t *testing.T) {
	start := time.Now()

	var buffer bytes.Buffer

	Parallel(func() {
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
	assert.True(t, time.Since(start) < 4*time.Second)
}
