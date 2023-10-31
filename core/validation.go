package core

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func InitValidation() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

func ValidateStruct(s interface{}) map[string]string {
	validationResult := validate.Struct(s)

	structErrors := make(map[string]string)
	if validationResult != nil {
		for _, err := range validationResult.(validator.ValidationErrors) {
			structErrors[err.Field()] = err.Tag()
		}
	}

	return structErrors
}
