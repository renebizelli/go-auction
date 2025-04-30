package rest_err

import "renebizelli/go-leilao/Internals/internal_error"

type RestErr struct {
	Message string   `json:"message"`
	Err     string   `json:"error"`
	Code    int      `json:"code"`
	Causes  []Causes `json:"causes"`
}

type Causes struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (r *RestErr) Error() string {
	return r.Message
}

func ConvertError(internal_error *internal_error.InternalError) *RestErr {
	switch internal_error.Err {
	case "bad_request":
		return NewBadRequestError(internal_error.Error())
	case "not_found":
		return NewNotFoundError(internal_error.Error())
	default:
		return InternalServerError(internal_error.Error())
	}
}

func NewBadRequestError(message string, causes ...Causes) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "bad_request",
		Code:    400,
		Causes:  causes,
	}
}

func NewNotFoundError(message string, causes ...Causes) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "not_found",
		Code:    404,
		Causes:  causes,
	}
}

func InternalServerError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "internal_server_error",
		Code:    500,
	}
}
