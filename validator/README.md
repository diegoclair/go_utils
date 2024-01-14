# validator Package

## Description

This package provides a custom validator for Go structs, with additional validations for Brazilian CPF and CNPJ numbers.  
The error returned is formated as the [resterrors package](../resterrors/README.md)

## Types

### Validator interface

The `Validator` interface defines a method for validating a struct.
The exported functions are:  
* `ValidateStruct(dataSet interface{}) error`

## Functions

### NewValidator

The `NewValidator` function returns a new instance of the `Validator` interface. It registers custom validations for CPF and CNPJ numbers.

### ValidateStruct

The `ValidateStruct` method validates the given data set using the validator instance. It returns an error if the validation fails, with detailed error messages for each validation rule that was not satisfied. The error message includes information about the field name and the specific validation rule that failed.

## Usage

To use this package, import it into your Go project and create a new validator with `NewValidator`. Then, use the `ValidateStruct` method to validate your structs.

## Contribution

Details on how to contribute to this project.

## License

Information about the license.