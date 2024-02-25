package middleware

import (
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func Logger(log *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		var err error

		response := make([]byte, 0)

		if c.Request.Response != nil {
			response, err = io.ReadAll(c.Request.Response.Body)
			if err != nil {
				log.Error().Err(err).Msg("read response body")
			}
		}

		// Log the request
		if c.Writer.Status() >= 300 {
			log.Error().
				Str("host", c.Request.Host).
				Str("ip", c.ClientIP()).
				Str("user-agent", c.Request.UserAgent()).
				Str("referer", c.Request.Referer()).
				Str("method", c.Request.Method).
				Str("path", c.Request.URL.Path).
				Int("status", c.Writer.Status()).
				Dur("latency", time.Since(start)).
				Msg(string(response))
				return
		}

		log.Info().
			Str("host", c.Request.Host).
			Str("ip", c.ClientIP()).
			Str("user-agent", c.Request.UserAgent()).
			Str("referer", c.Request.Referer()).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", c.Writer.Status()).
			Dur("latency", time.Since(start)).
			Msg("request handled")
	}
}
