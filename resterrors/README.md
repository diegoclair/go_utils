# Resterrors Package

## Description

This package provides a structure and functions for handling errors in a RESTful format.

## Structures

### restErr

The `restErr` structure is used to represent errors that occur during the processing of a REST request. It contains the following fields:

- `ErrMessage`: The error message.
- `ErrStatusCode`: The HTTP status code associated with the error.
- `ErrError`: A string representing the error.
- `ErrCauses`: Any additional causes that may have contributed to the error.

## Functions

### NewRestError

The `NewRestError` function returns a new instance of `restErr`.

### NewRestErrorFromBytes

The `NewRestErrorFromBytes` function attempts to deserialize a `restErr` from a byte slice. If the deserialization fails, it returns an error.

### NewBadRequestError

The `NewBadRequestError` function returns a new `restErr` with the HTTP status code set to 400 (Bad Request).

### NewNotFoundError

The `NewNotFoundError` function returns a new `restErr` with the HTTP status code set to 404 (Not Found).

### NewInternalServerError

The `NewInternalServerError` function returns a new `restErr` with the HTTP status code set to 500 (Internal Server Error).

### NewUnauthorizedError

The `NewUnauthorizedError` function returns a new `restErr` with the HTTP status code set to 401 (Unauthorized).

### NewUnprocessableEntity

The `NewUnprocessableEntity` function returns a new `restErr` with the HTTP status code set to 422 (Unprocessable Entity).

### NewConflictError

The `NewConflictError` function returns a new `restErr` with the HTTP status code set to 409 (Conflict).