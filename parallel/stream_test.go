package parallel_test

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	"github.com/go-fries/fries/parallel/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMapStreamPreservesInputIndexesAndItemErrors(t *testing.T) {
	const count = 100
	input := make(chan int, count)
	for i := range count {
		input <- i
	}
	close(input)
	wantErr := errors.New("item failed")
	stream := newTestStream(t, 8, input, func(_ context.Context, value int) (int, error) {
		if value%7 == 0 {
			return value * 2, wantErr
		}

		return value * 2, nil
	})
	seen := make(map[int]bool, count)
	for result := range stream.Results() {
		require.GreaterOrEqual(t, result.Index, 0)
		require.Less(t, result.Index, count)
		assert.False(t, seen[result.Index], "duplicate index %d", result.Index)
		seen[result.Index] = true
		assert.Equal(t, result.Index*2, result.Value)
		if result.Index%7 == 0 {
			assert.ErrorIs(t, result.Err, wantErr)
		} else {
			assert.NoError(t, result.Err)
		}
	}
	require.NoError(t, stream.Wait(t.Context()))
	assert.Len(t, seen, count)
}

func TestMapStreamDeliversWithoutWaitingForEarlierInputs(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		input := make(chan int, 2)
		input <- 0
		input <- 1
		close(input)
		release := make(chan struct{})
		stream := newTestStream(t, 2, input, func(ctx context.Context, value int) (int, error) {
			if value == 0 {
				select {
				case <-release:
				case <-ctx.Done():
					return 0, context.Cause(ctx)
				}
			}

			return value, nil
		})
		synctest.Wait()
		first := <-stream.Results()
		assert.Equal(t, parallel.StreamResult[int]{Index: 1, Value: 1}, first)
		close(release)
		second := <-stream.Results()
		assert.Equal(t, parallel.StreamResult[int]{Index: 0, Value: 0}, second)
		_, ok := <-stream.Results()
		assert.False(t, ok)
		require.NoError(t, stream.Wait(t.Context()))
	})
}

func TestMapStreamBoundsConcurrency(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		const limit = 4
		input := make(chan int, 20)
		for i := range 20 {
			input <- i
		}
		close(input)
		release := make(chan struct{})
		var active atomic.Int32
		stream := newTestStream(t, limit, input, func(ctx context.Context, value int) (int, error) {
			n := active.Add(1)
			defer active.Add(-1)
			assert.LessOrEqual(t, n, int32(limit))
			select {
			case <-release:
				return value, nil
			case <-ctx.Done():
				return 0, context.Cause(ctx)
			}
		})
		synctest.Wait()
		assert.Equal(t, int32(limit), active.Load())
		close(release)
		count := 0
		for range stream.Results() {
			count++
		}
		require.NoError(t, stream.Wait(t.Context()))
		assert.Equal(t, 20, count)
		assert.Zero(t, active.Load())
	})
}

func TestMapStreamBackpressureBoundsInputConsumption(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		const limit = 3
		input := make(chan int, 20)
		for i := range 20 {
			input <- i
		}
		var called atomic.Int32
		stream := newTestStream(t, limit, input, func(_ context.Context, value int) (int, error) {
			called.Add(1)

			return value, nil
		})
		synctest.Wait()
		assert.Equal(t, int32(limit), called.Load())
		assert.Len(t, input, 20-limit)
		<-stream.Results()
		synctest.Wait()
		assert.Equal(t, int32(limit+1), called.Load())
		assert.Len(t, input, 20-limit-1)
		stream.Close()
		require.ErrorIs(t, stream.Wait(t.Context()), context.Canceled)
		_, ok := <-stream.Results()
		assert.False(t, ok)
		// The caller still owns the unconsumed input and its channel.
		assert.Len(t, input, 20-limit-1)
		close(input)
	})
}

func TestMapStreamEmptyInputCompletesNormally(t *testing.T) {
	input := make(chan int)
	close(input)
	stream := newTestStream(t, 4, input, func(context.Context, int) (int, error) {
		t.Error("callback invoked for an empty input")

		return 0, nil
	})
	_, ok := <-stream.Results()
	assert.False(t, ok)
	require.NoError(t, stream.Wait(t.Context()))
	assert.ErrorIs(t, context.Cause(stream.Context()), context.Canceled)
	stream.Close()
	require.NoError(t, stream.Wait(t.Context()))
}

func TestStreamCloseUnblocksIdleReceivers(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		input := make(chan int)
		stream := newTestStream(t, 8, input, func(context.Context, int) (int, error) {
			t.Error("callback invoked without an input")

			return 0, nil
		})
		synctest.Wait()
		stream.Close()
		require.ErrorIs(t, stream.Wait(t.Context()), context.Canceled)
		_, ok := <-stream.Results()
		assert.False(t, ok)
		select {
		case <-input:
			t.Fatal("stream closed the producer's input channel")
		default:
		}
		close(input)
	})
}

