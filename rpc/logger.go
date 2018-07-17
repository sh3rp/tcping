package rpc

import (
	"os"

	"github.com/rs/zerolog/log"

	"github.com/rs/zerolog"
)

// default instance, we can declare others later if necessary

var LOGGER = newLogger()

// Logger is a basic logging interface to be used for logging throughout the rpc server
type Logger interface {
	Debug(string, ...interface{})
	Info(string, ...interface{})
	Error(string, ...interface{})
}

// zeroLogger is an implementation of the Logger interface that uses ZeroLog as the underlying
// implementation
type zeroLogger struct {
	logger zerolog.Logger
}

// newLogger creates a new instance of zero logger with pretty printing enabled
func newLogger() Logger {
	return zeroLogger{
		logger: log.Output(zerolog.ConsoleWriter{Out: os.Stderr}),
	}
}

func (l zeroLogger) Debug(fmt string, vars ...interface{}) {
	l.logger.Debug().Msgf(fmt, vars)
}

func (l zeroLogger) Info(fmt string, vars ...interface{}) {
	l.logger.Info().Msgf(fmt, vars)
}

func (l zeroLogger) Error(fmt string, vars ...interface{}) {
	l.logger.Error().Msgf(fmt, vars)
}
