package resterrors

import (
	"fmt"
	"net/http"
)

// RestErr interface
type RestErr interface {
	Message() string
	StatusCode() int
	Error() string
}

type restErr struct {
	message    string `json:"message"`
	statusCode int    `json:"status_code"`
	error      string `json:"error"`
}

func (e restErr) Message() string {
	return e.message
}
func (e restErr) StatusCode() int {
	return e.statusCode
}

func (e restErr) Error() string {
	return fmt.Sprintf("message: %s - statusCode: %d - error: %s", e.message, e.statusCode, e.error)
}

// NewRestError returns a instace of type restErr
func NewRestError(message string, status int, err string) RestErr {
	return restErr{
		message:    message,
		statusCode: status,
		error:      err,
	}
}

// NewBadRequestError returns a bad_request code error with your string message error
func NewBadRequestError(message string) RestErr {
	return restErr{
		message:    message,
		statusCode: http.StatusBadRequest,
		error:      "bad_request",
	}
}

// NewNotFoundError returns a not_found code error with your string message error
func NewNotFoundError(message string) RestErr {
	return restErr{
		message:    message,
		statusCode: http.StatusNotFound,
		error:      "not_found",
	}
}

// NewInternalServerError returns a not_found code error with your string message error
func NewInternalServerError(message string) RestErr {
	return &restErr{
		message:    message,
		statusCode: http.StatusInternalServerError,
		error:      "internal_server_error",
	}
}

// NewUnauthorizedError returns a unauthorized code error with your string message error
func NewUnauthorizedError(message string) RestErr {
	return &restErr{
		message:    message,
		statusCode: http.StatusUnauthorized,
		error:      "unauthorized",
	}
}
