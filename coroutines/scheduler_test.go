package coroutines

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestScheduler_Parallel(t *testing.T) {
	var fns []func()
	var wg sync.WaitGroup

	n, s := 300, 50
	ch := make(chan struct{}, n)
	wg.Add(n)
	for i := 0; i < n; i++ {
		ch <- struct{}{}
		fns = append(fns, func() {
			defer wg.Done()
			time.Sleep(100 * time.Millisecond)
			<-ch
		})
	}

	assert.Len(t, ch, n)
	assert.Len(t, fns, n)

	sc := NewScheduler(s)
	now := time.Now()

	sc.Parallel(fns...)
	wg.Wait()

	t.Logf("time: %v", time.Since(now))

	assert.True(t, time.Since(now) < time.Duration(n/s+1)*100*time.Millisecond)
	assert.True(t, time.Since(now) > time.Duration(n/s-1)*100*time.Millisecond)
	assert.Len(t, ch, 0)
}

func TestScheduler_Wait(t *testing.T) {
	var fns []func()

	n, s := 300, 50
	ch := make(chan struct{}, n)
	for i := 0; i < n; i++ {
		ch <- struct{}{}
		fns = append(fns, func() {
			time.Sleep(100 * time.Millisecond)
			<-ch
		})
	}

	assert.Len(t, ch, n)
	assert.Len(t, fns, n)

	sc := NewScheduler(s)
	now := time.Now()

	sc.Wait(fns...)

	t.Logf("time: %v", time.Since(now))

	assert.True(t, time.Since(now) < time.Duration(n/s+1)*100*time.Millisecond)
	assert.True(t, time.Since(now) > time.Duration(n/s-1)*100*time.Millisecond)
	assert.Len(t, ch, 0)
}
