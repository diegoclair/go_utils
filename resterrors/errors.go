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
	ErrMessage    string `json:"message"`
	ErrStatusCode int    `json:"status_code"`
	ErrError      string `json:"error"`
}

func (e restErr) Message() string {
	return e.ErrMessage
}
func (e restErr) StatusCode() int {
	return e.ErrStatusCode
}

func (e restErr) Error() string {
	return fmt.Sprintf("message: %s - statusCode: %d - error: %s", e.ErrMessage, e.ErrStatusCode, e.ErrError)
}

// NewRestError returns a instace of type restErr
func NewRestError(message string, status int, err string) RestErr {
	return restErr{
		ErrMessage:    message,
		ErrStatusCode: status,
		ErrError:      err,
	}
}

// NewBadRequestError returns a bad_request code error with your string message error
func NewBadRequestError(message string) RestErr {
	return restErr{
		ErrMessage:    message,
		ErrStatusCode: http.StatusBadRequest,
		ErrError:      "bad_request",
	}
}

// NewNotFoundError returns a not_found code error with your string message error
func NewNotFoundError(message string) RestErr {
	return restErr{
		ErrMessage:    message,
		ErrStatusCode: http.StatusNotFound,
		ErrError:      "not_found",
	}
}

// NewInternalServerError returns a not_found code error with your string message error
func NewInternalServerError(message string) RestErr {
	return &restErr{
		ErrMessage:    message,
		ErrStatusCode: http.StatusInternalServerError,
		ErrError:      "internal_server_error",
	}
}

// NewUnauthorizedError returns a unauthorized code error with your string message error
func NewUnauthorizedError(message string) RestErr {
	return &restErr{
		ErrMessage:    message,
		ErrStatusCode: http.StatusUnauthorized,
		ErrError:      "unauthorized",
	}
}
