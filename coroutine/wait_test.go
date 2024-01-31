package coroutine

import (
	"bytes"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWait(t *testing.T) {
	var (
		buffer bytes.Buffer
		mu     sync.Mutex
	)

	Wait(func() {
		time.Sleep(100 * time.Millisecond)
		mu.Lock()
		defer mu.Unlock()
		buffer.WriteString("hello")
	}, func() {
		time.Sleep(200 * time.Millisecond)
		mu.Lock()
		defer mu.Unlock()
		buffer.WriteString(" world")
	})

	assert.Equal(t, "hello world", buffer.String())
}
