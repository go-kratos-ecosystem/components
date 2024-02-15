package coordinator

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestManager(t *testing.T) {
	wg := sync.WaitGroup{}
	ch := make(chan struct{}, 2)

	c1 := Until("foo")
	c2 := Until("foo")

	assert.Same(t, c1, c2)
	assert.Equal(t, 0, len(ch))

	wg.Add(2)
	go func() {
		defer wg.Done()

		if <-c1.Done(); true {
			ch <- struct{}{}
			return
		}
	}()
	go func() {
		defer wg.Done()

		if <-c2.Done(); true {
			ch <- struct{}{}
			return
		}
	}()

	c1.Close()

	wg.Wait()
	assert.Equal(t, 2, len(ch))
}
