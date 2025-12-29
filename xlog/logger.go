package xlog

import (
	"context"
	"log/slog"
	"os"
	"sync/atomic"
)

var defaultLogger atomic.Pointer[Logger]

func init() {
	defaultLogger.Store(NewText(LevelInfo))
}

func Debug(msg string, fields ...slog.Attr) {
	Defualt().Debug(context.Background(), msg, fields...)
}

func Info(msg string, fields ...slog.Attr) {
	Defualt().Info(context.Background(), msg, fields...)
}

func Warn(msg string, fields ...slog.Attr) {
	Defualt().Warn(context.Background(), msg, fields...)
}
func Error(msg string, fields ...slog.Attr) {
	Defualt().Error(context.Background(), msg, fields...)
}

type Logger struct {
	json bool
	s    *slog.Logger
}

const (
	LevelDebug slog.Level = slog.LevelDebug
	LevelInfo  slog.Level = slog.LevelInfo
	LevelWarn  slog.Level = slog.LevelWarn
	LevelError slog.Level = slog.LevelError
)

var (
	Int  = slog.Int
	I64  = slog.Int64
	U64  = slog.Uint64
	F64  = slog.Float64
	Str  = slog.String
	Dur  = slog.Duration
	Any  = slog.Any
	Bool = slog.Bool
	Time = slog.Time
)

func I32(key string, i int32) slog.Attr {
	return slog.Int(key, int(i))
}
func U32(key string, i uint32) slog.Attr {
	return slog.Uint64(key, uint64(i))
}
func F32(key string, f float32) slog.Attr {
	return slog.Float64(key, float64(f))
}
func Err(e error) slog.Attr {
	return slog.Any("error", e)
}
func Uid(id string) slog.Attr {
	return slog.String("userId", id)
}
func With(args ...any) *Logger {
	return Defualt().With(args...)
}
func WithLevel(level slog.Level) *Logger {
	return Defualt().WithLevel(level)
}
func ParseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
func NewText(level slog.Level) *Logger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})
	return &Logger{s: slog.New(newRequestContextHandler(handler)), json: false}
}
func NewJSON(level slog.Level) *Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})
	return &Logger{s: slog.New(newRequestContextHandler(handler)), json: true}
}

func Defualt() *Logger {
	return defaultLogger.Load()
}
func SetDefault(l *Logger) {
	defaultLogger.Store(l)
}
func (l *Logger) With(args ...any) *Logger {
	return &Logger{s: l.s.With(args...)}
}
func (l *Logger) WithLevel(level slog.Level) *Logger {
	if l.json {
		return NewJSON(level)
	}
	return NewText(level)
}
func (l *Logger) Debug(ctx context.Context, msg string, fields ...slog.Attr) {
	l.s.LogAttrs(ctx, slog.LevelDebug, msg, fields...)
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...slog.Attr) {
	l.s.LogAttrs(ctx, slog.LevelInfo, msg, fields...)
}

func (l *Logger) Warn(ctx context.Context, msg string, fields ...slog.Attr) {
	l.s.LogAttrs(ctx, slog.LevelWarn, msg, fields...)
}
func (l *Logger) Error(ctx context.Context, msg string, fields ...slog.Attr) {
	l.s.LogAttrs(ctx, slog.LevelError, msg, fields...)
}

type remoteIpKey struct{}
type requestIdKey struct{}

type requestContextHandler struct {
	handler slog.Handler
}

var (
	KeyRequestID              = requestIdKey{}
	KeyRemoteIP               = remoteIpKey{}
	_            slog.Handler = requestContextHandler{}
)

func newRequestContextHandler(handler slog.Handler) slog.Handler {
	return requestContextHandler{
		handler: handler,
	}
}

func (h requestContextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h requestContextHandler) Handle(ctx context.Context, record slog.Record) error {
	if v, ok := ctx.Value(KeyRemoteIP).(string); ok {
		record.AddAttrs(slog.String("remote_ip", v))
	}
	if v, ok := ctx.Value(KeyRequestID).(string); ok {
		record.AddAttrs(slog.String("request_id", v))
	}
	return h.handler.Handle(ctx, record)
}

func (h requestContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return requestContextHandler{h.handler.WithAttrs(attrs)}
}

func (h requestContextHandler) WithGroup(name string) slog.Handler {
	return h.handler.WithGroup(name)
}
