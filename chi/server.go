package chi

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-kratos/kratos/v2/log"
)

type Server struct {
	*chi.Mux
	server *http.Server
	addr   string
}

type Option func(*Server)

func Addr(addr string) Option {
	return func(s *Server) {
		s.addr = addr
	}
}

func NewServer(c *chi.Mux, opts ...Option) *Server {
	srv := &Server{
		Mux:  c,
		addr: ":8080",
	}

	for _, opt := range opts {
		opt(srv)
	}

	srv.server = &http.Server{
		Addr:    srv.addr,
		Handler: c,
	}

	return srv
}

func (s *Server) Start(_ context.Context) error {
	log.Infof("[go-chi] server listening on: %s", s.addr)
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	log.Info("[go-chi] server stopping")
	return s.server.Shutdown(ctx)
}
