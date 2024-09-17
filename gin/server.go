package gin

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
)

type Server struct {
	*gin.Engine

	server      *http.Server
	middlewares []gin.HandlerFunc

	addr string
}

type Option func(*Server)

// Deprecated: use Addr
func WithAddr(addr string) Option {
	return Addr(addr)
}

func Addr(addr string) Option {
	return func(s *Server) {
		s.addr = addr
	}
}

func Middleware(middlewares ...gin.HandlerFunc) Option {
	return func(s *Server) {
		s.middlewares = append(s.middlewares, middlewares...)
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

	// apply middlewares
	if len(srv.middlewares) > 0 {
		srv.Use(srv.middlewares...)
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
