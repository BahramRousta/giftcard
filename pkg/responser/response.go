package responser

type Response struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
	Data    any    `json:"data"`
}

type GiftCardResponse struct {
	StatusCode int
	Header     any
	Body       string
}
