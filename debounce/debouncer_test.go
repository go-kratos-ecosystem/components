package debounce_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"testing/synctest"
	"time"

	"github.com/go-fries/fries/debounce/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlushExecutesLatestValueOnce(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		calls := make(chan int, 10)
		d := newDebouncer(t, time.Hour, func(_ context.Context, value int) error {
			calls <- value

			return nil
		})
		require.NoError(t, d.Flush(t.Context()))
		require.NoError(t, d.Trigger(1))
		require.NoError(t, d.Trigger(2))
		require.NoError(t, d.Flush(t.Context()))
		require.Len(t, calls, 1)
		assert.Equal(t, 2, <-calls)
		time.Sleep(2 * time.Hour)
		synctest.Wait()
		require.NoError(t, d.Flush(t.Context()))
		assert.Empty(t, calls)
	})
}

func TestFlushFreezesValueBeforeLaterTriggers(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		release := make(chan struct{})
		calls := make(chan int, 10)
		wantErr := errors.New("frozen save failed")
		d := newDebouncer(t, time.Hour, func(_ context.Context, value int) error {
			calls <- value
			if value == 1 {
				<-release
			}
			if value == 2 {
				return wantErr
			}

			return nil
		})
		unblock := releaseOnCleanup(t, release)
		require.NoError(t, d.Trigger(1))
		first := flushAsync(t, t.Context(), d)
		synctest.Wait()
		require.NoError(t, d.Trigger(2))
		frozen := flushAsync(t, t.Context(), d)
		synctest.Wait()
		require.NoError(t, d.Trigger(3))
		unblock()
		require.NoError(t, <-first)
		require.ErrorIs(t, <-frozen, wantErr)
		synctest.Wait()
		require.Len(t, calls, 2)
		assert.Equal(t, 1, <-calls)
		assert.Equal(t, 2, <-calls)
		require.NoError(t, d.Flush(t.Context()))
		assert.Equal(t, 3, <-calls)
	})
}

func TestConcurrentFlushesShareRunningResult(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		release := make(chan struct{})
		wantErr := errors.New("save failed")
		calls := make(chan int, 10)
		d := newDebouncer(t, time.Hour, func(_ context.Context, value int) error {
			calls <- value
			<-release

			return wantErr
		})
		unblock := releaseOnCleanup(t, release)
		require.NoError(t, d.Trigger(1))
		first := flushAsync(t, t.Context(), d)
		synctest.Wait()
		second := flushAsync(t, t.Context(), d)
		synctest.Wait()
		assert.Empty(t, first)
		assert.Empty(t, second)
		unblock()
		require.ErrorIs(t, <-first, wantErr)
		require.ErrorIs(t, <-second, wantErr)
		assert.Len(t, calls, 1)
	})
}

func TestFlushWaitsForAdmissionBeforeCapturingLatest(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		firstRelease := make(chan struct{})
		secondRelease := make(chan struct{})
		calls := make(chan int, 10)
		d := newDebouncer(t, time.Hour, func(_ context.Context, value int) error {
			calls <- value
			switch value {
			case 1:
				<-firstRelease
			case 2:
				<-secondRelease
			}

			return nil
		})
		unblockFirst := releaseOnCleanup(t, firstRelease)
		unblockSecond := releaseOnCleanup(t, secondRelease)
		require.NoError(t, d.Trigger(1))
		first := flushAsync(t, t.Context(), d)
		synctest.Wait()
		require.NoError(t, d.Trigger(2))
		second := flushAsync(t, t.Context(), d)
		synctest.Wait()
		shared := flushAsync(t, t.Context(), d)
		synctest.Wait()
		require.NoError(t, d.Trigger(3))
		third := flushAsync(t, t.Context(), d)
		ctx, cancel := context.WithCancelCause(t.Context())
		defer cancel(nil)
		canceled := flushAsync(t, ctx, d)
		synctest.Wait()
		waitErr := errors.New("admission wait stopped")
		cancel(waitErr)
		require.ErrorIs(t, <-canceled, waitErr)
		require.NoError(t, d.Trigger(4))
		unblockFirst()
		require.NoError(t, <-first)
		synctest.Wait()
		assert.Empty(t, third)
		unblockSecond()
		require.NoError(t, <-second)
		require.NoError(t, <-shared)
		require.NoError(t, <-third)
		require.Len(t, calls, 3)
		assert.Equal(t, 1, <-calls)
		assert.Equal(t, 2, <-calls)
		assert.Equal(t, 4, <-calls)
	})
}

func TestCanceledFlushPreservesFrozenInvocation(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		release := make(chan struct{})
		calls := make(chan int, 10)
		d := newDebouncer(t, time.Hour, func(ctx context.Context, value int) error {
			if value == 1 {
				<-release
			}
			assert.NoError(t, context.Cause(ctx))
			calls <- value

			return nil
		})
		unblock := releaseOnCleanup(t, release)
		require.NoError(t, d.Trigger(1))
		first := flushAsync(t, t.Context(), d)
		synctest.Wait()
		require.NoError(t, d.Trigger(2))
		ctx, cancel := context.WithCancelCause(t.Context())
		defer cancel(nil)
		second := flushAsync(t, ctx, d)
		synctest.Wait()
		wantErr := errors.New("flush wait stopped")
		cancel(wantErr)
		require.ErrorIs(t, <-second, wantErr)
		unblock()
		require.NoError(t, <-first)
		synctest.Wait()
		require.Len(t, calls, 2)
		assert.Equal(t, 1, <-calls)
		assert.Equal(t, 2, <-calls)
	})
}

