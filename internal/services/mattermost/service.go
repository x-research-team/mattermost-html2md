package mattermost

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/go-resty/resty/v2"
	"github.com/mattermost/mattermost-server/v6/model"

	"github.com/x-research-team/mattermost-html2md/internal/config"
	"github.com/x-research-team/mattermost-html2md/pkg/models"
)

type Service interface {
	SendAPI(ctx context.Context, text, channel string) error
	SendWebhook(ctx context.Context, text, channel string) error
}

type service struct {
	cfg       *config.Config
	converter *md.Converter
	api       *model.Client4
	client    *resty.Client
}

func New(cfg *config.Config, converter *md.Converter, api *model.Client4, client *resty.Client) Service {
	return &service{
		cfg:       cfg,
		converter: converter,
		api:       api,
		client:    client.SetDebug(cfg.Mattermost.Debug),
	}
}

func (s service) SendAPI(ctx context.Context, html, channel string) error {
	result, err := s.converter.ConvertString(html)
	if err != nil {
		return fmt.Errorf("convert string: %w", err)
	}

	_, resp, err := s.api.CreatePost(&model.Post{
		ChannelId: channel,
		Message:   result,
	})

	if err != nil {
		return fmt.Errorf("create post: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return errors.New("create post: " + fmt.Sprintf("status code: %d", resp.StatusCode))
	}

	return nil
}

func (s *service) SendWebhook(ctx context.Context, text string, channel string) error {
	text, err := s.converter.ConvertString(text)
	if err != nil {
		return fmt.Errorf("convert string: %w", err)
	}

	text = strings.ReplaceAll(text, "<br>", "")
	text = strings.ReplaceAll(text, "<br/>", "")

	resp, err := s.client.R().
		EnableTrace().
		SetDebug(s.cfg.Mattermost.Debug).
		SetBody(models.Webhook{
			Text: text,
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
