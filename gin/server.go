package gin

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
)

type Server struct {
	*gin.Engine

	server *http.Server

	addr string
}

type Option func(*Server)

func WithAddr(addr string) Option {
	return func(s *Server) {
		s.addr = addr
	}
}

func NewServer(e *gin.Engine, opts ...Option) *Server {
	srv := &Server{
		Engine: e,
		addr:   ":8080",
	}

	for _, opt := range opts {
		opt(srv)
	}

	srv.server = &http.Server{
		Addr:    srv.addr,
		Handler: e,
	}

	return srv
}

func (s *Server) Start(_ context.Context) error {
	log.Infof("[GIN] server listening on: %s", s.addr)
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	log.Info("[GIN] server stopping")
	return s.server.Shutdown(ctx)
}
