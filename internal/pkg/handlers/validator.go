package handlers

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func ValidateInputs(dataSet interface{}) (bool, map[string]string) {

	validate = validator.New()
	err := validate.Struct(dataSet)

	if err != nil {
		if err, ok := err.(*validator.InvalidValidationError); ok {
			panic(err)
		}

		errors := make(map[string]string)
		reflected := reflect.ValueOf(dataSet)

		for _, err := range err.(validator.ValidationErrors) {
			field, _ := reflected.Type().FieldByName(err.StructField())

			var name string
			if name = field.Tag.Get("json"); name == "" {
				name = strings.ToLower(err.StructField())
			}

			switch err.ActualTag() {
			case "required":
				errors[name] = "Required " + name
				break

			case "email":
				errors[name] = "Email must be in the following format: mail@example.com"
				break
			}

			return false, errors
		}
	}
	return true, nil
}
