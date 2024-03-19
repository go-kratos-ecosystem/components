package coroutines

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestScheduler(t *testing.T) {
	ch := make(chan struct{}, 3)
	for i := 0; i < 3; i++ {
		ch <- struct{}{}
	}
	assert.Len(t, ch, 3)

	s := NewScheduler(2)
	now := time.Now()
	wg := sync.WaitGroup{}
	wg.Add(3)

	s.Run(func() {
		defer wg.Done()
		time.Sleep(100 * time.Millisecond)
		<-ch
	}, func() {
		defer wg.Done()
		time.Sleep(100 * time.Millisecond)
		<-ch
	}, func() {
		defer wg.Done()
		time.Sleep(100 * time.Millisecond)
		<-ch
	})
	wg.Wait()

	assert.True(t, time.Since(now) < 300*time.Millisecond)
	assert.True(t, time.Since(now) > 100*time.Millisecond)
	assert.Len(t, ch, 0)
}
