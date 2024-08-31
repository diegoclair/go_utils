package resterrors

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/diegoclair/goswag/models"
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

// NewRestError returns a instance of type restErr
func NewRestError(message string, status int, err string, causes ...interface{}) RestErr {
	return restErr{
		ErrMessage:    message,
		ErrStatusCode: status,
		ErrError:      err,
		ErrCauses:     causes,
	}
}

func GoSwagDefaultResponseErrors() []models.ReturnType {
	return []models.ReturnType{
		{
			StatusCode: http.StatusBadRequest,
			Body:       restErr{},
		},
		{
			StatusCode: http.StatusNotFound,
			Body:       restErr{},
		},
		{
			StatusCode: http.StatusInternalServerError,
			Body:       restErr{},
		},
		{
			StatusCode: http.StatusUnauthorized,
			Body:       restErr{},
		},
		{
			StatusCode: http.StatusUnprocessableEntity,
			Body:       restErr{},
		},
		{
			StatusCode: http.StatusConflict,
			Body:       restErr{},
		},
	}
}

// NewRestErrorFromBytes returns a RestErr response format by result from another api using resterror response
// example: if you call to another api and this api response a resterror response format, you can use this function to convert the response to RestErr  with response.Bytes()
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
		ErrError:      http.StatusText(http.StatusBadRequest),
		ErrCauses:     causes,
	}
}

// NewNotFoundError returns a not_found code error with your string message error
func NewNotFoundError(message string, causes ...interface{}) RestErr {
	return restErr{
		ErrMessage:    message,
		ErrStatusCode: http.StatusNotFound,
		ErrError:      http.StatusText(http.StatusNotFound),
		ErrCauses:     causes,
	}
}

// NewInternalServerError returns a not_found code error with your string message error
func NewInternalServerError(message string, causes ...interface{}) RestErr {
	return &restErr{
		ErrMessage:    message,
		ErrStatusCode: http.StatusInternalServerError,
		ErrError:      http.StatusText(http.StatusInternalServerError),
		ErrCauses:     causes,
	}
}

// NewUnauthorizedError returns a unauthorized code error with your string message error
func NewUnauthorizedError(message string, causes ...interface{}) RestErr {
	return &restErr{
		ErrMessage:    message,
		ErrStatusCode: http.StatusUnauthorized,
		ErrError:      http.StatusText(http.StatusUnauthorized),
		ErrCauses:     causes,
	}
}

// NewUnprocessableEntity returns a unprocessable_entity code error with your string message error
func NewUnprocessableEntity(message string, causes ...interface{}) RestErr {
	result := restErr{
		ErrMessage:    message,
		ErrStatusCode: http.StatusUnprocessableEntity,
		ErrError:      http.StatusText(http.StatusUnprocessableEntity),
		ErrCauses:     causes,
	}

	return result
}

func NewConflictError(message string, causes ...interface{}) RestErr {
	result := restErr{
		ErrMessage:    message,
		ErrStatusCode: http.StatusConflict,
		ErrError:      http.StatusText(http.StatusConflict),
		ErrCauses:     causes,
	}

	return result
}
