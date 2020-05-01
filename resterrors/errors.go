package resterrors

import "net/http"

// RestErr struct
type RestErr struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
}

// NewRestError returns a instace of type restErr
func NewRestError(message string, status int, err string) *RestErr {
	return &RestErr{
		Message:    message,
		StatusCode: status,
		Error:      err,
	}
}

// NewBadRequestError returns a bad_request code error with your string message error
func NewBadRequestError(message string) *RestErr {
	return &RestErr{
		Message:    message,
		StatusCode: http.StatusBadRequest,
		Error:      "bad_request",
	}
}

// NewNotFoundError returns a not_found code error with your string message error
func NewNotFoundError(message string) *RestErr {
	return &RestErr{
		Message:    message,
		StatusCode: http.StatusNotFound,
		Error:      "not_found",
	}
}

// NewInternalServerError returns a not_found code error with your string message error
func NewInternalServerError(message string) *RestErr {
	return &RestErr{
		Message:    message,
		StatusCode: http.StatusInternalServerError,
		Error:      "internal_server_error",
	}
}

// NewUnauthorizedError returns a unauthorized code error with your string message error
func NewUnauthorizedError(message string) *RestErr {
	return &RestErr{
		Message:    message,
		StatusCode: http.StatusUnauthorized,
		Error:      "unauthorized",
	}
}
