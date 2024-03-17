package coroutines

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWait(t *testing.T) {
	var (
		ch  = make(chan struct{}, 2)
		now = time.Now()
	)

	Wait(func() {
		time.Sleep(2 * time.Millisecond)
		ch <- struct{}{}
	}, func() {
		time.Sleep(2 * time.Millisecond)
		ch <- struct{}{}
	})

	assert.True(t, len(ch) == 2)
	assert.True(t, time.Since(now) < 4*time.Millisecond)
}

func TestParallel(t *testing.T) {
	var (
		ch  = make(chan struct{}, 3)
		now = time.Now()
	)

	Parallel(2, func() {
		time.Sleep(2 * time.Millisecond)
		ch <- struct{}{}
	}, func() {
		time.Sleep(2 * time.Millisecond)
		ch <- struct{}{}
	}, func() {
		time.Sleep(2 * time.Millisecond)
		ch <- struct{}{}
	})

	assert.True(t, len(ch) == 3)
	assert.True(t, time.Since(now) < 6*time.Millisecond)
	assert.True(t, time.Since(now) > 4*time.Millisecond)
}
