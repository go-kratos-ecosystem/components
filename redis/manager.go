package redis

import "github.com/redis/go-redis/v9"

type Manager struct {
	redis.Cmdable

	connections map[string]redis.Cmdable
}

func New(db redis.Cmdable) *Manager {
	return &Manager{
		Cmdable:     db,
		connections: make(map[string]redis.Cmdable),
	}
}

func (m *Manager) Register(name string, db redis.Cmdable) {
	m.connections[name] = db
}

func (m *Manager) Conn(names ...string) redis.Cmdable {
	var name string
	if len(names) > 0 {
		name = names[0]
	}

	if name == "" {
		return m.Cmdable
	}

	if c, ok := m.connections[name]; ok {
		return c
	}

	panic("redis: the connection [" + name + "] is not registered.")
}
