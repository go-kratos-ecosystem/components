package jet

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTestMiddleware(t *testing.T, result *[]string, no int) Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, name string, request any) (any, error) {
			*result = append(*result, fmt.Sprintf("Before: %d", no))
			defer func() {
				*result = append(*result, fmt.Sprintf("After: %d", no))
			}()

			assert.Equal(t, "name", name)
			assert.Equal(t, "request", request)
			return next(ctx, name, request)
		}
	}
}

func TestMiddleware_Chain(t *testing.T) {
	var result []string
	chain := Chain(
		createTestMiddleware(t, &result, 1),
		createTestMiddleware(t, &result, 2),
	)(func(_ context.Context, name string, request any) (any, error) {
		result = append(result, "Before: 3")
		defer func() {
			result = append(result, "After: 3")
		}()

		assert.Equal(t, "name", name)
		assert.Equal(t, "request", request)

		return "response", nil
	})
	response, err := chain(context.Background(), "name", "request")
	assert.NoError(t, err)
	assert.Equal(t, "response", response)

	assert.Equal(t, []string{
		"Before: 1",
		"Before: 2",
		"Before: 3",
		"After: 3",
		"After: 2",
		"After: 1",
	}, result)
}
