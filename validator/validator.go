package validator

import (
	"fmt"

	"github.com/diegoclair/go_utils-lib/v2/resterrors"
	"github.com/go-playground/validator/v10"
	"github.com/klassmann/cpfcnpj"
)

type Validator interface {
	ValidateStruct(dataSet interface{}) error
}

type validatorImpl struct {
	validator *validator.Validate
}

// NewValidator returns a new instance of validator interface with the custom validations:
// cpf - validate if the input is a valid cpf
// cnpj - validate if the input is a valid cnpj
func NewValidator() (Validator, error) {
	v := validator.New()

	err := v.RegisterValidation("cpf", func(fl validator.FieldLevel) bool {
		cpf := cpfcnpj.NewCPF(fl.Field().String())
		return cpf.IsValid()
	})
	if err != nil {
		return nil, resterrors.NewInternalServerError("Error trying to register cpf validation", err)
	}

	err = v.RegisterValidation("cnpj", func(fl validator.FieldLevel) bool {
		cnpj := cpfcnpj.NewCNPJ(fl.Field().String())
		return cnpj.IsValid()
	})
	if err != nil {
		return nil, resterrors.NewInternalServerError("Error trying to register cnpj validation", err)
	}

	return &validatorImpl{
		validator: v,
	}, nil
}

// ValidateStruct validates the given data set using the validator instance.
// It returns an error if the validation fails, with detailed error messages for each validation rule that was not satisfied.
// The error message includes information about the field name and the specific validation rule that failed.
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
