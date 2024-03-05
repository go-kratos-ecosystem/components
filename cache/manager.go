package cache

type Manager struct {
	Repository

	drivers map[string]Repository
}

func NewManager(repository Repository) *Manager {
	return &Manager{
		Repository: repository,
		drivers:    make(map[string]Repository),
	}
}

func (m *Manager) Register(name string, repository Repository) {
	m.drivers[name] = repository
}

func (m *Manager) Driver(names ...string) Repository {
	var name string
	if len(names) > 0 {
		name = names[0]
	}

	if name == "" {
		return m.Repository
	}

	if c, ok := m.drivers[name]; ok {
		return c
	}

	panic("cache: the repository [" + name + "] is not registered.")
}
