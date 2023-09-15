package cache

import (
	"testing"
	"time"

	"github.com/go-packagist/go-kratos-components/contract/cache"
	"github.com/stretchr/testify/assert"
)

func TestManager(t *testing.T) {
	m := NewManager(&Config{
		Default: "test1",
		Stores: map[string]cache.Store{
			"test1": newMockStore(),
			"test2": newMockStore(),
		},
	})

	var test1, test2, test3, test4 string

	// use default
	assert.NoError(t, m.Connect().Put(ctx, "test", "test", time.Second*10))
	assert.NoError(t, m.Connect().Get(ctx, "test", &test1))
	assert.Equal(t, "test", test1)

	// use test1
	assert.NoError(t, m.Connect("test1").Get(ctx, "test", &test2))
	assert.Equal(t, "test", test2)

	// use test2
	assert.Error(t, m.Connect("test2").Get(ctx, "test", &test3))
	assert.NotEqual(t, "test", test3)

	assert.NoError(t, m.Connect("test2").Put(ctx, "test", "test", time.Second*10))
	assert.NoError(t, m.Connect("test2").Get(ctx, "test", &test4))
	assert.Equal(t, "test", test4)

	// unknown
	assert.Panics(t, func() {
		m.Connect("unknown").Get(ctx, "test", &test3)
	})
}
