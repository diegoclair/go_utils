package validstruct

import (
	"fmt"

	"github.com/diegoclair/go_utils-lib/v2/resterrors"
	"github.com/go-playground/validator/v10"
)

// ValidateStruct - validate if the input is valid for requirements of a struct
func ValidateStruct(dataSet interface{}) resterrors.RestErr {

	validate := validator.New()

	err := validate.Struct(dataSet)
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

			default:
				errMessage = append(errMessage, fmt.Sprintf("The field '%s' is invalid.", name))
			}
		}

		return resterrors.NewUnprocessableEntity("Invalid input data", errMessage)
	}

	return nil
}
