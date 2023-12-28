package udp

import (
	"context"
	"log"
	"net"
	"sync"
)

type Server struct {
	address string

	bufSize int

	conn net.PacketConn
	mu   sync.Mutex // guards conn

	handler func(conn net.PacketConn, buf []byte, addr net.Addr)

	recoveryHandler func(conn net.PacketConn, buf []byte, addr net.Addr, err interface{})
}

type Option func(*Server)

func WithBufSize(bufSize int) Option {
	return func(s *Server) {
		if bufSize > 0 {
			s.bufSize = bufSize
		}
	}
}

func WithHandler(handler func(conn net.PacketConn, buf []byte, addr net.Addr)) Option {
	return func(s *Server) {
		if handler != nil {
			s.handler = handler
		}
	}
}

func WithRecoveryHandler(handler func(conn net.PacketConn, buf []byte, addr net.Addr, err interface{})) Option {
	return func(s *Server) {
		if handler != nil {
			s.recoveryHandler = handler
		}
	}
}

func NewServer(address string, opts ...Option) *Server {
	s := &Server{
		address: address,
		bufSize: 1024,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Server) Start(ctx context.Context) (err error) {
	s.conn, err = net.ListenPacket("udp", s.address)
	if err != nil {
		return
	}

	log.Printf("udp server: listening on %s\n", s.address)

	buf := make([]byte, s.bufSize)

	for {
		n, addr, err := s.conn.ReadFrom(buf)
		if err != nil {
			return err
		}

		if s.handler == nil {
			log.Printf("udp server: receive from %s: %s\n", addr.String(), string(buf))
			continue
		}

		go s.handle(buf[:n], addr)
	}

}

func (s *Server) handle(buf []byte, addr net.Addr) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.recoveryHandler != nil {
		defer func() {
			if err := recover(); err != nil {
				s.recoveryHandler(s.conn, buf, addr, err)
			}
		}()
	}

	s.handler(s.conn, buf, addr)
}

func (s *Server) Stop(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	log.Println("udp server: stopping")

	return s.conn.Close()
}
