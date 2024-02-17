package logger

import "context"

var CtxKey = "logger_ctx"

type Logger interface {
	SetFields(fields ...Field) Logger
	SetError(err error) Logger

	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
}

type Field struct {
	Key   string
	Value interface{}
}

// GetFromContext used to get active `logger` instance from the provided context
func GetFromContext(ctx context.Context) Logger {
	return ctx.Value(CtxKey).(Logger)
}
