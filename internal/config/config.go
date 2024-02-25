package config

import (
	"time"

	"github.com/lindenlab/env"
	"github.com/rs/zerolog"
)

type Config struct {
	Name        string `env:"APP_NAME" required:"true"`
	Description string `env:"APP_DESCRIPTION" required:"true"`
	Version     string `env:"APP_VERSION" required:"true"`

	Server struct {
		Host string `env:"SERVER_HOST" envDefault:"localhost"`
		Port int    `env:"SERVER_PORT" envDefault:"8080"`

		ReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT" envDefault:"10s"`
		WriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT" envDefault:"10s"`
		IdleTimeout  time.Duration `env:"SERVER_IDLE_TIMEOUT" envDefault:"120s"`
	}

	Heartbeat struct {
		TTL      time.Duration `env:"HEARTBEAT_TTL" envDefault:"10s"`
		Interval time.Duration `env:"HEARTBEAT_INTERVAL" envDefault:"5s"`
	}

	Mattermost struct {
		Webhook string        `env:"MATTERMOST_WEBHOOK" required:"true"`
		User    string        `env:"MATTERMOST_USER" required:"true"`
		Channel string        `env:"MATTERMOST_CHANNEL" required:"true"`
		Timeout time.Duration `env:"MATTERMOST_TIMEOUT" envDefault:"10s"`
		Debug   bool          `env:"MATTERMOST_DEBUG" envDefault:"false"`
	}
}

func Load(l *zerolog.Logger, paths ...string) (*Config, error) {
	cfg := &Config{}

	for _, path := range paths {
		if err := env.Load(path); err != nil {
			l.Error().Err(err).Msg("load")
		}
	}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
