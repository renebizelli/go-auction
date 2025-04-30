package internal_error

type InternalError struct {
	Err     string `json:"err"`
	Message string `json:"message"`
}

func (e *InternalError) Error() string {
	return e.Message
}

func NewInternalServerError(message string) *InternalError {
	return &InternalError{
		Err:     "internal_server_error",
		Message: message,
	}
}

func NewNotFoundError(message string) *InternalError {
	return &InternalError{
		Err:     "not_found",
		Message: message,
	}
}

func NewBadRequestError(message string) *InternalError {
	return &InternalError{
		Err:     "bad_request",
		Message: message,
	}
}
