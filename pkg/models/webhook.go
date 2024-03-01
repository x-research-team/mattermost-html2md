package models

type Webhook struct {
	Text        string              `json:"text,omitempty"`
	Username    string              `json:"username,omitempty"`
	IconURL     string              `json:"icon_url,omitempty"`
	Channel     string              `json:"channel,omitempty"`
	Attachments []WebhookAttachment `json:"attachments,omitempty"`
}

type WebhookAttachment struct {
	Fallback   string                   `json:"fallback,omitempty"`
	Color      string                   `json:"color,omitempty"`
	Pretext    string                   `json:"pretext,omitempty"`
	Text       string                   `json:"text,omitempty"`
	Title      string                   `json:"title,omitempty"`
	TitleLink  string                   `json:"title_link,omitempty"`
	AuthorName string                   `json:"author_name,omitempty"`
	AuthorIcon string                   `json:"author_icon,omitempty"`
	AuthorLink string                   `json:"author_link,omitempty"`
	Fields     []WebhookAttachmentField `json:"fields,omitempty"`
	ImageURL   string                   `json:"image_url,omitempty"`
}

// WebhookAttachmentField contains attachment fields for usage in attachments
type WebhookAttachmentField struct {
	Short bool   `json:"short"`
	Title string `json:"title,omitempty"`
	Value string `json:"value,omitempty"`
}
