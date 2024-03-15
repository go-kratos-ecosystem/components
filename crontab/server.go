package crontab

import (
	"context"

	"github.com/robfig/cron/v3"
)

type Server struct {
	*cron.Cron
}

func NewServer(c *cron.Cron) *Server {
	return &Server{
		Cron: c,
	}
}

func (s *Server) Start(context.Context) error {
	s.Cron.Run()
	return nil
}

func (s *Server) Stop(context.Context) error {
	s.Cron.Stop()
	return nil
}
