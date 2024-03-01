package mailbox

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message"
	"golang.org/x/net/html"

	"github.com/x-research-team/mattermost-html2md/internal/config"
)

type mailbox struct {
	cfg    *config.Config
	client *client.Client
}

func New(cfg *config.Config, c *client.Client) *mailbox {
	return &mailbox{
		cfg:    cfg,
		client: c,
	}
}

func (m *mailbox) Handle(ctx context.Context, send func(context.Context, string, string) error) error {
	mbox, err := m.client.Select("INBOX", false)
	if err != nil {
		return fmt.Errorf("select INBOX: %w", err)
	}

	uids, err := m.client.Search(&imap.SearchCriteria{WithoutFlags: []string{imap.SeenFlag}})
	if err != nil {
		return fmt.Errorf("search: %w", err)
	}

	if len(uids) > 0 {
		seqset := new(imap.SeqSet)
		seqset.AddNum(uids...)
		seqset.AddRange(mbox.UidNext, mbox.UnseenSeqNum)

		messages := make(chan *imap.Message, 10)
		done := make(chan error, 1)
		go func() {
			done <- m.client.Fetch(seqset, []imap.FetchItem{imap.FetchRFC822}, messages)
		}()

		for msg := range messages {
			for _, r := range msg.Body {
				entity, err := message.Read(r)
				if err != nil {
					return fmt.Errorf("read message: %w", err)
				}

				ct := entity.Header.Get("Content-Type")
				switch {
				case strings.HasPrefix(ct, "multipart"):
					multiPartReader := entity.MultipartReader()

					for {
						p, err := multiPartReader.NextPart()
						if err == io.EOF {
							break
						}

						kind, _, err := p.Header.ContentType()
						if err != nil {
							return fmt.Errorf("content type: %w", err)
						}

						if kind != "text/html" {
							continue
						}

						body, err := io.ReadAll(p.Body)
						if err != nil {
							return fmt.Errorf("read body: %w", err)
						}

						doc, err := html.Parse(strings.NewReader(strings.TrimSpace(string(body))))
						if err != nil {
							return fmt.Errorf("parse html: %w", err)
						}

						var buf bytes.Buffer
						if err := html.Render(&buf, doc); err != nil {
							return fmt.Errorf("render html: %w", err)
						}

						if err := send(ctx, buf.String(), m.cfg.Mattermost.Channel); err != nil {
							return fmt.Errorf("send: %w", err)
						}

						// Mark the message as read
						item := imap.FormatFlagsOp(imap.AddFlags, true)

						if err := m.client.Store(seqset, item, nil, messages); err != nil {
							return fmt.Errorf("store: %w", err)
						}
					}
				}
			}
		}

		if err := <-done; err != nil {
			return fmt.Errorf("fetch: %w", err)
		}
	}

	return nil
}
