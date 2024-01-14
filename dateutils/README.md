# dateutils Package

## Description

This package provides functions for handling and formatting dates and times.

## Constants

- `apiDateLayout`: The layout for the complete date and time in the format "2006-01-02T15:04:05Z".
- `onlyDateLayout`: The layout for the date in the format "2006-01-02".
- `apiDBLayout`: The layout for the complete date and time in the format "2006-01-02 15:04:05", typically used for database operations.

## Functions

### GetDateNowTime

The `GetDateNowTime` function returns the current date and time as a `time.Time` value in UTC.

### GetCompleteDateNowString

The `GetCompleteDateNowString` function returns the current date and time as a string in the `apiDateLayout` format.

### GetCompleteDateNowDBLayout

The `GetCompleteDateNowDBLayout` function returns the current date and time as a string in the `apiDBLayout` format, typically used for database operations.

### GetOnlyDateNowString

The `GetOnlyDateNowString` function returns the current date as a string in the `onlyDateLayout` format.

### GetFirstAndLastOfAMonth

The `GetFirstAndLastOfAMonth` function takes a month as a `time.Month` value and returns the first and last day of that month as strings in the `onlyDateLayout` format.