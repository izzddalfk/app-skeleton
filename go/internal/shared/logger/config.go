package logger

import "github.com/rs/zerolog"

type Config struct {
	LogLevel    string
	ServiceName string
}

func setDefaultConfig() {
	// UNIX time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.TimestampFieldName = "ts"

	// all sampled loggers will stop sampling and issue 100% of their log events
	zerolog.DisableSampling(true)

	// default log level is always to `DEBUG`
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}
