package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

var (
	l *Log
)

type Option struct {
	AppName      string
	Level        string
	ReportCaller bool
}

type Log struct {
	opt Option
}

func NewLogger(opt Option) *Log {
	if l == nil {
		l = &Log{opt: opt}
	}
	return l
}

func (l *Log) GetLog() *slog.Logger {
	level := l.parseLogLevel()
	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: l.opt.ReportCaller,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)

	logger := slog.New(handler)

	logger = logger.With(slog.String("application", l.opt.AppName))

	return logger
}

func (l *Log) parseLogLevel() slog.Level {
	switch l.opt.Level {
	case "INFO":
		return slog.LevelInfo
	case "ERROR":
		return slog.LevelError
	case "DEBUG":
		return slog.LevelDebug
	case "WARN":
		return slog.LevelWarn
	default:
		return slog.LevelInfo
	}
}

type log struct {
	ctx context.Context
}

func (l *log) Info(msg string, args ...any) {
	slog.InfoContext(l.ctx, msg, args...)
}

func (l *log) Infof(format string, args ...any) {
	slog.InfoContext(l.ctx, fmt.Sprintf(format, args...))
}

func (l *log) Warn(msg string, args ...any) {
	slog.WarnContext(l.ctx, msg, args...)
}

func (l *log) Warnf(format string, args ...any) {
	slog.WarnContext(l.ctx, fmt.Sprintf(format, args...))
}

func (l *log) Debug(msg string, args ...any) {
	slog.DebugContext(l.ctx, msg, args...)
}

func (l *log) Debugf(format string, args ...any) {
	slog.DebugContext(l.ctx, fmt.Sprintf(format, args...))
}

func (l *log) Error(msg string, args ...any) {
	slog.ErrorContext(l.ctx, msg, args...)
}

func (l *log) Errorf(format string, args ...any) {
	slog.ErrorContext(l.ctx, fmt.Sprintf(format, args...))
}

func FromContext(ctx context.Context) Logger {
	logger := &log{
		ctx,
	}
	return logger
}
