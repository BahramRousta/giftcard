package giftcard

type AuthErr struct {
	ErrMsg string
}

func (e *AuthErr) Error() string {
	return e.ErrMsg
}

type ForbiddenErr struct {
	ErrMsg string
}

func (e *ForbiddenErr) Error() string {
	return e.ErrMsg
}

type RequestErr struct {
	ErrMsg   string
	Response map[string]any
}

func (e *RequestErr) Error() string {
	return e.ErrMsg
}
