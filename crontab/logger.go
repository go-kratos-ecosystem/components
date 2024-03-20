package crontab

import (
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
)

type Logger struct {
	logger log.Logger
}

func NewLogger(logger log.Logger) *Logger {
	return &Logger{logger: logger}
}

func (l *Logger) Printf(format string, v ...interface{}) {
	_ = l.logger.Log(log.LevelInfo, "msg", fmt.Sprintf(format, v...))
}
