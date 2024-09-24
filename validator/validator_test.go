package validator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_validatorImpl_ValidateStruct(t *testing.T) {
	ctx := context.Background()
	v, err := NewValidator()
	assert.NoError(t, err)

	type args struct {
		dataSet interface{}
	}
	tests := []struct {
		name       string
		args       args
		checkError func(err error)
	}{
		{
			name: "should not return error for a valid struct",
			args: args{
				dataSet: struct {
					Name  string `validate:"required"`
					Email string `validate:"required,email"`
				}{
					Name:  "John Doe",
					Email: "email@test.com",
				},
			},
			checkError: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "should return error for tag required",
			args: args{
				dataSet: struct {
					Name string `validate:"required"`
				}{
					Name: "",
				},
			},
			checkError: func(err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "The field 'Name' is required")
			},
		},
		{
			name: "should return error for tag email",
			args: args{
				dataSet: struct {
					Email string `validate:"email"`
				}{
					Email: "email",
				},
			},
			checkError: func(err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "The field 'Email' should be a valid email")
			},
		},
		{
			name: "should return error for tag eq",
			args: args{
				dataSet: struct {
					Age int `validate:"eq=18"`
				}{
					Age: 17,
				},
			},
			checkError: func(err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "The value 'Age' should be equal to the 18")
			},
		},
		{
			name: "should return error for tag eqfield",
			args: args{
				dataSet: struct {
					Age  int `validate:"eqfield=Name"`
					Name int
				}{
					Age:  17,
					Name: 18,
				},
			},
			checkError: func(err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "The field 'Age' should be equal to the field Name")
			},
		},
		{
			name: "should return error for tag ne",
			args: args{
				dataSet: struct {
					Age int `validate:"ne=18"`
				}{
					Age: 18,
				},
			},
			checkError: func(err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "The value 'Age' should not be equal to the 18")
			},
		},
		{
			name: "should return error for tag gte",
			args: args{
				dataSet: struct {
					Age int `validate:"gte=18"`
				}{
					Age: 17,
				},
			},
			checkError: func(err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "The field 'Age' should be greater than or equal 18")
			},
		},
		{
			name: "should return error for tag gt",
			args: args{
				dataSet: struct {
					Age int `validate:"gt=18"`
				}{
					Age: 18,
				},
			},
			checkError: func(err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "The field 'Age' should be greater than 18")
			},
		},
		{
			name: "should return error for tag lte",
			args: args{
				dataSet: struct {
					Age int `validate:"lte=18"`
				}{
					Age: 19,
				},
			},
			checkError: func(err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "The field 'Age' should be less than or equal 18")
			},
		},
		{
			name: "should return error for tag lt",
			args: args{
				dataSet: struct {
					Age int `validate:"lt=18"`
				}{
					Age: 18,
				},
			},
			checkError: func(err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "The field 'Age' should be less than 18")
			},
		},
		{
			name: "should return error for tag max",
			args: args{
				dataSet: struct {
					Age int `validate:"max=18"`
				}{
					Age: 19,
				},
			},
			checkError: func(err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "The field 'Age' should have the max lenhgt or value: 18")
			},
		},
		{
			name: "should return error for tag min",
			args: args{
				dataSet: struct {
					Age int `validate:"min=18"`
				}{
					Age: 17,
				},
			},
			checkError: func(err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "The field 'Age' should have the minimun lenhgt or value: 18")
			},
		},
		{
			name: "should return error for tag uuid4",
			args: args{
				dataSet: struct {
					UUID string `validate:"uuid4"`
				}{
					UUID: "123",
				},
			},
			checkError: func(err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "The format of 'UUID' should be uuid4: ")
			},
		},
		{
			name: "should return error for tag cpf",
			args: args{
				dataSet: struct {
					CPF string `validate:"cpf"`
				}{
					CPF: "12345678910",
				},
			},
			checkError: func(err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "The field 'CPF' should be a valid cpf")
			},
		},
		{
			name: "should return error for tag cnpj",
			args: args{
				dataSet: struct {
					CNPJ string `validate:"cnpj"`
				}{
					CNPJ: "1234567891011",
				},
			},
			checkError: func(err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "The field 'CNPJ' should be a valid cnpj")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateStruct(ctx, tt.args.dataSet)
			tt.checkError(err)
		})
	}
}
