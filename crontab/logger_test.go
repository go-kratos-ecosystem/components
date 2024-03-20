package crontab

import (
	"fmt"
	"testing"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/assert"
)

type testLogger struct {
	ch chan string
}

func newTestLogger() *testLogger {
	return &testLogger{
		ch: make(chan string, 1),
	}
}

func (t *testLogger) Log(level log.Level, keyvals ...interface{}) error {
	t.ch <- fmt.Sprintf("level: %s, keyvals: %v", level, keyvals)
	return nil
}

func TestLogger(t *testing.T) {
	_, ok := interface{}(NewLogger(nil)).(interface{ Printf(string, ...interface{}) })
	assert.True(t, ok)

	logger := newTestLogger()
	l := NewLogger(logger)
	l.Printf("test %s", "logger")
	assert.Equal(t, "level: INFO, keyvals: [msg test logger]", <-logger.ch)
}
