package recovery

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecovery(t *testing.T) {
	testName := "test"
	testRequest := "request"
	testError := "error"

	recovery := NewRecovery(
		Handler(func(_ context.Context, name string, request any, err any) error {
			assert.Equal(t, testName, name)
			assert.Equal(t, testRequest, request)
			assert.Equal(t, testError, err)
			return fmt.Errorf("name: %s, request: %v, error: %v", name, request, err)
		}),
	)

	response, err := recovery(func(context.Context, string, any) (response any, err error) {
		panic(testError)
	})(context.Background(), testName, testRequest)
	assert.Error(t, err)
	assert.Nil(t, response)
}

func TestRecovery_DefaultHandler(t *testing.T) {
	recovery := NewRecovery()

	response, err := recovery(func(context.Context, string, any) (response any, err error) {
		panic("error")
	})(context.Background(), "test", "request")
	assert.Error(t, err)
	assert.Nil(t, response)
}
