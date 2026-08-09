package validator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"

	"github.com/Nurlan270/weather-go/internal/view"
)

type Validate struct {
	*validator.Validate
}

func New(validate *validator.Validate) *Validate {
	return &Validate{validate}
}

func (v *Validate) MapErrors(err error) view.ErrorMessages {
	var vErrors validator.ValidationErrors

	errors.As(err, &vErrors)

	var messages = make(view.ErrorMessages)

	for _, e := range vErrors {
		messages[e.StructField()] = v.getMessage(e)
	}

	return messages
}

func (v *Validate) getMessage(f validator.FieldError) string {
	fName := v.formatFieldName(f.StructField())

	switch f.Tag() {
	case "required":
		return fmt.Sprintf(Required, fName)
	case "max":
		return fmt.Sprintf(Max, fName, f.Param())
	case "min":
		return fmt.Sprintf(Min, fName, f.Param())
	case "eqfield":
		return fmt.Sprintf(Eqfield, fName, f.Param())
	case "login":
		return fmt.Sprintf(Login, fName)
	default:
		return fmt.Sprintf(Default, fName, f.Tag())
	}
}

// formatFieldName formats field's name
// from "SomeFieldName" to "Some field name".
func (v *Validate) formatFieldName(n string) string {
	s := strcase.ToDelimited(n, ' ')
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}
