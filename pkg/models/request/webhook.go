package request

type Webhook struct {
	Body struct {
		Text    string `json:"text" required:"true" example:"Hello World"`
		Channel string `json:"channel" required:"true" example:"test"`
	}
}
