package cache

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/go-packagist/go-kratos-components/contract/cache"
	"github.com/stretchr/testify/assert"
)

var ctx = context.Background()

type mockStore struct {
	items map[string]interface{}
	rw    sync.RWMutex
}

func newMockStore() cache.Store {
	return &mockStore{
		items: make(map[string]interface{}),
	}
}

func (m *mockStore) Has(ctx context.Context, key string) (bool, error) {
	m.rw.RLock()
	defer m.rw.RUnlock()

	_, ok := m.items[key]

	return ok, nil
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

func newMockRepository() cache.Repository {
	return NewRepository(newMockStore())
}

func TestRepository_Has(t *testing.T) {
	r := newMockRepository()

	r.Has(ctx, "test")
}

func TestRepository_Add(t *testing.T) {
	r := newMockRepository()

	added, err := r.Add(ctx, "test", 1, time.Second*30)
	assert.NoError(t, err)
	assert.True(t, added)

	added2, err2 := r.Add(ctx, "test", 1, time.Second*30)
	assert.NoError(t, err2)
	assert.False(t, added2)
}

func TestRepository_Remember(t *testing.T) {
	r := newMockRepository()

	var value string
	err1 := r.Remember(ctx, "remember", &value, func() interface{} {
		return "test"
	}, time.Second*10)

	assert.NoError(t, err1)
	assert.Equal(t, "test", value)

	err2 := r.Remember(ctx, "remember", &value, func() interface{} {
		return "test2"
	}, time.Second*10)
	assert.NoError(t, err2)
	assert.Equal(t, "test", value)
}
