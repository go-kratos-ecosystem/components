package coordinator

import "sync"

type Manager struct {
	coordinators map[string]*Coordinator
	mu           sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		coordinators: make(map[string]*Coordinator),
	}
}

func (m *Manager) Until(identifier string) *Coordinator {
	m.mu.Lock()
	defer m.mu.Unlock()

	if c, ok := m.coordinators[identifier]; ok {
		return c
	}

	c := NewCoordinator()
	m.coordinators[identifier] = c

	return c
}

func (m *Manager) Close(identifier string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if c, ok := m.coordinators[identifier]; ok {
		c.Close()
		delete(m.coordinators, identifier)
	}
}

func (m *Manager) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, c := range m.coordinators {
		c.Close()
	}

	m.coordinators = make(map[string]*Coordinator)
}

var defaultManager = NewManager()

func Until(identifier string) *Coordinator {
	return defaultManager.Until(identifier)
}

func Close(identifier string) {
	defaultManager.Close(identifier)
}

func Clear() {
	defaultManager.Clear()
}
