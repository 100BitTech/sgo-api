package db

import (
	"context"
	"errors"
	"fmt"
	"sgo-api/base"
	"time"

	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

const (
	infoFormat       = "%s\n"
	warnFormat       = "%s\n"
	errorFormat      = "%s\n"
	traceFormat      = "%s\n[%.3fms] [rows:%v] %s"
	traceWarnFormat  = "%s\n%s\n[%.3fms] [rows:%v] %s"
	traceErrorFormat = "%s\n%s\n[%.3fms] [rows:%v] %s"
)

type Logger struct {
	log  base.Logger
	conf LoggerConfig
}

type LoggerConfig struct {
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
	TraceContextKey           any
}

func NewLogger(log base.Logger, conf LoggerConfig) Logger {
	if conf.SlowThreshold <= 0 {
		conf.SlowThreshold = 3000 * time.Millisecond
	}

	return Logger{
		log:  log.WithTag("GORM"),
		conf: conf,
	}
}

func (l Logger) Log() base.Logger {
	return l.log
}

func (l Logger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

func (l Logger) Info(ctx context.Context, msg string, data ...any) {
	log := l.log.WithTrace(ctx, l.conf.TraceContextKey)
	log.Infof(infoFormat+msg, append([]any{utils.FileWithLineNum()}, data...)...)
}

func (l Logger) Warn(ctx context.Context, msg string, data ...any) {
	log := l.log.WithTrace(ctx, l.conf.TraceContextKey)
	log.Warnf(warnFormat+msg, append([]any{utils.FileWithLineNum()}, data...)...)
}

func (l Logger) Error(ctx context.Context, msg string, data ...any) {
	log := l.log.WithTrace(ctx, l.conf.TraceContextKey)
	log.Errorf(errorFormat+msg, append([]any{utils.FileWithLineNum()}, data...)...)
}

func (l Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	log := l.log.WithTrace(ctx, l.conf.TraceContextKey)

	elapsed := time.Since(begin)
	switch {
	case err != nil && (!errors.Is(err, logger.ErrRecordNotFound) || !l.conf.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			log.Errorf(traceErrorFormat, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			log.Errorf(traceErrorFormat, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > l.conf.SlowThreshold && l.conf.SlowThreshold != 0:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.conf.SlowThreshold)
		if rows == -1 {
			log.Warnf(traceWarnFormat, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			log.Warnf(traceWarnFormat, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	default:
		sql, rows := fc()
		if rows == -1 {
			log.Debugf(traceFormat, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			log.Debugf(traceFormat, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}
