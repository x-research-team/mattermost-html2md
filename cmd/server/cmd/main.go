package main

import (
	"context"
	"os"

	"github.com/x-research-team/mattermost-html2md/cmd/server"

	"github.com/rs/zerolog"
)

func main() {
	l := zerolog.New(os.Stderr).With().Timestamp().Caller().Logger()
	if err := server.Run(context.Background(), l); err != nil {
		l.Fatal().Err(err).Msg("startup")
	}
}
