# mysqlutils Package

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
