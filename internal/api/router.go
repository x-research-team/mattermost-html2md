package api

import (
	"context"
	"net/http"

	"github.com/x-research-team/mattermost-html2md/internal/api/middleware"
	"github.com/x-research-team/mattermost-html2md/internal/config"
	"github.com/x-research-team/mattermost-html2md/pkg/models/response"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/x-research-team/swagin"
	"github.com/x-research-team/swagin/router"
	"github.com/x-research-team/swagin/security"
	"github.com/x-research-team/swagin/swagger"
)

type api struct {
	*swagin.SwaGin

	cfg     *config.Config
	service MettermostService
}

type MettermostService interface {
	SendAPI(ctx context.Context, text, channel string) error
	SendWebhook(ctx context.Context, text, channel string) error
}

type Router interface {
	Router() *gin.Engine
}

func New(cfg *config.Config, logger *zerolog.Logger, server *http.Server, service MettermostService) Router {
	return &api{
		SwaGin: swagin.New(
			swagger.New(
				cfg.Name,
				cfg.Description,
				cfg.Version,
				swagger.Contact(&openapi3.Contact{
					Name:  "Adel Urazov",
					URL:   "https://github.com/x-research-team/mattermost-html2md",
					Email: "adel.i.urazov@ya.ru",
				}),
				swagger.TermsOfService("https://github.com/x-research-team/mattermost-html2md"),
			),
			swagin.Server(server),
		).Middlewares(
			middleware.Logger(logger),
			gin.Recovery(),
			middleware.Tracer(cfg),
			middleware.Metrics(cfg, logger),
		),

		cfg:     cfg,
		service: service,
	}
}

func (a *api) Router() *gin.Engine {
	a.POST("/api/v1/webhook", router.New(a.Webhook,
		router.Summary("Send HTML to Markdown"),
		router.Description("Send HTML to Markdown"),
		router.ContentType("application/json", router.ContentTypeRequest),
		router.ContentType("application/json", router.ContentTypeResponse),
		router.Responses(router.Response{
			"204": router.ResponseItem{
				Description: "OK",
				Model:       response.Empty{},
			},
			"401": router.ResponseItem{
				Description: "Unauthorized",
				Model:       response.Err{},
			},
			"500": router.ResponseItem{
				Description: "Internal Server Error",
				Model:       response.Err{},
			},
		}),
		router.Security(&security.ApiKey{Name: "X-API-KEY"}),
	))

	a.Init()

	return a.SwaGin.Engine
}
