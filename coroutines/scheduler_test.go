package coroutines

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestScheduler(t *testing.T) {
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

	sc.Run(fns...)

	t.Logf("time: %v", time.Since(now))

	assert.True(t, time.Since(now) < time.Duration(n/s+1)*100*time.Millisecond)
	assert.True(t, time.Since(now) > time.Duration(n/s-1)*100*time.Millisecond)
	assert.Len(t, ch, 0)
}
