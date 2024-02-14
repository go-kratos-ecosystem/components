package v2

import (
	"context"
	"testing"
	"time"

	"github.com/robfig/cron/v3"
)

type mockJob struct{}

func RegisterMockJob(srv *Server) {
	n := newMockJob()
	_, _ = srv.AddJob(n.Exp(), n)
}

func newMockJob() *mockJob {
	return &mockJob{}
}

func (j *mockJob) Exp() string {
	return "* * * * * *"
}

func (j *mockJob) Run() {
	data <- "Hello World!"
}

var (
	ctx  = context.Background()
	data = make(chan string, 1)
)

func TestCrontab(t *testing.T) {
	srv := NewServer(cron.New(
		cron.WithSeconds(),
	))

	RegisterMockJob(srv)

	go srv.Start(ctx)
	defer srv.Stop(ctx)

	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	select {
	case <-ctx.Done():
		t.Error("Crontab test timeout")
		return
	case msg := <-data:
		t.Log(msg)
	}
}
