package debounce_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-fries/fries/debounce/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewValidatesInputs(t *testing.T) {
	handler := func(context.Context, int) error { return nil }
	var nilContext context.Context
	d, err := debounce.New(nilContext, time.Second, handler)
	require.ErrorIs(t, err, debounce.ErrInvalidContext)
	assert.Nil(t, d)
	ctx, cancel := context.WithCancelCause(t.Context())
	wantErr := errors.New("already stopped")
	cancel(wantErr)
	d, err = debounce.New(ctx, time.Second, handler)
	require.ErrorIs(t, err, wantErr)
	assert.Nil(t, d)
	for _, delay := range []time.Duration{0, -time.Second} {
		d, err = debounce.New(t.Context(), delay, handler)
		require.ErrorIs(t, err, debounce.ErrInvalidDelay)
		assert.Nil(t, d)
	}
	d, err = debounce.New[int](t.Context(), time.Second, nil)
	require.ErrorIs(t, err, debounce.ErrNilHandler)
	assert.Nil(t, d)
}

func TestOptionsUseLastErrorCallback(t *testing.T) {
	wantErr := errors.New("save failed")
	var received error
	d := newDebouncer(t, time.Hour, func(context.Context, int) error { return wantErr },
		debounce.WithOnError(func(error) { t.Error("replaced error callback ran") }),
		nil,
		debounce.WithOnError(func(err error) { received = err }),
	)
	require.NoError(t, d.Trigger(1))
	require.ErrorIs(t, d.Flush(t.Context()), wantErr)
	assert.ErrorIs(t, received, wantErr)
}

func TestNilErrorCallbackDisablesNotifications(t *testing.T) {
	wantErr := errors.New("save failed")
	d := newDebouncer(t, time.Hour, func(context.Context, int) error { return wantErr },
		debounce.WithOnError(func(error) { t.Error("disabled error callback ran") }),
		debounce.WithOnError(nil),
	)
	require.NoError(t, d.Trigger(1))
	require.ErrorIs(t, d.Flush(t.Context()), wantErr)
}
