package cache

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/go-packagist/go-kratos-components/contracts/cache"
	"github.com/stretchr/testify/assert"
)

type mockStore struct {
	items map[string]interface{}
	rw    sync.RWMutex
}

func newMockStore() cache.Store {
	return &mockStore{
		items: make(map[string]interface{}),
	}
}

func (m *mockStore) Has(ctx context.Context, key string) bool {
	m.rw.RLock()
	defer m.rw.RUnlock()

	_, ok := m.items[key]

	return ok
}

func (m *mockStore) Get(ctx context.Context, key string, dest interface{}) error {
	m.rw.RLock()
	defer m.rw.RUnlock()

	if item, ok := m.items[key]; ok {
		return valueOf(item, dest)
	}

	return cache.ErrKeyNotFound
}

func (m *mockStore) Put(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	m.rw.Lock()
	defer m.rw.Unlock()

	m.items[key] = value

	return nil
}

func (m *mockStore) Increment(ctx context.Context, key string, value int) (int, error) {
	panic("implement me")
}

func (m *mockStore) Decrement(ctx context.Context, key string, value int) (int, error) {
	panic("implement me")
}

func (m *mockStore) Forever(ctx context.Context, key string, value interface{}) error {
	panic("implement me")
}

func (m *mockStore) Forget(ctx context.Context, key string) error {
	panic("implement me")
}

func (m *mockStore) Flush(ctx context.Context) error {
	panic("implement me")
}

func (m *mockStore) GetPrefix() string {
	panic("implement me")
}

var _ cache.Store = (*mockStore)(nil)

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
