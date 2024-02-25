package middleware

import (
	"github.com/x-research-team/mattermost-html2md/internal/config"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func Tracer(cfg *config.Config) gin.HandlerFunc {
	tracer := otel.Tracer(cfg.Name)
	return func(c *gin.Context) {
		ctx, span := tracer.Start(c.Request.Context(), c.Request.URL.Path,
			trace.WithAttributes(attribute.String("method", c.Request.Method)))
		defer span.End()

		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
