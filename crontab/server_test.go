package crontab

import (
	"context"
	"testing"
	"time"

	"github.com/robfig/cron/v3"
)

func TestCrontab(t *testing.T) {
	var (
		ctx  = context.Background()
		data = make(chan string, 1)
		srv  = NewServer(cron.New(
			cron.WithSeconds(),
		))
	)

	_, _ = srv.AddFunc("* * * * * *", func() {
		data <- "Hello World!"
	})

	go srv.Start(ctx)   //nolint:errcheck
	defer srv.Stop(ctx) //nolint:errcheck

	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	select {
	case <-ctx.Done():
		t.Error("Crontab test timeout")
		return
	case msg := <-data:
		if msg != "Hello World!" {
			t.Errorf("Crontab test failed: %s", msg)
		}
	}
}
