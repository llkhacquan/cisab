package service

var ErrNotFound = Error{
	Code:    404,
	Message: "not found",
}

var ErrUnauthorized = Error{
	Code:    401,
	Message: "unauthorized",
}

func NewInvalidInputError(message string) Error {
	return Error{
		Code:    400,
		Message: message,
	}
}

type Error struct {
	Code    int    `json:"code"` // HTTP status code
	Message string `json:"message"`
}

func (s Error) Error() string {
	return s.Message
}
