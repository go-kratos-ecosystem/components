package batcher_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-fries/fries/batcher/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewValidatesInputs(t *testing.T) {
	handler := func(context.Context, []int) error { return nil }
	t.Run("nil context", func(t *testing.T) {
		var nilContext context.Context
		b, err := batcher.New(nilContext, handler)
		require.ErrorIs(t, err, batcher.ErrInvalidContext)
		assert.Nil(t, b)
	})
	t.Run("canceled context", func(t *testing.T) {
		ctx, cancel := context.WithCancelCause(t.Context())
		wantErr := errors.New("stopped")
		cancel(wantErr)
		b, err := batcher.New(ctx, handler)
		require.ErrorIs(t, err, wantErr)
		assert.Nil(t, b)
	})
	t.Run("nil handler", func(t *testing.T) {
		b, err := batcher.New[int](t.Context(), nil)
		require.ErrorIs(t, err, batcher.ErrNilHandler)
		assert.Nil(t, b)
	})
	for name, option := range map[string]batcher.Option{
		"zero batch":        batcher.WithBatchSize(0),
		"negative batch":    batcher.WithBatchSize(-1),
		"zero interval":     batcher.WithFlushInterval(0),
		"negative interval": batcher.WithFlushInterval(-time.Second),
		"negative queue":    batcher.WithQueueSize(-1),
	} {
		t.Run(name, func(t *testing.T) {
			b, err := batcher.New(t.Context(), handler, option)
			require.ErrorIs(t, err, batcher.ErrInvalidOption)
			assert.Nil(t, b)
		})
	}
	t.Run("nil options and callback", func(t *testing.T) {
		b, err := batcher.New(t.Context(), handler, nil, batcher.WithOnResult(nil))
		require.NoError(t, err)
		require.NoError(t, b.Shutdown(t.Context()))
	})
}
