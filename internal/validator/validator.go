package validator

import (
	"errors"
	"fmt"
	"github.com/Nurlan270/weather-go/internal/renderer/response"
	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
	"strings"
)

type Validate struct {
	*validator.Validate
}

func New(validate *validator.Validate) *Validate {
	return &Validate{validate}
}

func (v *Validate) MapErrors(err error) response.ErrorMessages {
	var vErrors validator.ValidationErrors
	errors.As(err, &vErrors)

	var messages = make(response.ErrorMessages)

	for _, e := range vErrors {
		messages[e.StructField()] = v.getMessage(e)
	}

	return messages
}

func (v *Validate) getMessage(f validator.FieldError) string {
	fName := v.formatFieldName(f.StructField())

	switch f.Tag() {
	case "required":
		return fmt.Sprintf("%s field is required", fName)
	case "max":
		return fmt.Sprintf("%s field must be at most %s characters long", fName, f.Param())
	case "min":
		return fmt.Sprintf("%s field must be at least %s characters long", fName, f.Param())
	case "eqfield":
		return fmt.Sprintf("%s field must be equal to %s field", fName, f.Param())
	case "login":
		return fmt.Sprintf(
			"%s field must be a valid email address or a username, allowed following characters: ._-", fName)
	default:
		return fmt.Sprintf("%s field failed on %s validation rule", fName, f.Tag())
	}
}

// formatFieldName formats field's name
// from "SomeFieldName" to "Some field name"
func (v *Validate) formatFieldName(n string) string {
	s := strcase.ToDelimited(n, ' ')
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}
