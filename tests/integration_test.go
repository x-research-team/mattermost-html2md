package tests

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/x-research-team/mattermost-html2md/cmd/server"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const wait = 5 * time.Second

type request struct {
	Text string `json:"text"`
}

func TestMain(t *testing.T) {
	go func() {
		err := server.Run(context.Background(), zerolog.New(os.Stderr).With().Timestamp().Caller().Logger())
		require.NoError(t, err)
	}()

	time.Sleep(wait)

	client := resty.New()

	t.Run("send success", func(t *testing.T) {
		resp, err := client.R().
			SetHeader("X-API-KEY", "test").
			SetBody(request{Text: "<h1>Hello World</h1><p>This is a simple HTML document.</p>"}).
			Post("http://localhost:8080/send")
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode())
	})

	t.Run("send auth error", func(t *testing.T) {
		resp, err := client.R().
			SetBody(request{Text: "<h1>Hello World</h1><p>This is a simple HTML document.</p>"}).
			Post("http://localhost:8080/send")
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode())
	})

	t.Run("send invalid request", func(t *testing.T) {
		resp, err := client.R().
			SetHeader("X-API-KEY", "test").
			Post("http://localhost:8080/send")
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode())
	})
}
