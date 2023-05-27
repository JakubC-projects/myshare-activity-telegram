package log

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var L *zerolog.Logger

func init() {
	L = configureGlobalLogger()
}

func configureGlobalLogger() *zerolog.Logger {
	zerolog.LevelFieldName = "severity"
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = time.RFC3339Nano

	logger := zerolog.New(os.Stderr).
		With().
		Timestamp().
		Caller().
		Logger()

	serviceName := os.Getenv("K_SERVICE")
	if serviceName != "" {
		logger = logger.
			Level(zerolog.DebugLevel)
	} else {
		logger = logger.Output(
			zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"},
		)
	}

	return &logger
}
