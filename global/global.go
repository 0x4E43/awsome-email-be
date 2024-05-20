package global

type Response struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Data    any    `json:"data"`
}

func PrepareResponse(message string, status int, data any) *Response {
	res := Response{
		Message: message,
		Status:  status,
		Data:    data,
	}
	return &res
}
