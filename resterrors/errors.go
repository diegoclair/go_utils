package resterrors

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// RestErr interface
type RestErr interface {
	Message() string
	StatusCode() int
	Error() string
	Causes() interface{}
}

type restErr struct {
	ErrMessage    string      `json:"message"`
	ErrStatusCode int         `json:"status_code"`
	ErrError      string      `json:"error"`
	ErrCauses     interface{} `json:"causes"`
}

func (e restErr) Message() string {
	return e.ErrMessage
}

func (e restErr) StatusCode() int {
	return e.ErrStatusCode
}

func (e restErr) Error() string {
	return fmt.Sprintf("message: %s - statusCode: %d - error: %s - causes: %v", e.ErrMessage, e.ErrStatusCode, e.ErrError, e.ErrCauses)
}

func (e restErr) Causes() interface{} {
	return e.ErrCauses
}

// NewRestError returns a instace of type restErr
func NewRestError(message string, status int, err string, causes ...interface{}) RestErr {
	return restErr{
		ErrMessage:    message,
		ErrStatusCode: status,
		ErrError:      err,
		ErrCauses:     causes,
	}
}

// NewRestErrorFromBytes returns a RestErr response format by result from another api using resterror response
func NewRestErrorFromBytes(bytes []byte) (RestErr, error) {
	var apiErr restErr
	if err := json.Unmarshal(bytes, &apiErr); err != nil {
		return nil, errors.New("invalid restError json")
	}
	return apiErr, nil
}

// NewBadRequestError returns a bad_request code error with your string message error
func NewBadRequestError(message string, causes ...interface{}) RestErr {
	return restErr{
		ErrMessage:    message,
		ErrStatusCode: http.StatusBadRequest,
		ErrError:      "bad_request",
		ErrCauses:     causes,
	}
}

// NewNotFoundError returns a not_found code error with your string message error
func NewNotFoundError(message string, causes ...interface{}) RestErr {
	return restErr{
		ErrMessage:    message,
		ErrStatusCode: http.StatusNotFound,
		ErrError:      "not_found",
		ErrCauses:     causes,
	}
}

// NewInternalServerError returns a not_found code error with your string message error
func NewInternalServerError(message string, causes ...interface{}) RestErr {
	return &restErr{
		ErrMessage:    message,
		ErrStatusCode: http.StatusInternalServerError,
		ErrError:      "internal_server_error",
		ErrCauses:     causes,
	}
}

// NewUnauthorizedError returns a unauthorized code error with your string message error
func NewUnauthorizedError(message string, causes ...interface{}) RestErr {
	return &restErr{
		ErrMessage:    message,
		ErrStatusCode: http.StatusUnauthorized,
		ErrError:      "unauthorized",
		ErrCauses:     causes,
	}
}

// NewUnprocessableEntity returns a unprocessable_entity code error with your string message error
func NewUnprocessableEntity(message string, causes ...interface{}) RestErr {
	result := restErr{
		ErrMessage:    message,
		ErrStatusCode: http.StatusUnprocessableEntity,
		ErrError:      "unprocessable_entity",
		ErrCauses:     causes,
	}

	return result
}

func NewConflictError(message string, causes ...interface{}) RestErr {
	result := restErr{
		ErrMessage:    message,
		ErrStatusCode: http.StatusConflict,
		ErrError:      "conflict",
		ErrCauses:     causes,
	}

	return result
}
