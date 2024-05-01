package iris

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/kataras/iris/v12"
)

type Server struct {
	*iris.Application
	addr          string
	configurators []iris.Configurator
}

type Option func(*Server)

func Addr(addr string) Option {
	return func(s *Server) {
		s.addr = addr
	}
}

func WithConfigurators(configurators ...iris.Configurator) Option {
	return func(s *Server) {
		s.configurators = append(s.configurators, configurators...)
	}
}

func NewServer(app *iris.Application, opts ...Option) *Server {
	srv := &Server{
		Application: app,
		addr:        ":8080",
	}

	for _, opt := range opts {
		opt(srv)
	}

	return srv
}

func (s *Server) Start(_ context.Context) error {
	log.Infof("[iris] server listening on: %s", s.addr)
	return s.Application.Listen(s.addr, s.configurators...)
}

func (s *Server) Stop(ctx context.Context) error {
	log.Info("[iris] server stopping")
	return s.Application.Shutdown(ctx)
}
