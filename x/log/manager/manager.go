package manager

import (
	"errors"
	"sync"

	"github.com/go-kratos/kratos/v2/log"
)

var (
	ErrNoDefaultLogger = errors.New("log manager: no default logger")
)

type Config struct {
	Default string

	Channels map[string]func() log.Logger
}

type Manager struct {
	config *Config

	channels map[string]log.Logger // resolved logger by name
	rw       sync.RWMutex
}

var _ log.Logger = (*Manager)(nil)

func New(config *Config) *Manager {
	return &Manager{
		config:   config,
		channels: make(map[string]log.Logger),
	}
}

func (m *Manager) Log(level log.Level, keyvals ...interface{}) error {
	if m.config.Default != "" {
		return m.Channel(m.config.Default).Log(level, keyvals...)
	}

	return ErrNoDefaultLogger
}

func (m *Manager) Channel(name string) log.Logger {
	m.rw.RLock()

	if logger, ok := m.channels[name]; ok {
		m.rw.RUnlock()
		return logger
	}

	m.rw.RUnlock()

	m.rw.Lock()
	defer m.rw.Unlock()

	if logger, ok := m.config.Channels[name]; ok {
		l := logger()
		m.channels[name] = l

		return l
	}

	panic("log manager: unknown channel " + name)
}
