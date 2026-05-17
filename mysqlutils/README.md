# mysqlutils Package

> **⚠️ DEPRECATED**
>
> This package is **deprecated** and is no longer maintained. Projects have migrated:
> - From **MySQL** → **PostgreSQL** (so MySQL-specific error parsing is no longer needed).
> - From [`resterrors`](../resterrors/README.md) → [`apperr`](https://github.com/diegoclair/apperr) (transport-agnostic errors).
>
> **Recommended replacement:** handle database errors at the repository layer and return `apperr` definitions directly. For example:
>
> ```go
> if errors.Is(err, sql.ErrNoRows) {
>     return apperr.ErrNotFound
> }
> var pgErr *pgconn.PgError
> if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
>     return apperr.ErrDuplicateEntry
> }
> return apperr.ErrInternal.Wrap(err)
> ```

## Description

This package provides a function for handling MySQL errors and converting them into [RESTful errors package](../resterrors/README.md).


## Functions

### HandleMySQLError

The `HandleMySQLError` function handles the MySQL errors and returns a corresponding REST error.
It takes an error as input and checks if it is a MySQL error. If it is not a MySQL error,
it checks if the error message contains a specific string indicating a "no rows" error.
If it does, it returns a NotFoundError with a custom message. Otherwise, it returns an
InternalServerError with the error message.
If the error is a MySQL error, it checks the error number and handles specific cases.
For example, if the error number indicates a duplicated key error, it extracts the
duplicated key and value from the error message and returns a ConflictError with a
custom message. If none of the specific cases match, it returns an InternalServerError
with the error message.

### SQLNotFound 

The `SQLNotFound` checks if the given error message indicates that no SQL rows or records were found.
It returns true if no rows or records were found, otherwise false.

### between

The `between` function takes a string and two delimiters as input and returns the substring between the two delimiters.
