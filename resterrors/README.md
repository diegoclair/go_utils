# Resterrors Package

> **⚠️ DEPRECATED**
>
> This package is **deprecated**. New code should use [`github.com/diegoclair/apperr`](https://github.com/diegoclair/apperr), which is transport-agnostic: business code returns `Kind` + `Code`, and the transport layer maps to HTTP (via [`apperr/httpmap`](https://github.com/diegoclair/apperr/tree/main/httpmap)), gRPC, GraphQL, etc.
>
> Benefits of migrating:
> - Business code stops importing `net/http` and HTTP status codes.
> - Errors carry stable `Code` strings (e.g. `"USER_NOT_FOUND"`) — frontend can do i18n.
> - Structured `meta` for dynamic data (e.g. `retry_after_minutes`, validation `fields`).
> - Works with `errors.Is` / `errors.As` natively.
>
> This package remains here for backward compatibility with existing consumers (e.g. `lybel`, `rrr`). For validation specifically, see [`appvalidator/apperrmap`](https://github.com/diegoclair/appvalidator/tree/main/apperrmap).
>
> ### Migration cheat sheet
>
> | Old (resterrors)                | New (apperr)                                            |
> |---------------------------------|---------------------------------------------------------|
> | `NewBadRequestError`            | `apperr.ErrValidation` / `apperr.ErrInvalidInput`       |
> | `NewNotFoundError`              | `apperr.ErrNotFound` / `apperr.ErrRecordNotFound`       |
> | `NewInternalServerError`        | `apperr.ErrInternal.Wrap(err)`                          |
> | `NewUnauthorizedError`          | `apperr.ErrUnauthenticated` / `apperr.ErrTokenExpired`  |
> | `NewUnprocessableEntity`        | `apperr.ErrValidation.WithMeta("fields", …)`            |
> | `NewConflictError`              | `apperr.ErrConflict` / `apperr.ErrDuplicateEntry`       |
> | `RestErr` (interface)           | `apperr.AppError` (interface)                           |

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
