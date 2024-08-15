package logging

import (
	"context"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/assert"
)

type mockLogger struct {
	t *testing.T
}

func newMockLogger(t *testing.T) log.Logger {
	return &mockLogger{t: t}
}

func (m *mockLogger) Log(level log.Level, keyvals ...interface{}) error {
	if !assert.Len(m.t, keyvals, 14) {
		return nil
	}

	assert.Equal(m.t, "kind", keyvals[0])
	assert.Equal(m.t, "jet", keyvals[1])
	assert.Equal(m.t, "service", keyvals[2])
	assert.Equal(m.t, "{XXX}", keyvals[3])
	assert.Equal(m.t, "name", keyvals[4])

	if keyvals[5] == "no-error" {
		assert.Equal(m.t, log.LevelInfo, level)
		assert.Equal(m.t, "response", keyvals[9])
		assert.Nil(m.t, keyvals[11])
	} else if keyvals[5] == "with-error" {
		assert.Equal(m.t, log.LevelError, level)
		assert.Nil(m.t, keyvals[9])
		assert.Equal(m.t, assert.AnError, keyvals[11])
	} else {
		assert.Fail(m.t, "unexpected name")
	}

	assert.Equal(m.t, "request", keyvals[6])
	assert.Nil(m.t, keyvals[7])
	assert.Equal(m.t, "response", keyvals[8])
	assert.Equal(m.t, "error", keyvals[10])
	assert.Equal(m.t, "latency", keyvals[12])
	assert.IsType(m.t, time.Duration(0), keyvals[13])
	return nil
}

var _ log.Logger = (*mockLogger)(nil)

func TestLogging(t *testing.T) {
	logging := New(
		Logger(newMockLogger(t)),
	)

	// no error
	_, _ = logging(func(ctx context.Context, name string, request any) (response any, err error) {
		return "response", nil
	})(context.Background(), "no-error", nil)

	// with error
	_, _ = logging(func(ctx context.Context, name string, request any) (response any, err error) {
		return nil, assert.AnError
	})(context.Background(), "with-error", nil)
}
