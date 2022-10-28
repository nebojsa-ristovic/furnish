/*Package log contains the improved logger interface and it's nil implementation
Log levels / verbosity go from lowest to highest: Debug, Info, Warn, Error.
Setting the logger verbosity to a value will print that or higher level logs.
*/
package log

import (
	"context"
)

const (
	// DebugVerbosity is used for printing debug logs and higher.
	DebugVerbosity = "DEBUG"
	// InfoVerbosity is used for printing info logs and higher.
	InfoVerbosity = "INFO"
	// WarnVerbosity is used for printing warning logs and higher.
	WarnVerbosity = "WARN"
	// ErrorVerbosity is used for printing error logs only.
	ErrorVerbosity = "ERROR"
)

var global ContextLogger = NewZapLogger(&Config{})

// Logger is an interface which declares logging method of different levels.
// Every logger method accepts a message as the first parameter and key value pairs as varargs.
// If any of the keys do not have a pair they will not be printed.
// Keys must be of type string, if any of the keys are not, that key value pair will be skipped.
type Logger interface {
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})

	Fatal(msg string, fields ...interface{})
}

// ContextLogger is an interface which declares logging with context
// It contains all regular Logger methods
// Context logs will try to extract the current logging scope and group the logs.
type ContextLogger interface {
	Logger

	DebugCtx(ctx context.Context, msg string, fields ...interface{})
	InfoCtx(ctx context.Context, msg string, fields ...interface{})
	WarnCtx(ctx context.Context, msg string, fields ...interface{})
	ErrorCtx(ctx context.Context, msg string, fields ...interface{})

	FatalCtx(ctx context.Context, msg string, fields ...interface{})
}

// NilLogger is a nil logger implementation used for testing.
type NilLogger struct{}

var nilLogger ContextLogger = &NilLogger{}

// NewNilLogger returns an empty nil logger.
func NewNilLogger() Logger {
	return nilLogger
}

// Debug is a empty stub for nil logger. It does nothing.
func (*NilLogger) Debug(_ string, _ ...interface{}) {}

// Info is a empty stub for nil logger. It does nothing.
func (*NilLogger) Info(_ string, _ ...interface{}) {}

// Warn is a empty stub for nil logger. It does nothing.
func (*NilLogger) Warn(_ string, _ ...interface{}) {}

// Error is a empty stub for nil logger. It does nothing.
func (*NilLogger) Error(_ string, _ ...interface{}) {}

// Fatal is a empty stub for nil logger. It does nothing.
func (*NilLogger) Fatal(_ string, _ ...interface{}) {}

// Debug is a empty stub for nil logger. It does nothing.
func (*NilLogger) DebugCtx(_ context.Context, _ string, _ ...interface{}) {}

// InfoCtx is a empty stub for nil logger. It does nothing.
func (*NilLogger) InfoCtx(_ context.Context, _ string, _ ...interface{}) {}

// WarnCtx is a empty stub for nil logger. It does nothing.
func (*NilLogger) WarnCtx(_ context.Context, _ string, _ ...interface{}) {}

// ErrorCtx is a empty stub for nil logger. It does nothing.
func (*NilLogger) ErrorCtx(_ context.Context, _ string, _ ...interface{}) {}

// FatalCtx is a empty stub for nil logger. It does nothing.
func (*NilLogger) FatalCtx(_ context.Context, _ string, _ ...interface{}) {}

// SetGlobal sets the global logger, unsafe to call if logger is already in use.
func SetGlobal(l ContextLogger) {
	global = l
}

// Debug is a wrapper around global's log Debug
func Debug(msg string, fields ...interface{}) { global.Debug(msg, fields...) }

// Info is a wrapper around global's log Info
func Info(msg string, fields ...interface{}) { global.Info(msg, fields...) }

// Warn is a wrapper around global's log Warn
func Warn(msg string, fields ...interface{}) { global.Warn(msg, fields...) }

// Error is a wrapper around global's log Error
func Error(msg string, fields ...interface{}) { global.Error(msg, fields...) }

// Fatal is a wrapper around global's log Fatal
func Fatal(msg string, fields ...interface{}) { global.Fatal(msg, fields...) }

// DebugCtx is a wrapper around global's log DebugCtx
func DebugCtx(ctx context.Context, msg string, fields ...interface{}) {
	global.DebugCtx(ctx, msg, fields...)
}

// InfoCtx is a wrapper around global's log InfoCtx
func InfoCtx(ctx context.Context, msg string, fields ...interface{}) {
	global.InfoCtx(ctx, msg, fields...)
}

// WarnCtx is a wrapper around global's log Warn
func WarnCtx(ctx context.Context, msg string, fields ...interface{}) {
	global.WarnCtx(ctx, msg, fields...)
}

// ErrorCtx is a wrapper around global's log ErrorCtx
func ErrorCtx(ctx context.Context, msg string, fields ...interface{}) {
	global.ErrorCtx(ctx, msg, fields...)
}

// FatalCtx is a wrapper around global's log FatalCtx
func FatalCtx(ctx context.Context, msg string, fields ...interface{}) {
	global.FatalCtx(ctx, msg, fields...)
}
