package recovery

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-kratos-ecosystem/components/v2/hyperf/jet"
)

func TestRecovery(t *testing.T) {
	testName := "test"
	testRequest := "request"
	testError := "error"

	recovery := New(
		Handler(func(_ context.Context, _ *jet.Client, name string, request any, err any) error {
			assert.Equal(t, testName, name)
			assert.Equal(t, testRequest, request)
			assert.Equal(t, testError, err)
			return fmt.Errorf("name: %s, request: %v, error: %v", name, request, err)
		}),
	)

	response, err := recovery(func(context.Context, *jet.Client, string, any) (response any, err error) {
		panic(testError)
	})(context.Background(), nil, testName, testRequest)
	assert.Error(t, err)
	assert.Nil(t, response)
}

func TestRecovery_DefaultHandler(t *testing.T) {
	recovery := New()

	response, err := recovery(func(context.Context, *jet.Client, string, any) (response any, err error) {
		panic("error")
	})(context.Background(), nil, "test", "request")
	assert.Error(t, err)
	assert.Nil(t, response)
}
