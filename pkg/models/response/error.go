package response

type Err struct {
	Message string `json:"message" example:"Internal Server Error"`
}

func Error(err error) *Err {
	return &Err{
		Message: err.Error(),
	}
}
