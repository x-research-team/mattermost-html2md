package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/JohannesKaufmann/html-to-markdown/plugin"
	"github.com/emersion/go-imap/client"
	rest "github.com/go-micro/plugins/v4/server/http"
	"github.com/go-resty/resty/v2"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/server"
	"golang.org/x/sync/errgroup"

	"github.com/x-research-team/mattermost-html2md/internal/api"
	"github.com/x-research-team/mattermost-html2md/internal/config"
	"github.com/x-research-team/mattermost-html2md/internal/flags"
	"github.com/x-research-team/mattermost-html2md/internal/services/mailbox"
	"github.com/x-research-team/mattermost-html2md/internal/services/mattermost"
	"github.com/x-research-team/mattermost-html2md/pkg/log"
	"github.com/x-research-team/mattermost-html2md/pkg/slice"
)

func Run(ctx context.Context, l zerolog.Logger) error {
	ctx, stop := signal.NotifyContext(ctx, os.Kill, os.Interrupt)
	defer stop()

	group, ctx := errgroup.WithContext(ctx)

	cfg, err := config.Load(&l, ".env")
	if err != nil {
		return fmt.Errorf("load: %w", err)
	}

	logger.DefaultLogger = log.NewMicroAdapter(l)

	converter := md.NewConverter("", true, nil)
	converter.Use(plugin.GitHubFlavored())
	cli := model.NewAPIv4Client(cfg.Mattermost.URL)
	cli.SetToken(cfg.Mattermost.Token)

	service := mattermost.New(cfg, converter, cli, resty.NewWithClient(&http.Client{
		Timeout: cfg.Mattermost.Timeout,
	}))

	adapter := log.NewCronAdapter(l)

	scheduler := cron.New(
		cron.WithChain(cron.Recover(adapter)),
		cron.WithLogger(adapter),
	)
	if _, err = scheduler.AddFunc(cfg.Cron.Interval, func() {
		c, err := client.DialTLS(net.JoinHostPort(cfg.IMAP.Host, strconv.Itoa(cfg.IMAP.Port)), nil)
		if err != nil {
			l.Error().Err(err).Msg("dial")
		}

		if err := c.Login(cfg.IMAP.User, cfg.IMAP.Pass); err != nil {
			l.Err(err).Msg("login")
		}

		mail := mailbox.New(cfg, c)
		if err := mail.Handle(ctx, service.SendWebhook); err != nil {
			l.Error().Err(err).Msg("handle")
		}

		if err := c.Logout(); err != nil {
			l.Error().Err(err).Msg("logout")
		}
	}); err != nil {
		return fmt.Errorf("add func: %w", err)
	}

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
		micro.Flags(slice.Merge(flags.Test)...),
	)

	svc.Init()

	group.Go(start(ctx, l, svc))
	group.Go(shutdown(ctx, l))
	group.Go(func() error {
		scheduler.Start()
		return nil
	})

	if err := group.Wait(); err != nil {
		return fmt.Errorf("wait: %w", err)
	}

	return nil
}

func start(_ context.Context, l zerolog.Logger, svc micro.Service) func() error {
	return func() error {
		l.Info().Msg("starting")
		return svc.Run()
	}
}

func shutdown(ctx context.Context, l zerolog.Logger) func() error {
	return func() error {
		<-ctx.Done()

		l.Info().Msg("shutting down")

		return ctx.Err()
	}
}
