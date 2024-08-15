package retry

import (
	"context"
	"testing"
)

func TestRetry(t *testing.T) {
	handler := NewRetry(
		Attempts(3),
		Backoff(LinearBackoff(1)),
		Allow(AllowChain(DefaultAllow)),
	)(func(ctx context.Context, name string, request any) (response any, err error) {
		return nil, nil
	})

	_, _ = handler(context.Background(), "test", nil)
}
