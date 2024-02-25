package log

import (
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
)

func InitHooks(client *resty.Client) {
	client.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		// Log the request details
		log.Info().
			Str("method", req.Method).
			Str("url", req.URL).
			Any("body", req.Body).
			Msg("request sent")
		return nil // return nil to let the request proceed
	})

	client.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		if resp.StatusCode() >= 300 {
			log.Error().
				Str("url", resp.Request.URL).
				Str("status", resp.Status()).
				Str("body", resp.String()).
				Msg("response error")
		} else {
			// Log the response details
			log.Info().
				Int("status", resp.StatusCode()).
				Dur("latency", resp.Time()).
				Str("body", resp.String()).
				Msg("response received")
		}
		return nil // return nil to let the execution continue
	})
}
