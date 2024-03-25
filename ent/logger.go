package ent

import (
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
)

type Logger func(...any)

func NewLogger(logger log.Logger, levels ...log.Level) Logger {
	level := log.LevelDebug
	if len(levels) > 0 {
		level = levels[0]
	}

	return func(args ...any) {
		_ = logger.Log(level, "msg", fmt.Sprint(args...))
	}
}
