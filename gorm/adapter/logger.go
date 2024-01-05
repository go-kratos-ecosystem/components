package adapter

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type loggerAdapter struct {
	log.Logger
	logger.Config

	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

func NewLogger(l log.Logger, config logger.Config) logger.Interface {
	var (
		infoStr      = "%s\n[info] "
		warnStr      = "%s\n[warn] "
		errStr       = "%s\n[error] "
		traceStr     = "%s\n[%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
	)

	if config.Colorful {
		infoStr = logger.Green + "%s\n" + logger.Reset + logger.Green + "[info] " + logger.Reset
		warnStr = logger.BlueBold + "%s\n" + logger.Reset + logger.Magenta + "[warn] " + logger.Reset
		errStr = logger.Magenta + "%s\n" + logger.Reset + logger.Red + "[error] " + logger.Reset
		traceStr = logger.Green + "%s\n" + logger.Reset + logger.Yellow + "[%.3fms] " + logger.BlueBold + "[rows:%v]" + logger.Reset + " %s"
		traceWarnStr = logger.Green + "%s " + logger.Yellow + "%s\n" + logger.Reset + logger.RedBold + "[%.3fms] " + logger.Yellow + "[rows:%v]" + logger.Magenta + " %s" + logger.Reset
		traceErrStr = logger.RedBold + "%s " + logger.MagentaBold + "%s\n" + logger.Reset + logger.Yellow + "[%.3fms] " + logger.BlueBold + "[rows:%v]" + logger.Reset + " %s"
	}

	return &loggerAdapter{
		Logger:       l,
		Config:       config,
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
}

func (l *loggerAdapter) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level

	return &newLogger
}

func (l *loggerAdapter) Info(_ context.Context, s string, i ...interface{}) {
	if l.LogLevel >= logger.Info {
		_ = l.Log(log.LevelInfo, fmt.Sprintf(l.infoStr+s, append([]interface{}{utils.FileWithLineNum()}, i...)...))
	}
}

func (l *loggerAdapter) Warn(_ context.Context, s string, i ...interface{}) {
	if l.LogLevel >= logger.Warn {
		_ = l.Log(log.LevelWarn, fmt.Sprintf(l.warnStr+s, append([]interface{}{utils.FileWithLineNum()}, i...)...))
	}
}

func (l *loggerAdapter) Error(_ context.Context, s string, i ...interface{}) {
	if l.LogLevel >= logger.Error {
		_ = l.Log(log.LevelError, fmt.Sprintf(l.errStr+s, append([]interface{}{utils.FileWithLineNum()}, i...)...))
	}
}

func (l *loggerAdapter) Trace(_ context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) { //nolint:lll
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= logger.Error && (!errors.Is(err, logger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError): //nolint:lll
		sql, rows := fc()
		if rows == -1 {
			_ = l.Log(log.LevelFatal, fmt.Sprintf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)) //nolint:gomnd,lll
		} else {
			_ = l.Log(log.LevelFatal, fmt.Sprintf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)) //nolint:gomnd,lll
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			_ = l.Log(log.LevelWarn, fmt.Sprintf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)) //nolint:gomnd,lll
		} else {
			_ = l.Log(log.LevelWarn, fmt.Sprintf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)) //nolint:gomnd,lll
		}
	case l.LogLevel == logger.Info:
		sql, rows := fc()
		if rows == -1 {
			_ = l.Log(log.LevelInfo, fmt.Sprintf(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)) //nolint:gomnd,lll
		} else {
			_ = l.Log(log.LevelInfo, fmt.Sprintf(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)) //nolint:gomnd,lll
		}
	}
}
