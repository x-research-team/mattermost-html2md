package middleware

import (
	"context"
	"time"

	"github.com/x-research-team/mattermost-html2md/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func Metrics(cfg *config.Config, l *zerolog.Logger) gin.HandlerFunc {
	meter := otel.GetMeterProvider().Meter(cfg.Name)

	// Create a counter instrument
	success, err := meter.Int64Counter(
		"success_total",
		metric.WithDescription("The total number of processed success"),
	)
	if err != nil {
		l.Panic().Err(err).Msg("create success counter")
	}

	errors, err := meter.Int64Counter(
		"errors_total",
		metric.WithDescription("The total number of processed errors"),
	)
	if err != nil {
		l.Panic().Err(err).Msg("create errors counter")
	}

	// Create a histogram instrument
	latency, err := meter.Float64Histogram(
		"request_latency",
		metric.WithDescription("The latency of requests"), metric.WithUnit("seconds"))
	if err != nil {
		l.Panic().Err(err).Msg("create latency histogram")
	}

	return func(c *gin.Context) {
		now := time.Now()
		c.Next()

		// Record the latency
		latency.Record(context.Background(), time.Since(now).Seconds(), metric.WithAttributes(
			attribute.String("method", c.Request.Method),
			attribute.String("path", c.Request.URL.Path),
			attribute.Int("status", c.Writer.Status()),
		))

		if c.Writer.Status() >= 300 {
			success.Add(context.Background(), 1, metric.WithAttributes(
				attribute.String("method", c.Request.Method),
				attribute.String("path", c.Request.URL.Path),
				attribute.Int("status", c.Writer.Status()),
			))
		} else {
			errors.Add(context.Background(), 1, metric.WithAttributes(
				attribute.String("method", c.Request.Method),
				attribute.String("path", c.Request.URL.Path),
				attribute.Int("status", c.Writer.Status()),
			))
		}
	}
}
