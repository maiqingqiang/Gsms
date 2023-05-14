package gsms

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// Silent silent log level
	Silent LogLevel = iota + 1
	// Error error log level
	Error
	// Warn warn log level
	Warn
	// Info info log level
	Info
)

type logger struct {
	level  LogLevel
	writer *zap.Logger
}

const callerSkipOffset = 2

func NewLogger(opts ...zap.Option) Logger {
	opts = append(opts, zap.AddCallerSkip(callerSkipOffset))

	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Encoding:    "console",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	writer, err := config.Build(opts...)
	if err != nil {
		panic(err)
	}

	return &logger{
		level:  Warn,
		writer: writer,
	}
}

func (l *logger) LogMode(level LogLevel) Logger {
	newLogger := *l
	newLogger.level = level
	return &newLogger
}

func (l *logger) Info(v ...interface{}) {
	if l.level >= Info {
		l.writer.Info(fmt.Sprint(v...))
	}
}

func (l *logger) Infof(format string, v ...interface{}) {
	l.Info(fmt.Sprintf(format, v...))
}

func (l *logger) Warn(v ...interface{}) {
	if l.level >= Warn {
		l.writer.Warn(fmt.Sprint(v...))
	}
}

func (l *logger) Warnf(format string, v ...interface{}) {
	l.Warn(fmt.Sprintf(format, v...))
}

func (l *logger) Error(v ...interface{}) {
	if l.level >= Error {
		l.writer.Error(fmt.Sprint(v...))
	}
}

func (l *logger) Errorf(format string, v ...interface{}) {
	l.Error(fmt.Sprintf(format, v...))
}
