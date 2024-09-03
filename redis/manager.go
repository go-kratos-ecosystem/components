package redis

import "github.com/redis/go-redis/v9"

type Manager struct {
	redis.UniversalClient

	connections map[string]redis.UniversalClient
}

func New(db redis.UniversalClient) *Manager {
	return &Manager{
		UniversalClient: db,
		connections:     make(map[string]redis.UniversalClient),
	}
}

func (m *Manager) Register(name string, db redis.UniversalClient) {
	m.connections[name] = db
}

func (m *Manager) Conn(names ...string) redis.UniversalClient {
	var name string
	if len(names) > 0 {
		name = names[0]
	}

	if name == "" {
		return m.UniversalClient
	}

	if c, ok := m.connections[name]; ok {
		return c
	}

	panic("redis: the connection [" + name + "] is not registered.")
}
