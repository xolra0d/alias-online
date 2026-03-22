package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

type logMessage struct {
	level  slog.Level
	prefix string
	msg    string
	args   []any
}

// BaseLogger helps to create async loggers, using WithPrefix.
type BaseLogger struct {
	logger *slog.Logger
	q      chan *logMessage
}

// NewBaseLogger creates new empty logger.
func NewBaseLogger(queueLen int) *BaseLogger {
	return &BaseLogger{
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
		q:      make(chan *logMessage, queueLen),
	}
}

// StartLogging starts to retrieve messages received from all PrefixLogger's.
func (l *BaseLogger) StartLogging() {
	for msg := range l.q {
		t := fmt.Sprintf("[%s] %s", msg.prefix, msg.msg)
		l.logger.Log(context.Background(), msg.level, t, msg.args...)
	}
}

// EndLogging closes the channel for messages.
func (l *BaseLogger) EndLogging() {
	close(l.q)
}

// PrefixLogger asyncronously logs logs.
type PrefixLogger struct {
	Prefix string
	q      chan *logMessage // points to BaseLogger.q
}

// WithPrefix creates new async logger.
func (l *BaseLogger) WithPrefix(prefix string) *PrefixLogger {
	return &PrefixLogger{
		Prefix: prefix,
		q:      l.q,
	}
}

// CopyWithPrefix is the same as WithPrefix
func (l *PrefixLogger) CopyWithPrefix(prefix string) *PrefixLogger {
	return &PrefixLogger{
		Prefix: prefix,
		q:      l.q,
	}
}

// Info logs data with INFO level
func (l *PrefixLogger) Info(msg string, args ...any) {
	m := &logMessage{
		level:  slog.LevelInfo,
		prefix: l.Prefix,
		msg:    msg,
		args:   args,
	}
	l.q <- m
}

// Warn logs data with WARN level
func (l *PrefixLogger) Warn(msg string, args ...any) {
	m := &logMessage{
		level:  slog.LevelWarn,
		prefix: l.Prefix,
		msg:    msg,
		args:   args,
	}
	l.q <- m
}

// Error logs data with ERROR level
func (l *PrefixLogger) Error(msg string, args ...any) {
	m := &logMessage{
		level:  slog.LevelError,
		prefix: l.Prefix,
		msg:    msg,
		args:   args,
	}
	l.q <- m
}

// Debug logs data with DEBUG level
func (l *PrefixLogger) Debug(msg string, args ...any) {
	m := &logMessage{
		level:  slog.LevelDebug,
		prefix: l.Prefix,
		msg:    msg,
		args:   args,
	}
	l.q <- m
}