func TestFlushAndCloseWaitForErrorCallback(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		release := make(chan struct{})
		wantErr := errors.New("save failed")
		d := newDebouncer(t, time.Hour, func(context.Context, int) error { return wantErr },
			debounce.WithOnError(func(err error) {
				assert.ErrorIs(t, err, wantErr)
				<-release
			}),
		)
		unblock := releaseOnCleanup(t, release)
		require.NoError(t, d.Trigger(1))
		flushed := flushAsync(t, t.Context(), d)
		synctest.Wait()
		closed := make(chan struct{})
		go func() {
			d.Close()
			close(closed)
		}()
		synctest.Wait()
		assert.Empty(t, flushed)
		select {
		case <-closed:
			t.Fatal("Close returned before the error callback finished")
		default:
		}
		unblock()
		require.ErrorIs(t, <-flushed, wantErr)
		<-closed
	})
}

func TestCloseDiscardsPendingAndFrozenValues(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		release := make(chan struct{})
		calls := make(chan int, 10)
		d := newDebouncer(t, time.Hour, func(ctx context.Context, value int) error {
			calls <- value
			<-ctx.Done()
			<-release

			return context.Cause(ctx)
		})
		unblock := releaseOnCleanup(t, release)
		require.NoError(t, d.Trigger(1))
		first := flushAsync(t, t.Context(), d)
		synctest.Wait()
		require.NoError(t, d.Trigger(2))
		second := flushAsync(t, t.Context(), d)
		synctest.Wait()
		require.NoError(t, d.Trigger(3))
		closed := make(chan struct{})
		go func() {
			d.Close()
			close(closed)
		}()
		synctest.Wait()
		require.ErrorIs(t, <-second, debounce.ErrClosed)
		require.ErrorIs(t, d.Trigger(4), debounce.ErrClosed)
		require.ErrorIs(t, d.Flush(t.Context()), debounce.ErrClosed)
		select {
		case <-closed:
			t.Fatal("Close returned before the handler exited")
		default:
		}
		unblock()
		<-closed
		require.ErrorIs(t, <-first, context.Canceled)
		require.Len(t, calls, 1)
		assert.Equal(t, 1, <-calls)
		d.Close()
	})
}

func TestParentCancellationStopsProcessing(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancelCause(t.Context())
		defer cancel(nil)
		calls := make(chan int, 10)
		d, err := debounce.New(ctx, time.Hour, func(ctx context.Context, value int) error {
			calls <- value
			<-ctx.Done()

			return context.Cause(ctx)
		})
		require.NoError(t, err)
		t.Cleanup(d.Close)
		require.NoError(t, d.Trigger(1))
		first := flushAsync(t, t.Context(), d)
		synctest.Wait()
		require.NoError(t, d.Trigger(2))
		wantErr := errors.New("application stopped")
		cancel(wantErr)
		require.ErrorIs(t, <-first, wantErr)
		d.Close()
		require.ErrorIs(t, d.Trigger(3), debounce.ErrClosed)
		assert.Len(t, calls, 1)
	})
}

func TestFlushValidatesContextWithoutConsumingPendingValue(t *testing.T) {
	calls := make(chan int, 1)
	d := newDebouncer(t, time.Hour, func(_ context.Context, value int) error {
		calls <- value

		return nil
	})
	require.NoError(t, d.Trigger(1))
	var nilContext context.Context
	require.ErrorIs(t, d.Flush(nilContext), debounce.ErrInvalidContext)
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	require.ErrorIs(t, d.Flush(ctx), context.Canceled)
	require.NoError(t, d.Flush(t.Context()))
	assert.Equal(t, 1, <-calls)
}

func TestConcurrentTriggerFlushAndClose(t *testing.T) {
	var values []int
	d := newDebouncer(t, time.Hour, func(_ context.Context, value int) error {
		values = append(values, value)

		return nil
	})
	require.NoError(t, d.Trigger(-1))
	require.NoError(t, d.Flush(t.Context()))
	start := make(chan struct{})
	var callers sync.WaitGroup
	for i := range 100 {
		callers.Go(func() {
			<-start
			err := d.Trigger(i)
			assert.True(t, err == nil || errors.Is(err, debounce.ErrClosed))
			err = d.Flush(t.Context())
			assert.True(t, err == nil || errors.Is(err, debounce.ErrClosed))
		})
	}
	for range 10 {
		callers.Go(func() {
			<-start
			d.Close()
		})
	}
	close(start)
	callers.Wait()
	d.Close()
	seen := make(map[int]bool, len(values))
	for _, value := range values {
		assert.False(t, seen[value], "duplicate invocation for %d", value)
		seen[value] = true
	}
	assert.True(t, seen[-1])
}

func newDebouncer[T any](t *testing.T, delay time.Duration, handler debounce.Handler[T], options ...debounce.Option) *debounce.Debouncer[T] {
	t.Helper()
	d, err := debounce.New(t.Context(), delay, handler, options...)
	require.NoError(t, err)
	t.Cleanup(d.Close)

	return d
}

func releaseOnCleanup(t *testing.T, release chan struct{}) func() {
	t.Helper()
	var once sync.Once
	unblock := func() { once.Do(func() { close(release) }) }
	t.Cleanup(unblock)

	return unblock
}

func flushAsync[T any](t *testing.T, ctx context.Context, d *debounce.Debouncer[T]) <-chan error {
	t.Helper()
	result := make(chan error, 1)
	go func() { result <- d.Flush(ctx) }()

	return result
}
