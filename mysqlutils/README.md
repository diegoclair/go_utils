# mysqlutils Package

## Description

This package provides a function for handling MySQL errors and converting them into [RESTful errors package](../resterrors/README.md).

## Constants

- `errorNoRows`: Represents the error message when no rows are returned from a query.
- `duplicatedKeyCode`: The MySQL error code for a duplicate key entry.

## Functions

### HandleMySQLError

The `HandleMySQLError` function takes a MySQL error as input and returns a RESTful error. It checks if the error is a MySQL error and if it is, it checks the error number and returns a corresponding RESTful error. If the error is not a MySQL error, it checks if the error message contains "no rows in result set" and returns a NotFoundError. Otherwise, it returns an InternalServerError.

### between

The `between` function takes a string and two delimiters as input and returns the substring between the two delimiters.

### before

The `before` function takes a string and a delimiter as input and returns the substring before the delimiter.

### after

The `after` function takes a string and a delimiter as input and returns the substring after the delimiter.