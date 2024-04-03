package crontab

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
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
	log.Info("[Crontab] server starting")
	s.Cron.Run()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	log.Info("[Crontab] server stopping")

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-s.Cron.Stop().Done():
			return nil
		}
	}
}
