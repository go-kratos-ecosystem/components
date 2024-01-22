package logging

import (
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
)

var ErrInvalidLogger = fmt.Errorf("logging: invalid logger")

type Logger struct {
	log.Logger
	loggers map[string]log.Logger
}

func New(logger log.Logger) *Logger {
	return &Logger{
		Logger:  logger,
		loggers: make(map[string]log.Logger),
	}
}

func (l *Logger) Log(level log.Level, keyvals ...interface{}) error {
	if l.Logger != nil {
		return l.Logger.Log(level, keyvals...)
	}

	return ErrInvalidLogger
}

func (l *Logger) Register(name string, logger log.Logger) {
	l.loggers[name] = logger
}

func (l *Logger) Channel(names ...string) log.Logger {
	if len(names) <= 0 {
		return l.Logger
	}

	name := names[0]

	if logger, ok := l.loggers[name]; ok {
		return logger
	}

	panic(fmt.Errorf("logging: unknown logger %s", name))
}
