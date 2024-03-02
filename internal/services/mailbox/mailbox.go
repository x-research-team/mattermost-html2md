package mailbox

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message"
	"github.com/rs/zerolog"
	"golang.org/x/net/html"
	"golang.org/x/text/encoding/charmap"

	"github.com/x-research-team/mattermost-html2md/internal/config"
)

var encodings = map[string]*charmap.Charmap{
	"windows-1251": charmap.Windows1251,
	"windows-1252": charmap.Windows1252,
	"koi8-r":       charmap.KOI8R,
	"koi8-u":       charmap.KOI8U,
	"iso-8859-1":   charmap.ISO8859_1,
	"iso-8859-2":   charmap.ISO8859_2,
	"iso-8859-3":   charmap.ISO8859_3,
	"iso-8859-4":   charmap.ISO8859_4,
	"iso-8859-5":   charmap.ISO8859_5,
	"iso-8859-6":   charmap.ISO8859_6,
	"iso-8859-7":   charmap.ISO8859_7,
	"iso-8859-8":   charmap.ISO8859_8,
	"iso-8859-9":   charmap.ISO8859_9,
	"iso-8859-10":  charmap.ISO8859_10,
	"iso-8859-13":  charmap.ISO8859_13,
	"iso-8859-14":  charmap.ISO8859_14,
	"iso-8859-15":  charmap.ISO8859_15,
	"iso-8859-16":  charmap.ISO8859_16,
}

type mailbox struct {
	cfg    *config.Config
	client *client.Client
	logger *zerolog.Logger
}

func New(cfg *config.Config, c *client.Client, l *zerolog.Logger) *mailbox {
	return &mailbox{
		cfg:    cfg,
		client: c,
		logger: l,
	}
}

func (m *mailbox) Handle(ctx context.Context, send func(context.Context, string, string) error) error {
	_, err := m.client.Select("INBOX", false)
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

		messages := make(chan *imap.Message, 10)
		done := make(chan error, 1)
		go func() {
			done <- m.client.Fetch(seqset, []imap.FetchItem{imap.FetchRFC822, imap.FetchBodyStructure}, messages)
		}()

		for msg := range messages {
			for _, literal := range msg.Body {
				entity, err := message.Read(literal)
				if err != nil {
					return fmt.Errorf("read message: %w", err)
				}

				ct := entity.Header.Get("Content-Type")
				switch {
				case strings.HasPrefix(ct, "text"):
					buffer, err := io.ReadAll(entity.Body)
					if err != nil {
						return fmt.Errorf("read body: %w", err)
					}

					result := string(buffer)
					if result == "" {
						continue
					}

					if err := send(ctx, result, m.cfg.Mattermost.Channel); err != nil {
						return fmt.Errorf("send: %w", err)
					}
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

						switch kind {
						case "text/html":
							body, err := io.ReadAll(p.Body)
							if err != nil {
								return fmt.Errorf("read body: %w", err)
							}

							doc, err := html.Parse(strings.NewReader(strings.TrimSpace(string(body))))
							if err != nil {
								return fmt.Errorf("parse html: %w", err)
							}

							removeBr(doc)
							removeImages(doc)

							var buf bytes.Buffer
							if err := html.Render(&buf, doc); err != nil {
								return fmt.Errorf("render html: %w", err)
							}

							var utf8Str string
							match := regexp.MustCompile(`charset="([^"]+)`).FindStringSubmatch(p.Header.Get("Content-Type"))
							if len(match) > 1 {
								enc, ok := encodings[strings.ToLower(match[1])]

								if ok {
									decoder := enc.NewDecoder()
									utf8Str, err = decoder.String(buf.String())
									if err != nil {
										return fmt.Errorf("convert string: %w", err)
									}
								} else {
									if match[1] != "utf-8" {
										m.logger.Info().Str("charset", match[1]).Msg("unknown charset")
									}
									utf8Str = buf.String()
								}
							}

							utf8Str = strings.ReplaceAll(utf8Str, "|", "")
							utf8Str = strings.Replace(utf8Str, "---", "", -1)
							utf8Str = strings.ReplaceAll(utf8Str, "<br>", "")
							utf8Str = strings.ReplaceAll(utf8Str, "<br/>", "")
							utf8Str = strings.ReplaceAll(utf8Str, "б═", " ")
							utf8Str = strings.TrimSpace(utf8Str)

							m.logger.Debug().
								Str("charset", match[1]).
								Str("content-type", p.Header.Get("Content-Type")).
								Str("body", utf8Str)

							if err := send(ctx, utf8Str, m.cfg.Mattermost.Channel); err != nil {
								return fmt.Errorf("send: %w", err)
							}
						default:
							continue
						}
						break
					}
				}
				break
			}
		}

		if err := <-done; err != nil {
			return fmt.Errorf("fetch: %w", err)
		}
	}

	return nil
}

func removeBr(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "br" {
		parent := n.Parent
		if parent != nil {
			parent.RemoveChild(n)
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		removeBr(c)
	}
}

func removeImages(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "img" {
		parent := n.Parent
		if parent != nil {
			parent.RemoveChild(n)
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		removeImages(c)
	}
}
