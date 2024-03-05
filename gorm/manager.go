package gorm

import "gorm.io/gorm"

type ConnectionFunc func() *gorm.DB

type Manager struct {
	*gorm.DB

	connections map[string]*gorm.DB
}

func NewManager(db *gorm.DB) *Manager {
	return &Manager{
		DB:          db,
		connections: make(map[string]*gorm.DB),
	}
}

func (m *Manager) Register(name string, db *gorm.DB) {
	m.connections[name] = db
}

func (m *Manager) Conn(names ...string) *gorm.DB {
	var name string
	if len(names) > 0 {
		name = names[0]
	}

	if name == "" {
		return m.DB
	}

	if c, ok := m.connections[name]; ok {
		return c
	}

	panic("The connection [" + name + "] is not registered.")
}
