package mattermost

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/x-research-team/mattermost-html2md/internal/config"

	md "github.com/JohannesKaufmann/html-to-markdown"
)

type Service interface {
	Send(ctx context.Context, text, channel string) error
}

type service struct {
	cfg       *config.Config
	converter *md.Converter
	client    *model.Client4
}

func New(cfg *config.Config, converter *md.Converter, client *model.Client4) Service {
	return &service{
		cfg:       cfg,
		converter: converter,
		client:    client,
	}
}

func (s service) Send(ctx context.Context, html, channel string) error {
	result, err := s.converter.ConvertString(html)
	if err != nil {
		return fmt.Errorf("convert string: %w", err)
	}

	_, resp, err := s.client.CreatePost(&model.Post{
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
