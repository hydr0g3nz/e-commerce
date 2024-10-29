package handlers

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func formatValidationError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("Minimum length is %s", err.Param())
	default:
		return "Invalid value"
	}
}