func TestStreamCloseWaitsForRunningCallbacks(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		input := make(chan int, 1)
		input <- 1
		release := make(chan struct{})
		canceled := make(chan struct{})
		stream := newTestStream(t, 1, input, func(ctx context.Context, value int) (int, error) {
			<-ctx.Done()
			close(canceled)
			<-release

			return value, context.Cause(ctx)
		})
		var once sync.Once
		unblock := func() { once.Do(func() { close(release) }) }
		t.Cleanup(unblock)
		synctest.Wait()
		closed := make(chan struct{})
		go func() {
			stream.Close()
			close(closed)
		}()
		<-canceled
		synctest.Wait()
		select {
		case <-closed:
			t.Fatal("Close returned before the callback exited")
		default:
		}
		unblock()
		<-closed
		require.ErrorIs(t, stream.Wait(t.Context()), context.Canceled)
		close(input)
	})
}

func TestStreamContextStopsProducerOnEarlyClose(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		input := make(chan int)
		stream := newTestStream(t, 2, input, func(_ context.Context, value int) (int, error) {
			return value, nil
		})
		producerDone := make(chan struct{})
		go func() {
			defer close(producerDone)
			defer close(input)
			for i := 0; ; i++ {
				select {
				case input <- i:
				case <-stream.Context().Done():
					return
				}
			}
		}()
		<-stream.Results()
		stream.Close()
		<-producerDone
		require.ErrorIs(t, stream.Wait(t.Context()), context.Canceled)
	})
}

func TestStreamWaitCancellationDoesNotCancelProcessing(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		input := make(chan int, 1)
		input <- 42
		close(input)
		stream := newTestStream(t, 1, input, func(_ context.Context, value int) (int, error) {
			return value, nil
		})
		ctx, cancel := context.WithTimeout(t.Context(), time.Second)
		defer cancel()
		require.ErrorIs(t, stream.Wait(ctx), context.DeadlineExceeded)
		assert.NoError(t, context.Cause(stream.Context()))
		assert.Equal(t, 42, (<-stream.Results()).Value)
		require.NoError(t, stream.Wait(t.Context()))
		require.NoError(t, stream.Wait(ctx))
		stream.Close()
		require.NoError(t, stream.Wait(ctx))
	})
}

func TestStreamPreservesParentCancellationCause(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancelCause(t.Context())
		defer cancel(nil)
		input := make(chan int, 1)
		input <- 1
		var called atomic.Bool
		stream, err := parallel.MapStream(ctx, 2, input, func(ctx context.Context, value int) (int, error) {
			called.Store(true)
			<-ctx.Done()

			return value, context.Cause(ctx)
		})
		require.NoError(t, err)
		t.Cleanup(stream.Close)
		synctest.Wait()
		require.True(t, called.Load())
		wantErr := errors.New("source stopped")
		cancel(wantErr)
		require.ErrorIs(t, stream.Wait(t.Context()), wantErr)
		require.ErrorIs(t, context.Cause(stream.Context()), wantErr)
		_, ok := <-stream.Results()
		assert.False(t, ok)
		close(input)
	})
}

func TestStreamConcurrentCloseAndWait(t *testing.T) {
	input := make(chan int)
	stream := newTestStream(t, 8, input, func(_ context.Context, value int) (int, error) { return value, nil })
	start := make(chan struct{})
	var callers sync.WaitGroup
	for range 20 {
		callers.Go(func() {
			<-start
			stream.Close()
		})
		callers.Go(func() {
			<-start
			assert.ErrorIs(t, stream.Wait(t.Context()), context.Canceled)
		})
	}
	close(start)
	callers.Wait()
	close(input)
}

func TestMapStreamValidatesInputs(t *testing.T) {
	fn := func(_ context.Context, value int) (int, error) { return value, nil }
	input := make(chan int)
	close(input)
	for _, limit := range []int{0, -1} {
		stream, err := parallel.MapStream(t.Context(), limit, input, fn)
		require.ErrorIs(t, err, parallel.ErrInvalidLimit)
		assert.Nil(t, stream)
	}
	stream, err := parallel.MapStream[int, int](t.Context(), 1, input, nil)
	require.ErrorIs(t, err, parallel.ErrNilFunc)
	assert.Nil(t, stream)
	stream, err = parallel.MapStream(t.Context(), 1, nil, fn)
	require.ErrorIs(t, err, parallel.ErrNilInput)
	assert.Nil(t, stream)
	ctx, cancel := context.WithCancelCause(t.Context())
	wantErr := errors.New("already stopped")
	cancel(wantErr)
	stream, err = parallel.MapStream(ctx, 1, input, fn)
	require.ErrorIs(t, err, wantErr)
	assert.Nil(t, stream)
}

func newTestStream[T, R any](t *testing.T, limit int, input <-chan T, fn func(context.Context, T) (R, error)) *parallel.Stream[R] {
	t.Helper()
	stream, err := parallel.MapStream(t.Context(), limit, input, fn)
	require.NoError(t, err)
	t.Cleanup(stream.Close)

	return stream
}
