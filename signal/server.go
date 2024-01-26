package signal

import (
	"context"
	"os"
	"os/signal"

	"github.com/go-kratos/kratos/v2/log"
)

var DefaultRecoveryHandler = func(err interface{}, sig os.Signal, handler Handler) {
	log.Errorf("[Signal] handler panic (%s): %v", sig, err)
}

type Server struct {
	handlers        []Handler
	stoped          chan struct{}
	recoveryHandler func(interface{}, os.Signal, Handler)
}

type Option func(*Server)

func AddHandler(handler ...Handler) Option {
	return func(s *Server) {
		s.handlers = append(s.handlers, handler...)
	}
}

func WithRecoveryHandler(handler func(interface{}, os.Signal, Handler)) Option {
	return func(s *Server) {
		if handler != nil {
			s.recoveryHandler = handler
		}
	}
}

func NewServer(opts ...Option) *Server {
	server := &Server{
		handlers: make([]Handler, 0),
		stoped:   make(chan struct{}),
	}

	for _, opt := range opts {
		opt(server)
	}

	if server.recoveryHandler == nil {
		server.recoveryHandler = DefaultRecoveryHandler
	}

	return server
}

func (s *Server) Start(ctx context.Context) error {
	var (
		signals  = make([]os.Signal, 0)
		handlers = make(map[os.Signal][]Handler)
	)

	for _, h := range s.handlers {
		for _, sig := range h.Listen() {
			if _, ok := handlers[sig]; !ok {
				handlers[sig] = make([]Handler, 0)
			}
			handlers[sig] = append(handlers[sig], h)
		}
		signals = append(signals, h.Listen()...)
	}

	signals = s.uniqueSignals(signals)

	ch := make(chan os.Signal, len(signals))
	signal.Notify(ch, signals...)

	log.Infof("[Signal] server starting")

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-s.stoped:
			return nil
		case sig := <-ch:
			if hs, ok := handlers[sig]; ok {
				for _, h := range hs {
					// if Support AsyncFeature
					if async, ok := h.(AsyncFeature); ok && async.Async() {
						go s.handle(sig, h)
					} else {
						s.handle(sig, h)
					}
				}
			}
		}
	}
}

func (s *Server) Register(handlers ...Handler) {
	s.handlers = append(s.handlers, handlers...)
}

func (s *Server) Stop(context.Context) error {
	log.Infof("[Signal] server stopping")
	close(s.stoped)
	return nil
}

func (s *Server) handle(sig os.Signal, handler Handler) {
	defer func() {
		if err := recover(); err != nil {
			s.recoveryHandler(err, sig, handler)
		}
	}()

	handler.Handle(sig)
}

func (s *Server) uniqueSignals(signals []os.Signal) []os.Signal {
	m := make(map[os.Signal]struct{})
	for _, sig := range signals {
		m[sig] = struct{}{}
	}
	signals = make([]os.Signal, 0)
	for sig := range m {
		signals = append(signals, sig)
	}
	return signals
}
