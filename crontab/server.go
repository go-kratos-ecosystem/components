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

func (s *Server) Stop(context.Context) error {
	log.Info("[Crontab] server stopping")
	<-s.Cron.Stop().Done()

	return nil
}
