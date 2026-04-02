package config

import (
	"context"
	"fmt"
	"io"
	"log/slog"
)

type logMessage struct {
	level  slog.Level
	prefix string
	msg    string
	args   []any
}

// Logger helps to create async loggers, using WithPrefix.
type Logger struct {
	logger *slog.Logger
	q      chan *logMessage
}

// NewLogger creates new empty logger.
func NewLogger(queueLen int, out io.Writer) *Logger {
	return &Logger{
		logger: slog.New(slog.NewTextHandler(out, nil)),
		q:      make(chan *logMessage, queueLen),
	}
}

// StartLogging starts to retrieve messages received from all PrefixLogger's.
func (l *Logger) StartLogging() {
	for msg := range l.q {
		t := fmt.Sprintf("[%s] %s", msg.prefix, msg.msg)
		l.logger.Log(context.Background(), msg.level, t, msg.args...)
	}
}

// EndLogging closes the channel for messages.
func (l *Logger) EndLogging() {
	close(l.q)
	for {
		if len(l.q) == 0 {
			break
		}
	}
}

// Info logs data with INFO level
func (l *Logger) Info(prefix string, msg string, args ...any) {
	m := &logMessage{
		level:  slog.LevelInfo,
		prefix: prefix,
		msg:    msg,
		args:   args,
	}
	l.q <- m
}

// Warn logs data with WARN level
func (l *Logger) Warn(prefix string, msg string, args ...any) {
	m := &logMessage{
		level:  slog.LevelWarn,
		prefix: prefix,
		msg:    msg,
		args:   args,
	}
	l.q <- m
}

// Error logs data with ERROR level
func (l *Logger) Error(prefix string, msg string, args ...any) {
	m := &logMessage{
		level:  slog.LevelError,
		prefix: prefix,
		msg:    msg,
		args:   args,
	}
	l.q <- m
}

// Debug logs data with DEBUG level
func (l *Logger) Debug(prefix string, msg string, args ...any) {
	m := &logMessage{
		level:  slog.LevelDebug,
		prefix: prefix,
		msg:    msg,
		args:   args,
	}
	l.q <- m
}
