package log

import (
	"github.com/rs/zerolog"
	"go-micro.dev/v4/logger"
)

type Adapter struct {
	log zerolog.Logger
}

func NewAdapter(l zerolog.Logger) *Adapter {
	return &Adapter{
		log: l,
	}
}

func (z *Adapter) Init(opts ...logger.Option) error {
	// Initialize your logger with options if necessary
	return nil
}

func (z *Adapter) Options() logger.Options {
	// Add options to your logger
	return logger.Options{}
}

func (z *Adapter) Fields(fields map[string]interface{}) logger.Logger {
	z.log = z.log.With().Fields(fields).Logger()
	return z
}

func (z *Adapter) Log(level logger.Level, args ...interface{}) {
	// Log messages at the specified level
	switch level {
	case logger.InfoLevel:
		z.log.Info().Msgf("%v", args)
	case logger.ErrorLevel:
		z.log.Error().Msgf("%v", args)
	case logger.WarnLevel:
		z.log.Warn().Msgf("%v", args)
	case logger.DebugLevel:
		z.log.Debug().Msgf("%v", args)
	case logger.TraceLevel:
		z.log.Trace().Msgf("%v", args)
	case logger.FatalLevel:
		z.log.Fatal().Msgf("%v", args)
	}
}

func (z *Adapter) Logf(level logger.Level, format string, args ...interface{}) {
	// Log formatted messages at the specified level
	switch level {
	case logger.InfoLevel:
		z.log.Info().Msgf(format, args...)
	case logger.ErrorLevel:
		z.log.Error().Msgf(format, args...)
	case logger.WarnLevel:
		z.log.Warn().Msgf(format, args...)
	case logger.DebugLevel:
		z.log.Debug().Msgf(format, args...)
	case logger.TraceLevel:
		z.log.Trace().Msgf(format, args...)
	case logger.FatalLevel:
		z.log.Fatal().Msgf(format, args...)
	}
}

func (z *Adapter) String() string {
	return "zerolog"
}
