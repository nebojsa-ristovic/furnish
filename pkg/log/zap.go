package log

import (
	"context"
	"os"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	FieldServiceName string = "service"
	FieldLogScopeID  string = "log_scope_id"
)

const LogIDHeader string = "x-log-scope-id"

type scopeKey string
type hubKey string
type idKey string

const (
	CtxScopeKey scopeKey = "logScope"
	CtxHubKey   hubKey   = "logHub"
	CtxIDKey    idKey    = "logID"
)

var (
	_ Logger        = (*zapLogger)(nil)
	_ ContextLogger = (*zapLogger)(nil)
)

// zapLogger is the log.Logger implementation.
// It contains a zap sugared logger.
// It has Logstash and Sentry hooks. These are enabled and configured in the configuration.
// Logstash hooks will be sent on log.Info level or higher, Sentry hooks on errors only.
// Invalid Sentry DSN or Logstash URL will just instantiate a logger without the hook.
type zapLogger struct {
	zlog        *zap.SugaredLogger
	serviceName string
}

// NewZapLogger parses the config and returns a new logger.
// It will return an error only if it fails initializing a zap logger.
// It will return the logger if it fails setting up hooks, but will print the errors.
func NewZapLogger(cfg *Config) *zapLogger {
	zlog := newZapLogger(cfg)
	logger := &zapLogger{}

	if cfg.ServiceName != "" {
		logger.serviceName = cfg.ServiceName
	}

	zap.ReplaceGlobals(zlog)
	logger.zlog = zlog.Sugar()
	return logger
}

// Debug prints debug logs. Only seen if verbosity level log.Debug.
func (log *zapLogger) Debug(msg string, fields ...interface{}) {
	log.zlog.Debugw(msg, log.addServiceName(fields...)...)
}

// Info prints info logs. Only seen if verbosity level log.Info or lower.
func (log *zapLogger) Info(msg string, fields ...interface{}) {
	log.zlog.Infow(msg, log.addServiceName(fields...)...)
}

// Warn prints warning logs. Only seen if verbosity level log.Warn or lower.
func (log *zapLogger) Warn(msg string, fields ...interface{}) {
	log.zlog.Warnw(msg, log.addServiceName(fields...)...)
}

// Error prints error logs. Only seen if verbosity level log.Error or lower.
func (log *zapLogger) Error(msg string, fields ...interface{}) {
	log.zlog.Errorw(msg, log.addServiceName(fields...)...)
}

// Fatal prints a log and shuts down the program. It uses os.Exit(1) to exit. Always visible.
func (log *zapLogger) Fatal(msg string, fields ...interface{}) {
	log.zlog.Fatalw(msg, log.addServiceName(fields...)...)
}

// DebugCtx prints debug logs. Only seen if verbosity level log.Debug.
// It will try to extract the scope from context before logging.
func (log *zapLogger) DebugCtx(ctx context.Context, msg string, fields ...interface{}) {
	log.zlog.Debugw(msg, log.enrichLog(ctx, fields...)...)
}

// Info prints info logs. Only seen if verbosity level log.Info or lower.
// It will try to extract the scope from context before logging.
func (log *zapLogger) InfoCtx(ctx context.Context, msg string, fields ...interface{}) {
	log.zlog.Infow(msg, log.enrichLog(ctx, fields...)...)
}

// Warn prints warning logs. Only seen if verbosity level log.Warn or lower.
// It will try to extract the scope from context before logging.
func (log *zapLogger) WarnCtx(ctx context.Context, msg string, fields ...interface{}) {
	log.zlog.Warnw(msg, log.enrichLog(ctx, fields...)...)
}

// Error prints error logs. Only seen if verbosity level log.Error or lower.
// It will try to extract the scope from context before logging.
func (log *zapLogger) ErrorCtx(ctx context.Context, msg string, fields ...interface{}) {
	log.zlog.Errorw(msg, log.enrichLog(ctx, fields...)...)
}

// Fatal prints a log and shuts down the program. It uses os.Exit(1) to exit. Always visible.
// It will try to extract the scope from context before logging.
func (log *zapLogger) FatalCtx(ctx context.Context, msg string, fields ...interface{}) {
	log.zlog.Fatalw(msg, log.enrichLog(ctx, fields...)...)
}

// revive:disable-next-line:confusing-naming this is called in the exported function
// it's unexported counterpart is used as seen in stdlib.
func newZapLogger(cfg *Config) *zap.Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	jsonEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	atom := zap.NewAtomicLevel()
	atom.SetLevel(mapVerbosityLevel(cfg.LogVerbosity))
	core := zapcore.NewCore(jsonEncoder, zapcore.Lock(os.Stdout), atom)

	logger := zap.New(core).WithOptions(
		zap.AddCallerSkip(2),
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.PanicLevel),
	)

	defer logger.Sync()
	return logger
}

func (log *zapLogger) enrichLog(ctx context.Context, fields ...interface{}) []interface{} {
	fields = log.addLogID(ctx, fields...)
	fields = log.addServiceName(fields...)
	return fields
}

func (log *zapLogger) addServiceName(fields ...interface{}) []interface{} {
	if log.serviceName == "" {
		return fields
	}
	return append(fields, zap.String(FieldServiceName, log.serviceName))
}

func (log *zapLogger) addLogID(ctx context.Context, fields ...interface{}) []interface{} {
	if val, ok := ctx.Value(CtxIDKey).(string); ok {
		return append(fields, zap.String(FieldLogScopeID, val))
	}
	return fields
}

func mapVerbosityLevel(verbosity string) zapcore.Level {
	switch verbosity {
	case DebugVerbosity:
		return zapcore.DebugLevel
	case InfoVerbosity:
		return zapcore.InfoLevel
	case WarnVerbosity:
		return zapcore.WarnLevel
	case ErrorVerbosity:
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// ContextWithScope will return the context with a new scope added.
// LogCtx functions will try to get the scope out of the context before logging so we can group logs
// by their scope.
func ContextWithScope(ctx context.Context) context.Context {
	if val := ctx.Value(CtxIDKey); val == nil {
		ctx = context.WithValue(ctx, CtxIDKey, uuid.New().String())
	}
	return ctx
}
