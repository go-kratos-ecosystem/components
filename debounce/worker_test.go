package debounce_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	"github.com/go-fries/fries/debounce/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkerTrailingEdgeAndReuse(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		calls := make(chan string, 10)
		d := newDebouncer(t, 500*time.Millisecond, func(_ context.Context, value string) error {
			calls <- value

			return nil
		})
		time.Sleep(10 * time.Second)
		synctest.Wait()
		assert.Empty(t, calls)
		require.NoError(t, d.Trigger("你"))
		time.Sleep(200 * time.Millisecond)
		require.NoError(t, d.Trigger("你好"))
		time.Sleep(200 * time.Millisecond)
		require.NoError(t, d.Trigger("你好呀"))
		time.Sleep(499 * time.Millisecond)
		synctest.Wait()
		assert.Empty(t, calls)
		time.Sleep(time.Millisecond)
		synctest.Wait()
		require.Len(t, calls, 1)
		assert.Equal(t, "你好呀", <-calls)
		time.Sleep(10 * time.Second)
		synctest.Wait()
		assert.Empty(t, calls)
		require.NoError(t, d.Trigger(""))
		time.Sleep(500 * time.Millisecond)
		synctest.Wait()
		require.Len(t, calls, 1)
		assert.Empty(t, <-calls)
	})
}

func TestWorkerContinuousTriggersPostponeExecution(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		calls := make(chan int, 10)
		d := newDebouncer(t, 100*time.Millisecond, func(_ context.Context, value int) error {
			calls <- value

			return nil
		})
		for i := range 10 {
			require.NoError(t, d.Trigger(i))
			time.Sleep(99 * time.Millisecond)
			synctest.Wait()
			assert.Empty(t, calls)
		}
		time.Sleep(time.Millisecond)
		synctest.Wait()
		require.Len(t, calls, 1)
		assert.Equal(t, 9, <-calls)
	})
}

func TestWorkerTimerRestartsAfterFlush(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		calls := make(chan int, 10)
		d := newDebouncer(t, time.Second, func(_ context.Context, value int) error {
			calls <- value

			return nil
		})
		require.NoError(t, d.Trigger(1))
		time.Sleep(600 * time.Millisecond)
		require.NoError(t, d.Flush(t.Context()))
		assert.Equal(t, 1, <-calls)
		require.NoError(t, d.Trigger(2))
		time.Sleep(400 * time.Millisecond)
		synctest.Wait()
		assert.Empty(t, calls)
		time.Sleep(600 * time.Millisecond)
		synctest.Wait()
		require.Len(t, calls, 1)
		assert.Equal(t, 2, <-calls)
	})
}

func TestWorkerCoalescesTriggersDuringExecution(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		calls := make(chan int, 10)
		release := make(chan struct{})
		var active atomic.Int32
		d := newDebouncer(t, 100*time.Millisecond, func(_ context.Context, value int) error {
			assert.Equal(t, int32(1), active.Add(1))
			defer active.Add(-1)
			calls <- value
			if value == 1 {
				<-release
			}

			return nil
		})
		unblock := releaseOnCleanup(t, release)
		require.NoError(t, d.Trigger(1))
		time.Sleep(100 * time.Millisecond)
		synctest.Wait()
		assert.Equal(t, 1, <-calls)
		require.NoError(t, d.Trigger(2))
		time.Sleep(50 * time.Millisecond)
		require.NoError(t, d.Trigger(3))
		time.Sleep(100 * time.Millisecond)
		synctest.Wait()
		assert.Empty(t, calls)
		unblock()
		synctest.Wait()
		require.Len(t, calls, 1)
		assert.Equal(t, 3, <-calls)
		assert.Zero(t, active.Load())
	})
}

func TestWorkerRetainsDeadlineAcrossRunningHandler(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		calls := make(chan int, 10)
		release := make(chan struct{})
		d := newDebouncer(t, time.Second, func(_ context.Context, value int) error {
			calls <- value
			if value == 1 {
				<-release
			}

			return nil
		})
		unblock := releaseOnCleanup(t, release)
		require.NoError(t, d.Trigger(1))
		time.Sleep(time.Second)
		synctest.Wait()
		assert.Equal(t, 1, <-calls)
		require.NoError(t, d.Trigger(2))
		time.Sleep(400 * time.Millisecond)
		unblock()
		time.Sleep(599 * time.Millisecond)
		synctest.Wait()
		assert.Empty(t, calls)
		time.Sleep(time.Millisecond)
		synctest.Wait()
		require.Len(t, calls, 1)
		assert.Equal(t, 2, <-calls)
	})
}

func TestWorkerReportsBackgroundErrorsAndContinues(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		wantErr := errors.New("save failed")
		reported := make(chan error, 10)
		calls := make(chan int, 10)
		d := newDebouncer(t, time.Second, func(_ context.Context, value int) error {
			calls <- value
			if value == 1 {
				return wantErr
			}

			return nil
		}, debounce.WithOnError(func(err error) { reported <- err }))
		require.NoError(t, d.Trigger(1))
		time.Sleep(time.Second)
		synctest.Wait()
		require.Len(t, reported, 1)
		require.ErrorIs(t, <-reported, wantErr)
		require.NoError(t, d.Flush(t.Context())) // Historical errors are not replayed.
		require.NoError(t, d.Trigger(2))
		time.Sleep(time.Second)
		synctest.Wait()
		assert.Empty(t, reported)
		require.Len(t, calls, 2)
		assert.Equal(t, 1, <-calls)
		assert.Equal(t, 2, <-calls)
	})
}

func TestWorkerCallbacksCanTriggerNextInvocation(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		calls := make(chan int, 10)
		wantErr := errors.New("retry explicitly")
		var d *debounce.Debouncer[int]
		d = newDebouncer(t, time.Second, func(_ context.Context, value int) error {
			calls <- value
			switch value {
			case 1:
				assert.NoError(t, d.Trigger(2))
			case 2:
				return wantErr
			}

			return nil
		}, debounce.WithOnError(func(err error) {
			assert.ErrorIs(t, err, wantErr)
			assert.NoError(t, d.Trigger(3))
		}))
		require.NoError(t, d.Trigger(1))
		for expected := 1; expected <= 3; expected++ {
			time.Sleep(time.Second)
			synctest.Wait()
			require.Len(t, calls, 1)
			assert.Equal(t, expected, <-calls)
		}
	})
}
