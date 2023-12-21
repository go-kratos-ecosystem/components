package crontab

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/transport"
	"github.com/robfig/cron/v3"
)

var _ transport.Server = (*Server)(nil)

type Server struct {
	cron *cron.Cron

	name  string
	mutex Mutex

	running   bool
	runningMu sync.Mutex

	stoped chan struct{}

	debug bool
}

type Option func(*Server)

func WithName(name string) Option {
	return func(s *Server) {
		s.name = name
	}
}

func WithMutex(m Mutex) Option {
	return func(s *Server) {
		s.mutex = m
	}
}

func WithDebug() Option {
	return func(s *Server) {
		s.debug = true
	}
}

func NewServer(c *cron.Cron, opts ...Option) *Server {
	s := &Server{
		cron:   c,
		name:   "cron:server",
		stoped: make(chan struct{}),
	}

	for _, opt := range opts {
		opt(s)
	}

	if s.mutex == nil {
		panic("crontab: mutex is nil")
	}

	return s
}

func (s *Server) Start(ctx context.Context) error {
	go s.run(ctx)

	return nil
}

func (s *Server) run(ctx context.Context) {
	timer := time.NewTicker(time.Second)
	defer func() {
		_ = s.mutex.Unlock(ctx, s.name)
		timer.Stop()
	}()

	for {
		select {
		case <-ctx.Done():
			s.log("crontab: server done")
			return
		case <-s.stoped:
			s.log("crontab: server stoped")
			return
		case <-timer.C:
			if err := s.mutex.Lock(ctx, s.name); err != nil {
				s.log(err)
				continue
			}

			s.start(ctx)
		}
	}
}

func (s *Server) start(ctx context.Context) {
	s.runningMu.Lock()
	defer s.runningMu.Unlock()

	if s.running {
		return
	}

	s.running = true
	s.cron.Start()

	s.log("crontab: server started")
}

func (s *Server) Stop(ctx context.Context) error {
	s.runningMu.Lock()
	defer s.runningMu.Unlock()

	if !s.running {
		return nil
	}

	s.running = false
	s.cron.Stop()

	close(s.stoped)

	return s.mutex.Unlock(ctx, s.name)
}

func (s *Server) log(v ...interface{}) {
	if s.debug {
		log.Println(v...)
	}
}
