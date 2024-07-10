package gorm

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	gl "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type logger struct {
	log.Logger
	gl.Config

	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

func New(l log.Logger, config gl.Config) gl.Interface {
	var (
		infoStr      = "%s\n[info] "
		warnStr      = "%s\n[warn] "
		errStr       = "%s\n[error] "
		traceStr     = "%s\n[%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
	)

	if config.Colorful {
		infoStr = gl.Green + "%s\n" + gl.Reset + gl.Green + "[info] " + gl.Reset //nolint:goconst
		warnStr = gl.BlueBold + "%s\n" + gl.Reset + gl.Magenta + "[warn] " + gl.Reset
		errStr = gl.Magenta + "%s\n" + gl.Reset + gl.Red + "[error] " + gl.Reset
		traceStr = gl.Green + "%s\n" + gl.Reset + gl.Yellow + "[%.3fms] " + gl.BlueBold + "[rows:%v]" + gl.Reset + " %s"                                     //nolint:goconst,lll
		traceWarnStr = gl.Green + "%s " + gl.Yellow + "%s\n" + gl.Reset + gl.RedBold + "[%.3fms] " + gl.Yellow + "[rows:%v]" + gl.Magenta + " %s" + gl.Reset //nolint:lll
		traceErrStr = gl.RedBold + "%s " + gl.MagentaBold + "%s\n" + gl.Reset + gl.Yellow + "[%.3fms] " + gl.BlueBold + "[rows:%v]" + gl.Reset + " %s"       //nolint:lll
	}

	return &logger{
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

func (l *logger) LogMode(level gl.LogLevel) gl.Interface {
	newLogger := *l
	newLogger.LogLevel = level

	return &newLogger
}

func (l *logger) Info(_ context.Context, s string, i ...interface{}) {
	if l.LogLevel >= gl.Info {
		_ = l.Log(log.LevelInfo, fmt.Sprintf(l.infoStr+s, append([]interface{}{utils.FileWithLineNum()}, i...)...))
	}
}

func (l *logger) Warn(_ context.Context, s string, i ...interface{}) {
	if l.LogLevel >= gl.Warn {
		_ = l.Log(log.LevelWarn, fmt.Sprintf(l.warnStr+s, append([]interface{}{utils.FileWithLineNum()}, i...)...))
	}
}

func (l *logger) Error(_ context.Context, s string, i ...interface{}) {
	if l.LogLevel >= gl.Error {
		_ = l.Log(log.LevelError, fmt.Sprintf(l.errStr+s, append([]interface{}{utils.FileWithLineNum()}, i...)...))
	}
}

func (l *logger) Trace(_ context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) { //nolint:lll
	if l.LogLevel <= gl.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= gl.Error && (!errors.Is(err, gl.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError): //nolint:lll
		sql, rows := fc()
		if rows == -1 {
			_ = l.Log(log.LevelFatal, fmt.Sprintf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)) //nolint:gomnd,lll
		} else {
			_ = l.Log(log.LevelFatal, fmt.Sprintf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)) //nolint:gomnd,lll
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gl.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			_ = l.Log(log.LevelWarn, fmt.Sprintf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)) //nolint:gomnd,lll
		} else {
			_ = l.Log(log.LevelWarn, fmt.Sprintf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)) //nolint:gomnd,lll
		}
	case l.LogLevel == gl.Info:
		sql, rows := fc()
		if rows == -1 {
			_ = l.Log(log.LevelInfo, fmt.Sprintf(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)) //nolint:gomnd,lll
		} else {
			_ = l.Log(log.LevelInfo, fmt.Sprintf(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)) //nolint:gomnd,lll
		}
	}
}
