package parallel_test

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/go-fries/fries/parallel/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubmitValueConcurrentWaits(t *testing.T) {
	pool := requirePool(t, 1)
	release := make(chan struct{})
	wantErr := errors.New("partial result")
	future, err := parallel.SubmitValue(t.Context(), pool, func(context.Context) (string, error) {
		<-release

		return "value", wantErr
	})
	require.NoError(t, err)
	ctx, cancel := context.WithCancelCause(t.Context())
	waitErr := errors.New("wait stopped")
	cancel(waitErr)
	value, err := future.Wait(ctx)
	require.ErrorIs(t, err, waitErr)
	assert.Empty(t, value)

	var waiters sync.WaitGroup
	for range 20 {
		waiters.Go(func() {
			value, err := future.Wait(t.Context())
			assert.ErrorIs(t, err, wantErr)
			assert.Equal(t, "value", value)
		})
	}
	close(release)
	waiters.Wait()
	<-future.Done()
	value, err = future.Wait(ctx)
	require.ErrorIs(t, err, wantErr)
	assert.Equal(t, "value", value)
	require.NoError(t, pool.Shutdown(t.Context()))
}

func TestSubmitValueQueuedCancellation(t *testing.T) {
	pool := requirePool(t, 1)
	release := make(chan struct{})
	started := make(chan struct{})
	_, err := pool.Submit(t.Context(), func(context.Context) error {
		close(started)
		<-release

		return nil
	})
	require.NoError(t, err)
	<-started
	ctx, cancel := context.WithCancelCause(t.Context())
	future, err := parallel.SubmitValue(ctx, pool, func(context.Context) (int, error) {
		t.Error("canceled queued callback ran")

		return 42, nil
	})
	require.NoError(t, err)
	wantErr := errors.New("task stopped")
	cancel(wantErr)
	close(release)
	value, err := future.Wait(t.Context())
	require.ErrorIs(t, err, wantErr)
	assert.Zero(t, value)
	require.NoError(t, pool.Shutdown(t.Context()))
}

func TestSubmitValueValidationAndSuccess(t *testing.T) {
	pool := requirePool(t, 1)
	fn := func(context.Context) (int, error) { return 42, nil }
	future, err := parallel.SubmitValue[int](t.Context(), pool, nil)
	require.ErrorIs(t, err, parallel.ErrNilFunc)
	assert.Nil(t, future)
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	future, err = parallel.SubmitValue(ctx, pool, fn)
	require.ErrorIs(t, err, context.Canceled)
	assert.Nil(t, future)
	future, err = parallel.SubmitValue(t.Context(), pool, fn)
	require.NoError(t, err)
	value, err := future.Wait(t.Context())
	require.NoError(t, err)
	assert.Equal(t, 42, value)
	require.NoError(t, pool.Shutdown(t.Context()))
	future, err = parallel.SubmitValue(t.Context(), pool, fn)
	require.ErrorIs(t, err, parallel.ErrPoolClosed)
	assert.Nil(t, future)
}
