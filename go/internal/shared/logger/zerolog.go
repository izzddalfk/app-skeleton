package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

type zerologLogger struct {
	instance *zerolog.Logger
}

func NewZerologLogger(config Config) (Logger, error) {
	setDefaultConfig()

	// override log level if passed in config
	if config.LogLevel != "" {
		var level zerolog.Level = zerolog.DebugLevel
		switch strings.ToLower(config.LogLevel) {
		case "info":
			level = zerolog.InfoLevel
		case "warn", "warning":
			level = zerolog.WarnLevel
		case "error":
			level = zerolog.ErrorLevel
		}
		zerolog.SetGlobalLevel(level)
	}
	// validate service name value
	if config.ServiceName == "" {
		return nil, fmt.Errorf("service name is empty")
	}

	logger := zerolog.New(os.Stderr).
		With().
		Str("service_name", config.ServiceName).
		Logger()

	return &zerologLogger{
		instance: &logger,
	}, nil
}

func (l *zerologLogger) SetFields(fields ...Field) Logger {
	l.instance.With().Fields(fields)

	return l
}

func (l *zerologLogger) SetError(err error) Logger {
	l.instance.With().Err(err)

	return l
}

func (l *zerologLogger) Debug(msg string, fields ...Field) {
	l.instance.Debug().Fields(l.toZerologFields(fields...)).Msg(msg)
}

func (l *zerologLogger) Info(msg string, fields ...Field) {
	l.instance.Info().Fields(l.toZerologFields(fields...)).Msg(msg)
}

func (l *zerologLogger) Warn(msg string, fields ...Field) {
	l.instance.Warn().Fields(l.toZerologFields(fields...)).Msg(msg)
}

func (l *zerologLogger) Error(msg string, fields ...Field) {
	l.instance.Error().Fields(l.toZerologFields(fields...)).Msg(msg)
}

func (l zerologLogger) toZerologFields(fields ...Field) map[string]interface{} {
	out := make(map[string]interface{}, len(fields))
	for _, f := range fields {
		out[f.Key] = f.Value
	}
	return out
}
