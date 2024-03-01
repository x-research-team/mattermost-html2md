package log

import (
	"github.com/rs/zerolog"
	"go-micro.dev/v4/logger"
)

type MicroAdapter struct {
	log zerolog.Logger
}

func NewMicroAdapter(l zerolog.Logger) *MicroAdapter {
	return &MicroAdapter{
		log: l,
	}
}

func (z *MicroAdapter) Init(opts ...logger.Option) error {
	// Initialize your logger with options if necessary
	return nil
}

func (z *MicroAdapter) Options() logger.Options {
	// Add options to your logger
	return logger.Options{}
}

func (z *MicroAdapter) Fields(fields map[string]interface{}) logger.Logger {
	z.log = z.log.With().Fields(fields).Logger()
	return z
}

func (z *MicroAdapter) Log(level logger.Level, args ...interface{}) {
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

func (z *MicroAdapter) Logf(level logger.Level, format string, args ...interface{}) {
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

func (z *MicroAdapter) String() string {
	return "zerolog"
}

type CronAdapter struct {
	log zerolog.Logger
}

func (c *CronAdapter) Error(err error, msg string, keysAndValues ...interface {
}) {
	c.log.Error().Err(err).Msgf(msg, keysAndValues...)
}

func (c *CronAdapter) Info(msg string, keysAndValues ...interface {
}) {
	c.log.Info().Msgf(msg, keysAndValues...)
}

func NewCronAdapter(l zerolog.Logger) *CronAdapter {
	return &CronAdapter{
		log: l,
	}
}
