package context

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContext(t *testing.T) {
	type (
		mockProviderStruct1 struct{}
		mockProviderStruct2 struct{}
		mockProviderStruct3 struct{}
	)

	var (
		mockProvider1 = func(ctx context.Context) (context.Context, error) {
			return context.WithValue(ctx, mockProviderStruct1{}, "mockProvider1"), nil
		}

		mockProvider2 = func(ctx context.Context) (context.Context, error) {
			return context.WithValue(ctx, mockProviderStruct2{}, "mockProvider2"), nil
		}

		mockProvider3 = func(ctx context.Context) (context.Context, error) {
			return ctx, errors.New("mockProvider3")
		}
	)

	ctx1, err1 := Chain(
		context.Background(),
		mockProvider1, mockProvider2,
	)
	assert.NoError(t, err1)
	assert.Equal(t, "mockProvider1", ctx1.Value(mockProviderStruct1{}))
	assert.Equal(t, "mockProvider2", ctx1.Value(mockProviderStruct2{}))

	ctx2, err2 := Chain(
		context.Background(),
		mockProvider1, mockProvider3,
	)
	assert.Error(t, err2)
	assert.NotNil(t, ctx2)
	assert.Equal(t, "mockProvider1", ctx2.Value(mockProviderStruct1{}))
	assert.Nil(t, ctx2.Value(mockProviderStruct3{}))
}
