package resterrors

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBadRequestError(t *testing.T) {
	err := NewBadRequestError("this is the message", nil)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
	assert.EqualValues(t, "this is the message", err.Message())
}

func TestNewNotFoundError(t *testing.T) {
	err := NewNotFoundError("this is the message", nil)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
	assert.EqualValues(t, "this is the message", err.Message())
}

func TestNewInternalServerError(t *testing.T) {
	err := NewInternalServerError("this is the message", nil)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode())
	assert.EqualValues(t, "this is the message", err.Message())
}
