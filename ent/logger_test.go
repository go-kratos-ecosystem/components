package ent

import (
	"testing"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	logger := NewLogger(log.DefaultLogger, log.LevelInfo)
	assert.NotPanics(t, func() {
		logger("test")
	})
}
