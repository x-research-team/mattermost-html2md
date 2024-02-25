package response

type E struct {
	Message string `json:"message" example:"Internal Server Error"`
}

func Error(err error) *E {
	return &E{
		Message: err.Error(),
	}
}
