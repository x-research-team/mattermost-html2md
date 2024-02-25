package mattermost

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/x-research-team/mattermost-html2md/internal/config"
	"github.com/x-research-team/mattermost-html2md/pkg/models/request"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/go-resty/resty/v2"
)

type Service interface {
	Send(ctx context.Context, text string) error
}

type service struct {
	cfg       *config.Config
	converter *md.Converter
	client    *resty.Client
}

func New(cfg *config.Config, converter *md.Converter, client *resty.Client) Service {
	return &service{
		cfg:       cfg,
		converter: converter,
		client:    client,
	}
}

func (s service) Send(ctx context.Context, text string) error {
	text, err := s.converter.ConvertString(text)
	if err != nil {
		return fmt.Errorf("convert string: %w", err)
	}

	resp, err := s.client.R().
		EnableTrace().
		SetDebug(s.cfg.Mattermost.Debug).
		SetBody(request.Webhook{
			Text:     text,
			Username: s.cfg.Mattermost.User,
			Channel:  s.cfg.Mattermost.Channel,
		}).
		Post(s.cfg.Mattermost.Webhook)

	if err != nil {
		return fmt.Errorf("send: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("send: %w", errors.New(resp.String()))
	}

	return nil
}