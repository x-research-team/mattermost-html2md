package request

type Send struct {
	Headers struct{}
	URI     struct{}
	Query   struct{}
	Body    struct {
		Text string `json:"text" required:"true" example:"Hello World"`
	}
}
