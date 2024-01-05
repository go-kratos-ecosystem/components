package manager

import (
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/assert"
)

func TestManager_Channel(t *testing.T) {
	m := New(&Config{
		Default: "test",
		Channels: map[string]func() log.Logger{
			"test": func() log.Logger {
				return log.DefaultLogger
			},
			"ts": func() log.Logger {
				return log.With(
					log.DefaultLogger,
					"ts", log.Timestamp(time.RFC3339),
				)
			},
		},
	})

	m.Log(log.LevelDebug, "test", "test")                 //nolint:errcheck
	m.Channel("test").Log(log.LevelDebug, "test", "test") //nolint:errcheck
	m.Channel("ts").Log(log.LevelDebug, "test", "test")   //nolint:errcheck

	assert.Panics(t, func() {
		m.Channel("unknown").Log(log.LevelDebug, "test", "test") //nolint:errcheck
	})
}

func TestManager_Log(t *testing.T) {
	m := New(&Config{})

	assert.EqualError(t, m.Log(log.LevelDebug, "test", "test"), ErrNoDefaultLogger.Error())
}
