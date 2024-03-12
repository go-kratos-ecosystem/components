package foundation

import (
	"github.com/go-kratos-ecosystem/components/v2/coordinator"
	"github.com/go-kratos-ecosystem/components/v2/event"
)

type options struct {
	m *coordinator.Manager
	d *event.Dispatcher
}

type Option func(*options)

func WithManager(m *coordinator.Manager) Option {
	return func(o *options) {
		o.m = m
	}
}

func WithDispatcher(d *event.Dispatcher) Option {
	return func(o *options) {
		o.d = d
	}
}
