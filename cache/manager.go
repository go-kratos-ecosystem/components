package cache

import (
	"sync"

	"github.com/go-packagist/go-kratos-components/contracts/cache"
)

type Config struct {
	Default string

	Stores map[string]cache.Store
}

type Manager struct {
	config *Config

	stores map[string]cache.Repository
	rw     sync.RWMutex
}

func NewManager(config *Config) *Manager {
	return &Manager{
		config: config,
		stores: make(map[string]cache.Repository),
	}
}

func (m *Manager) Connect(names ...string) cache.Repository {
	if len(names) == 0 {
		names = []string{m.config.Default}
	}

	name := names[0]

	m.rw.RLock()
	if store, ok := m.stores[name]; ok {
		m.rw.RUnlock()
		return store
	}
	m.rw.RUnlock()

	m.rw.Lock()
	defer m.rw.Unlock()
	m.stores[name] = m.resolve(name)

	return m.stores[name]
}

func (m *Manager) resolve(name string) cache.Repository {
	if store, ok := m.config.Stores[name]; ok {
		return NewRepository(store)
	}

	panic("cache: unknown store name: " + name)
}
