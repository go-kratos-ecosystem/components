package coroutine

import (
	"bytes"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConcurrent(t *testing.T) {
	var (
		start  = time.Now()
		buffer bytes.Buffer
		mu     sync.Mutex
	)

	Concurrent(2, func() {
		time.Sleep(1 * time.Second)
		mu.Lock()
		defer mu.Unlock()
		buffer.WriteString("1")
	}, func() {
		time.Sleep(2 * time.Second)
		mu.Lock()
		defer mu.Unlock()
		buffer.WriteString("2")
	}, func() {
		time.Sleep(3 * time.Second)
		mu.Lock()
		defer mu.Unlock()
		buffer.WriteString("3")
	})

	assert.Equal(t, "123", buffer.String())
	assert.True(t, time.Since(start) > 4*time.Second)
	assert.True(t, time.Since(start) < 6*time.Second)
}
