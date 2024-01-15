package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/diegoclair/go_utils-lib/v2/resterrors"
	"github.com/go-playground/validator/v10"
	"github.com/klassmann/cpfcnpj"
)

type Validator interface {
	// ValidateStruct validates the given data set using the validator instance.
	// It use the go-playground/validator/v10 package to validate the data set and return a better error message for some tags.
	// This function returns a error of type resterrors.RestErr.
	// It contains 3 more custom tags:
	// cpf - validate if the input is a valid cpf
	// cnpj - validate if the input is a valid cnpj
	// required_trim - validate the tag required after trim the input (only valid for string fields type)
	ValidateStruct(dataSet interface{}) error

	// Some default Methods from go-playground/validator/v10 package
	Var(field interface{}, tag string) error
	RegisterValidation(tag string, fn validator.Func) error
	RegisterAlias(alias string, tags string)
	StructExcept(current interface{}, fields ...string) error
	StructPartial(current interface{}, fields ...string) error
	StructFiltered(current interface{}, filter validator.FilterFunc) error
}

type validatorImpl struct {
	validator *validator.Validate
}

// NewValidator returns a new instance of validator interface with the custom validations tags validations.
func NewValidator() (Validator, error) {
	v := &validatorImpl{
		// the option validator.WithRequiredStructEnabled() will be default on v11 of go-playground/validator
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}

	err := v.registerCustomValidations()
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (v *validatorImpl) ValidateStruct(dataSet interface{}) error {
	err := v.validator.Struct(dataSet)
	if err != nil {

		invalidArgument, ok := err.(*validator.InvalidValidationError)
		if ok {
			return resterrors.NewInternalServerError("Invalid argument passed to struct: "+fmt.Sprint(invalidArgument), err)
		}

		var errMessage []string

		for _, err := range err.(validator.ValidationErrors) {

			name := err.StructField()

			switch err.Tag() {
			case "required":
				errMessage = append(errMessage, fmt.Sprintf("The field '%s' is required", name))

			case "email":
				errMessage = append(errMessage, fmt.Sprintf("The field '%s' should be a valid email", name))

			case "eq":
				errMessage = append(errMessage, fmt.Sprintf("The value '%s' should be equal to the %s", name, err.Param()))

			case "eqfield":
				errMessage = append(errMessage, fmt.Sprintf("The field '%s' should be equal to the field %s", name, err.Param()))

			case "ne":
				errMessage = append(errMessage, fmt.Sprintf("The value '%s' should not be equal to the %s", name, err.Param()))

			case "gte":
				errMessage = append(errMessage, fmt.Sprintf("The field '%s' should be greater than or equal %s", name, err.Param()))

			case "gt":
				errMessage = append(errMessage, fmt.Sprintf("The field '%s' should be greater than %s", name, err.Param()))

			case "lte":
				errMessage = append(errMessage, fmt.Sprintf("The field '%s' should be less than or equal %s", name, err.Param()))

			case "lt":
				errMessage = append(errMessage, fmt.Sprintf("The field '%s' should be less than %s", name, err.Param()))

			case "max":
				errMessage = append(errMessage, fmt.Sprintf("The field '%s' should have the max lenhgt or value: %s", name, err.Param()))

			case "min":
				errMessage = append(errMessage, fmt.Sprintf("The field '%s' should have the minimun lenhgt or value: %s", name, err.Param()))

			case "uuid4":
				errMessage = append(errMessage, fmt.Sprintf("The format of '%s' should be uuid4: %s", name, err.Param()))

			case "cpf":
				errMessage = append(errMessage, fmt.Sprintf("The field '%s' should be a valid cpf", name))

			case "cnpj":
				errMessage = append(errMessage, fmt.Sprintf("The field '%s' should be a valid cnpj", name))

			default:
				errMessage = append(errMessage, fmt.Sprintf("The field '%s' is invalid.", name))
			}
		}

		return resterrors.NewUnprocessableEntity("Invalid input data", errMessage)
	}

	return nil
}

func (v *validatorImpl) registerCustomValidations() error {
	err := v.validator.RegisterValidation("cpf", func(fl validator.FieldLevel) bool {
		cpf := cpfcnpj.NewCPF(fl.Field().String())
		return cpf.IsValid()
	})
	if err != nil {
		return resterrors.NewInternalServerError("Error trying to register cpf validation", err)
	}

	err = v.validator.RegisterValidation("cnpj", func(fl validator.FieldLevel) bool {
		cnpj := cpfcnpj.NewCNPJ(fl.Field().String())
		return cnpj.IsValid()
	})
	if err != nil {
		return resterrors.NewInternalServerError("Error trying to register cnpj validation", err)
	}

	err = v.validator.RegisterValidation("required_trim", func(fl validator.FieldLevel) bool {
		if fl.Field().Kind() != reflect.String {
			return false
		}

		trimmedValue := strings.TrimSpace(fl.Field().String())

		err = v.validator.Var(trimmedValue, "required")
		return err == nil
	})
	if err != nil {
		return resterrors.NewInternalServerError("Error trying to register required_trim validation", err)
	}

	return nil
}

func (v *validatorImpl) Var(field interface{}, tag string) error {
	return v.validator.Var(field, tag)
}

func (v *validatorImpl) RegisterValidation(tag string, fn validator.Func) error {
	return v.validator.RegisterValidation(tag, fn)
}

func (v *validatorImpl) RegisterAlias(alias string, tags string) {
	v.validator.RegisterAlias(alias, tags)
}

func (v *validatorImpl) StructExcept(current interface{}, fields ...string) error {
	return v.validator.StructExcept(current, fields...)
}

func (v *validatorImpl) StructPartial(current interface{}, fields ...string) error {
	return v.validator.StructPartial(current, fields...)
}

func (v *validatorImpl) StructFiltered(current interface{}, filter validator.FilterFunc) error {
	return v.validator.StructFiltered(current, filter)
}
