# validator Package

## Description

This package provides a custom validator for Go Structs, with additional validations:
* 1 - tag `cpf`: Validate CPF number with brazilian rules
* 2 - tag `cnpj`: Validate CNPJ number with brazilian rules
* 3 - tag `required_trim`: Validate tag required after trim the string value. Only valid for strings   
  
It use the `go-playground/validator/v10` lib to do the validations.

## Types

### Validator interface

The `Validator` interface defines a method for validating a struct.
- Exported custom functions:  
    * `ValidateStruct(dataSet interface{}) error`
        * this function with use the validator
- Default functions exported from `go-playground/validator/v10`
    - Var(field interface{}, tag string) error
	- RegisterValidation(tag string, fn validator.Func) error
	- RegisterAlias(alias string, tags string)
	- StructExcept(current interface{}, fields ...string) error
	- StructPartial(current interface{}, fields ...string) error
	- StructFiltered(current interface{}, filter validator.FilterFunc) error

## Functions

### NewValidator

The `NewValidator` function returns a new instance of the `Validator` interface. It registers custom validations for CPF and CNPJ numbers.

### ValidateStruct

The `ValidateStruct` method validates the given data set using the validator instance. It returns an error if the validation fails, with detailed error messages for each validation rule that was not satisfied. The error message includes information about the field name and the specific validation rule that failed.  
This function returns a error of type [resterrors.RestErr](../resterrors/README.md).

## Usage

To use this package, import it into your Go project and create a new validator with `NewValidator`. Then, use the `ValidateStruct` method to validate your structs.

## Contribution

Details on how to contribute to this project.

## License

Information about the license.