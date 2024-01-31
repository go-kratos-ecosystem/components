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
		mu.Lock()
		defer mu.Unlock()
		time.Sleep(100 * time.Millisecond)
		buffer.WriteString("hello")
	}, func() {
		mu.Lock()
		defer mu.Unlock()
		time.Sleep(200 * time.Millisecond)
		buffer.WriteString(" world")
	})

	assert.Equal(t, "hello world", buffer.String())
}
