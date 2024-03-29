package coroutines

import (
	"sync"
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
		time.Sleep(200 * time.Millisecond)
		ch <- struct{}{}
	}, func() {
		time.Sleep(200 * time.Millisecond)
		ch <- struct{}{}
	})

	assert.True(t, len(ch) == 2)
	assert.True(t, time.Since(now) < 400*time.Millisecond)
}

func TestRun(t *testing.T) {
	var (
		ch  = make(chan struct{}, 2)
		now = time.Now()
		wg  sync.WaitGroup
	)

	wg.Add(2)

	Run(func() {
		defer wg.Done()
		time.Sleep(200 * time.Millisecond)
		ch <- struct{}{}
	}, func() {
		defer wg.Done()
		time.Sleep(200 * time.Millisecond)
		ch <- struct{}{}
	})

	wg.Wait()
	assert.True(t, len(ch) == 2)
	assert.True(t, time.Since(now) < 400*time.Millisecond)
}

func TestParallel(t *testing.T) {
	var (
		ch  = make(chan struct{}, 3)
		now = time.Now()
		wg  sync.WaitGroup
	)
	wg.Add(3)

	Parallel(2, func() {
		defer wg.Done()
		time.Sleep(200 * time.Millisecond)
		ch <- struct{}{}
	}, func() {
		defer wg.Done()
		time.Sleep(200 * time.Millisecond)
		ch <- struct{}{}
	}, func() {
		defer wg.Done()
		time.Sleep(200 * time.Millisecond)
		ch <- struct{}{}
	})
	wg.Wait()

	assert.True(t, len(ch) == 3)
	assert.True(t, time.Since(now) < 600*time.Millisecond)
	assert.True(t, time.Since(now) > 400*time.Millisecond)
}

func TestParallelWait(t *testing.T) {
	var (
		ch  = make(chan struct{}, 3)
		now = time.Now()
	)

	ParallelWait(2, func() {
		time.Sleep(200 * time.Millisecond)
		ch <- struct{}{}
	}, func() {
		time.Sleep(200 * time.Millisecond)
		ch <- struct{}{}
	}, func() {
		time.Sleep(200 * time.Millisecond)
		ch <- struct{}{}
	})

	assert.True(t, len(ch) == 3)
	assert.True(t, time.Since(now) < 600*time.Millisecond)
	assert.True(t, time.Since(now) > 400*time.Millisecond)
}
