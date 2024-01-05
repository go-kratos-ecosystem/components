package crontab

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

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
	timer := time.NewTicker(time.Second)

	defer func() {
		if s.running {
			s.cron.Stop()
		}

		s.mutex.Unlock(ctx, s.name)
		timer.Stop()
	}()

	for {
		select {
		case <-ctx.Done():
			s.log("crontab: server done")
			return ctx.Err()
		case <-s.stoped:
			s.log("crontab: server stoped")
			return nil
		case <-timer.C:
			if err := s.mutex.Lock(ctx, s.name); err != nil {
				s.log(err)
				continue
			}

			s.start()
		}
	}
}

func (s *Server) start() {
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
	s.log("crontab: server stopping")

	close(s.stoped)

	return nil
}

func (s *Server) log(v ...interface{}) {
	if s.debug {
		log.Println(v...)
	}
}
