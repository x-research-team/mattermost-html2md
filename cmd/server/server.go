package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/x-research-team/mattermost-html2md/internal/api"
	"github.com/x-research-team/mattermost-html2md/internal/config"
	"github.com/x-research-team/mattermost-html2md/internal/flags"
	"github.com/x-research-team/mattermost-html2md/internal/services/mattermost"
	"github.com/x-research-team/mattermost-html2md/pkg/log"
	"github.com/x-research-team/mattermost-html2md/pkg/slice"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/JohannesKaufmann/html-to-markdown/plugin"
	rest "github.com/go-micro/plugins/v4/server/http"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/server"
	"golang.org/x/sync/errgroup"
)

func Run(ctx context.Context, l zerolog.Logger) error {
	ctx, stop := signal.NotifyContext(ctx, os.Kill, os.Interrupt)
	defer stop()

	group, ctx := errgroup.WithContext(ctx)

	cfg, err := config.Load(&l, ".env")
	if err != nil {
		return fmt.Errorf("load: %w", err)
	}

	logger.DefaultLogger = log.NewAdapter(l)

	converter := md.NewConverter("", true, nil)
	converter.Use(plugin.GitHubFlavored())
	client := model.NewAPIv4Client(cfg.Mattermost.URL)
	client.SetToken(cfg.Mattermost.Token)

	service := mattermost.New(cfg, converter, client)

	srv := rest.NewServer(
		server.Address(net.JoinHostPort(cfg.Server.Host, strconv.Itoa(cfg.Server.Port))),
		server.Context(ctx),
		server.Wait(nil),
	)

	handler := srv.NewHandler(api.New(cfg, &l, &http.Server{
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}, service).Router())

	if err := srv.Handle(handler); err != nil {
		return fmt.Errorf("handle: %w", err)
	}

	svc := micro.NewService(
		micro.Name(cfg.Name),
		micro.Version(cfg.Version),
		micro.Server(srv),
		micro.Context(ctx),
		micro.Registry(registry.NewRegistry()),
		micro.RegisterTTL(cfg.Heartbeat.TTL),
		micro.RegisterInterval(cfg.Heartbeat.Interval),
		micro.Flags(slice.Merge[cli.Flag](flags.Test)...),
	)

	svc.Init()

	group.Go(svc.Run)
	group.Go(shutdown(ctx, l, svc))

	if err := group.Wait(); err != nil {
		return fmt.Errorf("wait: %w", err)
	}

	return nil
}

func shutdown(ctx context.Context, l zerolog.Logger, svc micro.Service) func() error {
	return func() error {
		<-ctx.Done()

		l.Info().Msg("shutting down")
		return nil
	}
}
