package tests

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/x-research-team/mattermost-html2md/cmd/server"
	"github.com/x-research-team/mattermost-html2md/pkg/log"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const wait = 5 * time.Second

type request struct {
	Text    string `json:"text"`
	Channel string `json:"channel"`
}

func TestMain(t *testing.T) {
	go func() {
		err := server.Run(context.Background(), zerolog.New(os.Stderr).With().Timestamp().Caller().Logger())
		require.NoError(t, err)
	}()

	time.Sleep(wait)

	client := resty.New().SetDebug(true)
	log.SetHook(client)

	t.Run("send success", func(t *testing.T) {
		resp, err := client.R().
			SetHeader("X-API-KEY", "test").
			SetBody(request{
				Text:    "<h1>Hello World</h1><p>This is a simple HTML document.</p>",
				Channel: "pyx1obq8e7ympkm4eitq3uq89c", // Set your channel ID here.
			}).
			Post("http://localhost:8080/api/v1/webhook")
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode())
	})

	t.Run("send auth error", func(t *testing.T) {
		resp, err := client.R().
			SetBody(request{
				Text:    "<h1>Hello World</h1><p>This is a simple HTML document.</p>",
				Channel: "pyx1obq8e7ympkm4eitq3uq89c", // Set your channel ID here.
			}).
			Post("http://localhost:8080/api/v1/webhook")
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode())
	})

	t.Run("send invalid request", func(t *testing.T) {
		resp, err := client.R().
			SetHeader("X-API-KEY", "test").
			Post("http://localhost:8080/api/v1/webhook")
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode())
	})
}
