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
		URL     string        `env:"MATTERMOST_URL"`
		Timeout time.Duration `env:"MATTERMOST_TIMEOUT" envDefault:"10s"`
		Debug   bool          `env:"MATTERMOST_DEBUG" envDefault:"false"`
		Token   string        `env:"MATTERMOST_TOKEN"`
		Webhook string        `env:"MATTERMOST_WEBHOOK_URL" required:"true"`
		Channel string        `env:"MATTERMOST_CHANNEL" required:"true"`
	}

	IMAP struct {
		Host string `env:"IMAP_HOST" required:"true"`
		Port int    `env:"IMAP_PORT" envDefault:"993"`
		User string `env:"IMAP_USER" required:"true"`
		Pass string `env:"IMAP_PASS" required:"true"`
	}

	Cron struct {
		Interval string `env:"CRON_INTERVAL" envDefault:"* * * * *"`
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
