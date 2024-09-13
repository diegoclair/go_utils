package resterrors

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/diegoclair/go_utils/resterrors/internal/pb"
	"github.com/diegoclair/goswag/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RestErr interface
type RestErr interface {
	Message() string
	StatusCode() int
	Error() string
	GetError() string
	Causes() interface{}
}

type restErr struct {
	ErrMessage    string `json:"message"`
	ErrStatusCode int    `json:"status_code"`
	ErrError      string `json:"error"`
	ErrCauses     any    `json:"causes"`
}

func (e restErr) Message() string {
	return e.ErrMessage
}

func (e restErr) StatusCode() int {
	return e.ErrStatusCode
}

func (e restErr) GetError() string {
	return e.ErrError
}

func (e restErr) Error() string {
	return fmt.Sprintf("message: %s - statusCode: %d - error: %s - causes: %v", e.ErrMessage, e.ErrStatusCode, e.ErrError, e.ErrCauses)
}

func (e restErr) Causes() any {
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

// GoSwagDefaultResponseErrors returns the default response errors for the GoSwag library
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

// ToPb convert a RestErr to a gRPC error if the error is a RestErr
// It returns the original error if the error is not a RestErr
func ToPb(err error) (res error) {
	if err == nil {
		return nil
	}

	restErr, ok := err.(restErr)
	if !ok {
		return err
	}

	causes, jsonErr := json.Marshal(restErr.Causes())
	if jsonErr != nil {
		return err
	}

	pbErr := &pb.RestError{
		Message:    restErr.Message(),
		StatusCode: int32(restErr.StatusCode()),
		Error:      restErr.GetError(),
		Causes:     causes,
	}

	st := status.New(codes.Code(restErr.StatusCode()), restErr.Message())
	st, stErr := st.WithDetails(pbErr)
	if stErr != nil {
		return err
	}

	return st.Err()
}

// FromError try to convert an error to RestError
// It supports errors from grpc (protoRestError) and resterrors
// If the error is not a RestError, it returns the original error
func FromError(err error) error {
	if err == nil {
		return nil
	}

	st, stOk := status.FromError(err)

	restErr, restErrOk := err.(restErr)

	if !stOk && !restErrOk {
		return err
	}

	if stOk {
		if len(st.Details()) > 0 {
			return fromPb(st.Details()[0].(*pb.RestError))
		}
		return err
	}

	return fromRestErr(restErr)
}

func fromPb(pbErr *pb.RestError) error {
	var causes any
	err := json.Unmarshal(pbErr.GetCauses(), &causes)
	if err != nil {
		return err
	}

	return &restErr{
		ErrMessage:    pbErr.GetMessage(),
		ErrStatusCode: int(pbErr.GetStatusCode()),
		ErrError:      pbErr.GetError(),
		ErrCauses:     causes,
	}
}

func fromRestErr(err restErr) error {
	return &restErr{
		ErrMessage:    err.Message(),
		ErrStatusCode: err.StatusCode(),
		ErrError:      err.GetError(),
		ErrCauses:     err.Causes(),
	}
}
