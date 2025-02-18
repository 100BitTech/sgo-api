package base

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/samber/oops"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogLevel int8

const (
	LogLevelDebug LogLevel = iota - 1
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

type Logger interface {
	WithTrace(ctx context.Context, traceContextKey any) Logger
	WithTag(tag string) Logger

	Debug(args ...any)
	Info(args ...any)
	Warn(args ...any)
	Error(args ...any)

	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)

	// 日志使用了缓存，需要在结束时调用
	Sync()
}

type LogConfig struct {
	Level  int8   `json:"level"`   // 日志级别。默认为 Info
	Path   string `json:"path"`    // 日志路径。默认为 ./logs
	MaxAge int    `json:"max_age"` // 日志存活时间（小时）。默认为 24*30
}

// 初始化日志器
//   - 默认日志目录为 ./
//   - 可以通过入参指定目录（优先级最高）
//   - 可以编译时加上 -ldflags '-X github.com/100BitTech/sgo-api/base.LOGPATH=./logs' 指定目录
func InitLogger(conf LogConfig) Logger {
	if conf.Path == "" {
		conf.Path = LOGPATH
	}
	if conf.Path == "" {
		conf.Path = "./logs"
	}
	fmt.Printf("日志路径：%v\n", conf.Path)

	zap.ReplaceGlobals(newZapLogger(conf))
	return newMyZapLogger(zap.S())
}

func newZapLogger(conf LogConfig) *zap.Logger {
	return zap.New(zapcore.NewTee(
		newZapLoggerFileCore(conf),
		newZapLoggerConsoleCore(zapcore.Level(conf.Level)),
	), zap.AddCaller())
}

func newZapLoggerFileCore(conf LogConfig) zapcore.Core {
	filePath := strings.TrimRight(conf.Path, "/") + "/%Y-%m-%d.log"

	if conf.MaxAge <= 0 {
		conf.MaxAge = 30 * 24
	}

	encc := zap.NewProductionEncoderConfig()
	encc.EncodeTime = zapcore.ISO8601TimeEncoder
	encc.EncodeLevel = zapcore.CapitalLevelEncoder
	enc := zapcore.NewConsoleEncoder(encc)

	rotator, err := rotatelogs.New(
		filePath,
		rotatelogs.WithMaxAge(time.Duration(conf.MaxAge)*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour))
	if err != nil {
		panic(oops.Wrap(err))
	}

	ws := zapcore.AddSync(rotator)

	return zapcore.NewCore(enc, ws, zapcore.Level(conf.Level))
}

func newZapLoggerConsoleCore(level zapcore.Level) zapcore.Core {
	encc := zap.NewProductionEncoderConfig()
	encc.EncodeTime = zapcore.ISO8601TimeEncoder
	encc.EncodeLevel = zapcore.CapitalColorLevelEncoder
	enc := zapcore.NewConsoleEncoder(encc)

	ws := zapcore.AddSync(os.Stdout)

	return zapcore.NewCore(enc, ws, level)
}

type MyZapLogger struct {
	traceID string
	tag     string
	log     *zap.SugaredLogger
}

func newMyZapLogger(log *zap.SugaredLogger) Logger {
	return MyZapLogger{log: log}
}

func (l MyZapLogger) WithTrace(ctx context.Context, traceContextKey any) Logger {
	if traceContextKey == nil {
		traceContextKey = TraceContextKey
	}

	traceID := ""
	if value := ctx.Value(traceContextKey); value != nil {
		traceID = fmt.Sprintf("<%v>", value)
	}

	return MyZapLogger{
		traceID: traceID,
		tag:     l.tag,
		log:     l.log,
	}
}

func (l MyZapLogger) WithTag(tag string) Logger {
	if tag != "" {
		tag = "[" + tag + "]"
	}

	return MyZapLogger{
		traceID: l.traceID,
		tag:     tag,
		log:     l.log,
	}
}

func (l MyZapLogger) Debug(args ...any) {
	prefix := l.getPrefix()

	if prefix == "" {
		l.log.Debug(args...)
	} else {
		l.log.Debug(append([]any{prefix}, args...)...)
	}
}

func (l MyZapLogger) Info(args ...any) {
	prefix := l.getPrefix()

	if prefix == "" {
		l.log.Info(args...)
	} else {
		l.log.Info(append([]any{prefix}, args...)...)
	}
}

func (l MyZapLogger) Warn(args ...any) {
	prefix := l.getPrefix()

	if prefix == "" {
		l.log.Warn(args...)
	} else {
		l.log.Warn(append([]any{prefix}, args...)...)
	}
}

func (l MyZapLogger) Error(args ...any) {
	prefix := l.getPrefix()

	if prefix == "" {
		l.log.Error(args...)
	} else {
		l.log.Error(append([]any{prefix}, args...)...)
	}
}

func (l MyZapLogger) Debugf(format string, args ...any) {
	prefix := l.getPrefix()

	if prefix == "" {
		l.log.Debugf(format, args...)
	} else {
		l.log.Debugf("%s "+format, append([]any{prefix}, args...)...)
	}
}

func (l MyZapLogger) Infof(format string, args ...any) {
	prefix := l.getPrefix()

	if prefix == "" {
		l.log.Infof(format, args...)
	} else {
		l.log.Infof("%s "+format, append([]any{prefix}, args...)...)
	}
}

func (l MyZapLogger) Warnf(format string, args ...any) {
	prefix := l.getPrefix()

	if prefix == "" {
		l.log.Warnf(format, args...)
	} else {
		l.log.Warnf("%s "+format, append([]any{prefix}, args...)...)
	}
}

func (l MyZapLogger) Errorf(format string, args ...any) {
	prefix := l.getPrefix()

	if prefix == "" {
		l.log.Errorf(format, args...)
	} else {
		l.log.Errorf("%s "+format, append([]any{prefix}, args...)...)
	}
}

func (l MyZapLogger) Sync() {
	l.log.Sync()
}

func (l MyZapLogger) getPrefix() string {
	if l.traceID != "" && l.tag != "" {
		return l.traceID + " " + l.tag
	}

	return l.traceID + l.tag
}
