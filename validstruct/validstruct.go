package validstruct

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/diegoclair/go_utils-lib/resterrors"
	"github.com/go-playground/validator/v10"
)

// ValidateStruct - valid
func ValidateStruct(dataSet interface{}) resterrors.RestErr {

	validate := validator.New()

	err := validate.Struct(dataSet)
	if err != nil {

		invalidArgument, ok := err.(*validator.InvalidValidationError)
		if ok {
			return resterrors.NewInternalServerError("Invalid argument passed to struct: " + fmt.Sprint(invalidArgument))
		}

		reflected := reflect.ValueOf(dataSet)

		var name string
		var errMessage []string

		for _, err := range err.(validator.ValidationErrors) {

			field, _ := reflected.Type().FieldByName(err.StructField())

			name = field.Tag.Get("json")
			if name == "" {
				name = strings.ToLower(err.StructField())
			}
			fmt.Println(field.Type)
			switch err.Tag() {
			case "required":
				errMessage = append(errMessage, fmt.Sprintf("The field '%s' is required", name))

			case "email":
				errMessage = append(errMessage, fmt.Sprintf("The field '%s' should be a valid email", name))

			case "eqfield":
				errMessage = append(errMessage, fmt.Sprintf("The field '%s' should be equal to the %s", name, err.Param()))

			case "gte":
				errMessage = append(errMessage, fmt.Sprintf("The field '%s' should be greater than or equal %s", name, err.Param()))

			case "lte":
				errMessage = append(errMessage, fmt.Sprintf("The field '%s' should be less than or equal %s", name, err.Param()))

			case "max":
				errMessage = append(errMessage, fmt.Sprintf("The field '%s' should have the max lenhgt or value: %s", name, err.Param()))

			case "min":
				errMessage = append(errMessage, fmt.Sprintf("The field '%s' should have the minimun lenhgt or value: %s", name, err.Param()))

			default:
				errMessage = append(errMessage, fmt.Sprintf("The field '%s' is invalid.", name))
			}
		}

		return resterrors.NewUnprocessableEntity("Invalid input data", errMessage)
	}

	return nil
}
