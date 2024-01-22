package logging

import (
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/assert"
)

func TestManager_Channel(t *testing.T) {
	logger := New(log.DefaultLogger)
	logger.Register("ts", log.With(
		log.DefaultLogger,
		"ts", log.Timestamp(time.RFC3339),
	))

	logger.Log(log.LevelDebug, "test", "test") //nolint:errcheck
	logger.Channel().Log(log.LevelDebug, "test", "test")
	logger.Channel("ts").Log(log.LevelDebug, "test", "test")

	assert.Panics(t, func() {
		logger.Channel("unknown").Log(log.LevelDebug, "test", "test") //nolint:errcheck
	})
}

func TestManager_Log(t *testing.T) {
	m := New(nil)

	assert.EqualError(t, m.Log(log.LevelDebug, "test", "test"), ErrInvalidLogger.Error())
}
